package handlers

import (
	"encoding/csv"
	"errors"
	"fmt"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"net/http"
	"strconv"

	logger "github.com/chi-middleware/logrus-logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	log "github.com/sirupsen/logrus"

	_ "github.com/unbeman/av-prac-task/docs"
	"github.com/unbeman/av-prac-task/internal/database"
	"github.com/unbeman/av-prac-task/internal/model"
	"github.com/unbeman/av-prac-task/internal/services"
)

type HTTPHandlers struct {
	*chi.Mux
	userService    *services.UserService
	segmentService *services.SegmentService
}

func GetHandlers(userService *services.UserService, segmentService *services.SegmentService) (*HTTPHandlers, error) {
	h := &HTTPHandlers{
		Mux:            chi.NewMux(),
		userService:    userService,
		segmentService: segmentService,
	}
	h.Use(logger.Logger("router", log.New()))
	h.Get("/swagger/*", httpSwagger.Handler())
	h.Route("/api/v1/", func(router chi.Router) {
		router.Route("/user", func(r chi.Router) {
			r.Get("/segments", h.GetActiveUserSegments)
			r.Get("/segments/csv", h.GetUserSegmentsHistory)
			r.Post("/segments", h.UpdateUserSegments)
		})

		router.Route("/segment", func(r chi.Router) {
			r.Post("/", h.CreateSegment)
			r.Delete("/", h.DeleteSegment)
			r.Get("/", h.GetSegment)
		})
	})
	return h, nil
}

// CreateSegment godoc
// @Summary Creates new segment with given slug
// @Accept json
// @Produce json
// @Param slug body model.CreateSegment true "Segment input"
// @Success 200
// @Failure 400 {object} model.OutputError
// @Failure 409 {object} model.OutputError
// @Failure 500 {object} model.OutputError
// @Router /api/v1/segment [post]
func (h HTTPHandlers) CreateSegment(writer http.ResponseWriter, request *http.Request) {
	input := &model.CreateSegment{}
	err := render.Bind(request, input)
	if err != nil {
		log.Info(err)
		h.processError(writer, request, fmt.Errorf("%w: %v", ErrInvalidRequest, err))
		return
	}

	_, err = h.segmentService.CreateSegment(request.Context(), input)
	if err != nil {
		h.processError(writer, request, err)
		return
	}

	render.Status(request, http.StatusOK) // todo: status accepted for auto user adding
}

// DeleteSegment godoc
// @Summary Deletes segment with given slug
// @Accept json
// @Produce json
// @Param slug body model.CreateSegment true "Segment input"
// @Success 200
// @Failure 400 {object} model.OutputError
// @Failure 404 {object} model.OutputError
// @Failure 500 {object} model.OutputError
// @Router /api/v1/segment [delete]
func (h HTTPHandlers) DeleteSegment(writer http.ResponseWriter, request *http.Request) {
	segment := &model.SegmentInput{}
	err := render.Bind(request, segment)
	if err != nil {
		h.processError(writer, request, fmt.Errorf("%w: %v", ErrInvalidRequest, err))
		return
	}

	err = h.segmentService.DeleteSegment(request.Context(), segment)
	if err != nil {
		h.processError(writer, request, err)
		return
	}

	render.Status(request, http.StatusOK)
}

// todo: remove
func (h HTTPHandlers) GetSegment(writer http.ResponseWriter, request *http.Request) {
	input := &model.SegmentInput{}
	err := render.Bind(request, input)
	if err != nil {
		h.processError(writer, request, fmt.Errorf("%w: %v", ErrInvalidRequest, err))
		return
	}

	segment, err := h.segmentService.GetSegment(request.Context(), input)
	if err != nil {
		h.processError(writer, request, err)
		return
	}

	render.Status(request, http.StatusOK)
	render.Render(writer, request, segment)
}

// UpdateUserSegments godoc
// @Summary Updates user's segments
// @Accept json
// @Produce json
// @Param input body model.UserSegmentsInput true "User segments input"
// @Success 200
// @Failure 400 {object} model.OutputError
// @Failure 404 {object} model.OutputError
// @Failure 500 {object} model.OutputError
// @Router /api/v1/user/segments [post]
func (h HTTPHandlers) UpdateUserSegments(writer http.ResponseWriter, request *http.Request) {
	input := &model.UserSegmentsInput{}
	err := render.Bind(request, input)
	if err != nil {
		h.processError(writer, request, fmt.Errorf("%w: %v", ErrInvalidRequest, err))
		return
	}

	err = h.userService.UpdateUserSegments(request.Context(), input)
	if err != nil {
		h.processError(writer, request, err)
		return
	}

	render.Status(request, http.StatusOK)
}

// GetActiveUserSegments godoc
// @Summary Get user's active segments
// @Accept json
// @Produce json
// @Param input body model.User true "User"
// @Success 200 {object} model.Segments
// @Failure 400 {object} model.OutputError
// @Failure 404 {object} model.OutputError
// @Failure 500 {object} model.OutputError
// @Router /api/v1/user/segments [get]
func (h HTTPHandlers) GetActiveUserSegments(writer http.ResponseWriter, request *http.Request) {
	input := &model.UserInput{}
	err := render.Bind(request, input)
	if err != nil {
		h.processError(writer, request, fmt.Errorf("%w: %v", ErrInvalidRequest, err))
		return
	}

	user, err := h.userService.GetUserWithActiveSegments(request.Context(), input)
	if err != nil {
		h.processError(writer, request, err)
		return
	}

	render.Status(request, http.StatusOK)
	render.Render(writer, request, user.Segments)
}

// todo: desc
func (h HTTPHandlers) GetUserSegmentsHistory(writer http.ResponseWriter, request *http.Request) {
	input := &model.UserSegmentsHistoryInput{}
	err := render.Bind(request, input)
	if err != nil {
		h.processError(writer, request, fmt.Errorf("%w: %v", ErrInvalidRequest, err))
		return
	}

	//todo: go generateHist
	//todo: return file link to hist + status accepted

	user, err := h.userService.GetUserWithSegmentsHistory(request.Context(), input)
	if err != nil {
		h.processError(writer, request, err)
		return
	}

	history := make([][]string, 0, len(user.Segments)+1)

	head := []string{"user_id", "segment_slug", "operation", "date"}

	history = append(history, head)

	for _, s := range user.Segments {
		operation := "add"
		date := s.CreatedAt
		if s.DeletedAt != nil {
			operation = "delete"
			date = *s.DeletedAt
		}
		row := []string{strconv.Itoa(int(user.ID)), string(s.Slug), operation, date.String()}
		history = append(history, row)
	}

	hw := csv.NewWriter(writer)
	hw.WriteAll(history)

	render.Status(request, http.StatusOK)
}

func (h HTTPHandlers) processError(w http.ResponseWriter, r *http.Request, err error) {
	var httpCode int
	switch {
	case errors.Is(err, ErrInvalidRequest):
		httpCode = http.StatusBadRequest
	case errors.Is(err, database.ErrAlreadyExists):
		httpCode = http.StatusConflict
	case errors.Is(err, database.ErrNotFound):
		httpCode = http.StatusNotFound
	default:
		httpCode = http.StatusInternalServerError
	}
	oErr := model.OutputError{HttpCode: httpCode, Message: err.Error()}

	render.Render(w, r, oErr)
}
