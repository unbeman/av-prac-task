package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/unbeman/av-prac-task/internal/database"
	mock_database "github.com/unbeman/av-prac-task/internal/database/mock"
	"github.com/unbeman/av-prac-task/internal/model"
	"github.com/unbeman/av-prac-task/internal/services"
)

func setupHandler(t *testing.T, ctrl *gomock.Controller, setupDB func(db *mock_database.MockIDatabase)) *HTTPHandlers {
	database := mock_database.NewMockIDatabase(ctrl)
	setupDB(database)

	segmentServ, err := services.NewSegmentService(database)
	require.NoError(t, err)

	userServ, err := services.NewUserService(database)
	require.NoError(t, err)

	h, err := GetHandlers(userServ, segmentServ)
	require.NoError(t, err)

	return h
}

func generateSegment() model.Segment {
	userSel := rand.Float64()
	return model.Segment{
		Slug:      "GOOD-NAME",
		Selection: &userSel,
	}
}

func getSelection(n float64) *float64 {
	return &n
}

func TestHTTPHandlers_CreateSegment(t *testing.T) {
	segment := generateSegment()
	tests := []struct {
		name          string
		input         model.CreateSegment
		buildStubs    func(db *mock_database.MockIDatabase)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:  "OK",
			input: model.CreateSegment{Slug: segment.Slug, Selection: segment.Selection},
			buildStubs: func(db *mock_database.MockIDatabase) {
				db.EXPECT().
					CreateSegment(gomock.Any(), gomock.Any()).
					Return(&segment, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:  "Internal Error",
			input: model.CreateSegment{Slug: segment.Slug, Selection: segment.Selection},
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
			input: model.CreateSegment{Slug: segment.Slug, Selection: segment.Selection},
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
			input: model.CreateSegment{Slug: segment.Slug, Selection: getSelection(1.25)},
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
			input: model.CreateSegment{Slug: "", Selection: segment.Selection},
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
			t.Log(string(data))

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
	segment := generateSegment()
	tests := []struct {
		name          string
		input         model.CreateSegment
		buildStubs    func(db *mock_database.MockIDatabase)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:  "OK",
			input: model.CreateSegment{Slug: segment.Slug},
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
			input: model.CreateSegment{Slug: segment.Slug},
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
			name:  "Segment not found",
			input: model.CreateSegment{Slug: segment.Slug},
			buildStubs: func(db *mock_database.MockIDatabase) {
				db.EXPECT().
					DeleteSegment(gomock.Any(), gomock.Any()).
					Return(database.ErrNotFound)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:  "Empty slug",
			input: model.CreateSegment{Slug: ""},
			buildStubs: func(db *mock_database.MockIDatabase) {
				db.EXPECT().
					DeleteSegment(gomock.Any(), gomock.Any()).
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

			request, err := http.NewRequest(http.MethodDelete, "/api/v1/segment", bytes.NewBuffer(data))
			require.NoError(t, err)
			request.Header.Set("Content-Type", "application/json")

			//handler.DeleteSegment(recorder, request)

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

			request, err := http.NewRequest(http.MethodPost, "/api/v1/user/segments", bytes.NewBuffer(data))
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
	user.Segments = model.Segments{&segA, &segB}

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
					GetUserActiveSegments(gomock.Any(), gomock.Any()). //todo: fill
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
					GetUserActiveSegments(gomock.Any(), gomock.Any()). //todo: fill
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
					GetUserActiveSegments(gomock.Any(), gomock.Any()).
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

			data, err := json.Marshal(tt.input)
			require.NoError(t, err)
			t.Log(string(data))

			request, err := http.NewRequest(http.MethodGet, "/api/v1/user/segments", bytes.NewBuffer(data))
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
