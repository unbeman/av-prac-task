package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/unbeman/av-prac-task/internal/config"
	"github.com/unbeman/av-prac-task/internal/database"
	mock_database "github.com/unbeman/av-prac-task/internal/database/mock"
	"github.com/unbeman/av-prac-task/internal/model"
	"github.com/unbeman/av-prac-task/internal/services"
	"github.com/unbeman/av-prac-task/internal/worker"
)

func setupHandler(t *testing.T, ctrl *gomock.Controller, setupDB func(db *mock_database.MockIDatabase)) *HTTPHandler {
	database := mock_database.NewMockIDatabase(ctrl)
	setupDB(database)

	wp := worker.NewWorkersPool(config.NewWorkerPoolConfig())

	segmentServ, err := services.NewSegmentService(database)
	require.NoError(t, err)

	userServ, err := services.NewUserService(database, wp, t.TempDir())
	require.NoError(t, err)

	h, err := GetHandler(userServ, segmentServ)
	require.NoError(t, err)

	return h
}

func getSelection(n float64) *float64 {
	return &n
}

func TestHTTPHandlers_CreateSegment(t *testing.T) {
	segment := model.Segment{Slug: "SEGMENT-SLUG"}
	tests := []struct {
		name          string
		input         model.CreateSegmentInput
		buildStubs    func(db *mock_database.MockIDatabase)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:  "OK",
			input: model.CreateSegmentInput{Slug: segment.Slug, Selection: getSelection(0.5)},
			buildStubs: func(db *mock_database.MockIDatabase) {
				db.EXPECT().
					CreateSegment(gomock.Any(), gomock.Any()).
					Return(&segment, nil)
				db.EXPECT().
					AddSegmentToRandomUsers(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:  "Internal Error",
			input: model.CreateSegmentInput{Slug: segment.Slug},
			buildStubs: func(db *mock_database.MockIDatabase) {
				db.EXPECT().
					CreateSegment(gomock.Any(), gomock.Any()).
					Return(nil, database.ErrDB)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:  "Segment already exist",
			input: model.CreateSegmentInput{Slug: segment.Slug},
			buildStubs: func(db *mock_database.MockIDatabase) {
				db.EXPECT().
					CreateSegment(gomock.Any(), gomock.Any()).
					Return(nil, database.ErrAlreadyExists)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusConflict, recorder.Code)
			},
		},
		{
			name:  "Invalid selection value",
			input: model.CreateSegmentInput{Slug: segment.Slug, Selection: getSelection(1.25)},
			buildStubs: func(db *mock_database.MockIDatabase) {
				db.EXPECT().
					CreateSegment(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:  "Empty slug",
			input: model.CreateSegmentInput{Slug: ""},
			buildStubs: func(db *mock_database.MockIDatabase) {
				db.EXPECT().
					CreateSegment(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			handler := setupHandler(t, ctrl, tt.buildStubs)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tt.input)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, "/api/v1/segment", bytes.NewBuffer(data))
			require.NoError(t, err)
			request.Header.Set("Content-Type", "application/json")

			//handler.CreateSegment(recorder, request)

			handler.ServeHTTP(recorder, request)
			tt.checkResponse(t, recorder)
		})
	}
}

func TestHTTPHandlers_DeleteSegment(t *testing.T) {
	segment := model.Segment{Slug: "SEGMENT-SLUG"}

	tests := []struct {
		name          string
		input         model.SegmentInput
		buildStubs    func(db *mock_database.MockIDatabase)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:  "OK",
			input: model.SegmentInput{Slug: segment.Slug},
			buildStubs: func(db *mock_database.MockIDatabase) {
				db.EXPECT().
					DeleteSegment(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:  "Internal Error",
			input: model.SegmentInput{Slug: segment.Slug},
			buildStubs: func(db *mock_database.MockIDatabase) {
				db.EXPECT().
					DeleteSegment(gomock.Any(), gomock.Any()).
					Return(database.ErrDB)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:  "Empty slug",
			input: model.SegmentInput{Slug: ""},
			buildStubs: func(db *mock_database.MockIDatabase) {
				db.EXPECT().
					DeleteSegment(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusMethodNotAllowed, recorder.Code)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			handler := setupHandler(t, ctrl, tt.buildStubs)
			recorder := httptest.NewRecorder()

			request, err := http.NewRequest(
				http.MethodDelete,
				fmt.Sprintf("/api/v1/segment/%v", tt.input.Slug),
				nil,
			)
			require.NoError(t, err)

			handler.ServeHTTP(recorder, request)
			tt.checkResponse(t, recorder)
		})
	}
}

func TestHTTPHandlers_UpdateUserSegments(t *testing.T) {
	tests := []struct {
		name          string
		input         model.UserSegmentsInput
		buildStubs    func(db *mock_database.MockIDatabase)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			input: model.UserSegmentsInput{
				UserID:           1,
				SegmentsToAdd:    []model.Slug{"SEGMENT-3"},
				SegmentsToDelete: []model.Slug{"SEGMENT-1", "SEGMENT-2"},
			},
			buildStubs: func(db *mock_database.MockIDatabase) {
				db.EXPECT().
					CreateDeleteUserSegments(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()). //todo: fill
					Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Internal Error",
			input: model.UserSegmentsInput{
				UserID:           1,
				SegmentsToAdd:    []model.Slug{"SEGMENT-3"},
				SegmentsToDelete: []model.Slug{"SEGMENT-1", "SEGMENT-2"},
			},
			buildStubs: func(db *mock_database.MockIDatabase) {
				db.EXPECT().
					CreateDeleteUserSegments(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()). //todo: fill
					Return(database.ErrDB)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "User already in the segment",
			input: model.UserSegmentsInput{
				UserID:           1,
				SegmentsToAdd:    []model.Slug{"SEGMENT-3"},
				SegmentsToDelete: []model.Slug{"SEGMENT-1", "SEGMENT-2"},
			},
			buildStubs: func(db *mock_database.MockIDatabase) {
				db.EXPECT().
					CreateDeleteUserSegments(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(database.ErrAlreadyExists)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusConflict, recorder.Code)
			},
		},
		{
			name: "User not found",
			input: model.UserSegmentsInput{
				UserID:           1,
				SegmentsToAdd:    []model.Slug{"SEGMENT-3"},
				SegmentsToDelete: []model.Slug{"SEGMENT-1", "SEGMENT-2"},
			},
			buildStubs: func(db *mock_database.MockIDatabase) {
				db.EXPECT().
					CreateDeleteUserSegments(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(database.ErrNotFound)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "Empty segment slug",
			input: model.UserSegmentsInput{
				UserID:        1,
				SegmentsToAdd: []model.Slug{""},
			},
			buildStubs: func(db *mock_database.MockIDatabase) {
				db.EXPECT().
					CreateDeleteUserSegments(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			handler := setupHandler(t, ctrl, tt.buildStubs)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tt.input)
			require.NoError(t, err)
			t.Log(string(data))

			request, err := http.NewRequest(
				http.MethodPost,
				fmt.Sprintf("/api/v1/segments/user/%v", tt.input.UserID),
				bytes.NewBuffer(data),
			)
			require.NoError(t, err)
			request.Header.Set("Content-Type", "application/json")

			handler.ServeHTTP(recorder, request)
			tt.checkResponse(t, recorder)
		})
	}
}

func TestHTTPHandlers_GetActiveUserSegments(t *testing.T) {
	segA := model.Segment{Slug: "SEGMENT-A"}
	segB := model.Segment{Slug: "SEGMENT-B"}
	var user model.User
	user.ID = 1
	user.Segments = []model.Segment{segA, segB}

	tests := []struct {
		name          string
		input         model.User
		buildStubs    func(db *mock_database.MockIDatabase)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:  "OK",
			input: user,
			buildStubs: func(db *mock_database.MockIDatabase) {
				db.EXPECT().
					GetUserWithActiveSegments(gomock.Any(), gomock.Any()). //todo: fill
					Return(&user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:  "Internal Error",
			input: user,
			buildStubs: func(db *mock_database.MockIDatabase) {
				db.EXPECT().
					GetUserWithActiveSegments(gomock.Any(), gomock.Any()). //todo: fill
					Return(nil, database.ErrDB)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:  "User not found",
			input: user,
			buildStubs: func(db *mock_database.MockIDatabase) {
				db.EXPECT().
					GetUserWithActiveSegments(gomock.Any(), gomock.Any()).
					Return(nil, database.ErrNotFound)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			handler := setupHandler(t, ctrl, tt.buildStubs)
			recorder := httptest.NewRecorder()

			request, err := http.NewRequest(
				http.MethodGet,
				fmt.Sprintf("/api/v1/segments/user/%d", tt.input.ID),
				nil,
			)
			require.NoError(t, err)
			request.Header.Set("Content-Type", "application/json")

			handler.ServeHTTP(recorder, request)
			tt.checkResponse(t, recorder)
		})
	}
}

//func requireUserSegmentsEqual(t *testing.T, body *bytes.Buffer, segments model.Segments) {
//	data, err := io.ReadAll(body)
//	require.NoError(t, err)
//
//	var gotSegments model.Segments
//	err = json.Unmarshal(data, &gotSegments)
//	require.NoError(t, err)
//
//	require.Equal(t, segments, gotSegments) //todo: check slugs sets
//}
