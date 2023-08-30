package handlers

import (
	"encoding/csv"
	"errors"
	"fmt"
	logger "github.com/chi-middleware/logrus-logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	log "github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"net/http"

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
		router.Route("/segments/user", func(r chi.Router) {
			r.Get("/{user_id}", h.GetActiveUserSegments)
			r.Get("/{user_id}/csv", h.GetUserSegmentsHistory)
			r.Post("/{user_id}", h.UpdateUserSegments)
		})

		router.Route("/segment", func(r chi.Router) {
			r.Post("/", h.CreateSegment)
			r.Delete("/{slug}", h.DeleteSegment)
			r.Get("/{slug}", h.GetSegment)
		})
	})
	return h, nil
}

// CreateSegment godoc
// @Summary Creates new segment with given slug
// @Accept json
// @Produce json
// @Param segment body model.CreateSegmentInput true "Segment input"
// @Success 200
// @Failure 400 {object} model.OutputError
// @Failure 409 {object} model.OutputError
// @Failure 500 {object} model.OutputError
// @Router /api/v1/segment [post]
func (h HTTPHandlers) CreateSegment(writer http.ResponseWriter, request *http.Request) {
	input := &model.CreateSegmentInput{}
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
// @Param slug path string true "slug"
// @Success 200
// @Failure 400 {object} model.OutputError
// @Failure 404 {object} model.OutputError
// @Failure 500 {object} model.OutputError
// @Router /api/v1/segment/{slug} [delete]
func (h HTTPHandlers) DeleteSegment(writer http.ResponseWriter, request *http.Request) {
	segment := &model.SegmentInput{}

	err := segment.FromURI(request)
	if err != nil {
		h.processError(writer, request, fmt.Errorf("%w: %s", ErrInvalidRequest, err))
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

	err := input.FromURI(request)
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
// @Param user_id path uint true "User id"
// @Param input body model.UserSegmentsInput true "User segments input"
// @Success 200
// @Failure 400 {object} model.OutputError
// @Failure 404 {object} model.OutputError
// @Failure 500 {object} model.OutputError
// @Router /api/v1/segments/user/{user_id} [post]
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
// @Param user_id path uint true "User ID"
// @Success 200 {object} model.Slugs
// @Failure 400 {object} model.OutputError
// @Failure 404 {object} model.OutputError
// @Failure 500 {object} model.OutputError
// @Router /api/v1/segments/user/{user_id} [get]
func (h HTTPHandlers) GetActiveUserSegments(writer http.ResponseWriter, request *http.Request) {
	input := &model.UserInput{}

	err := input.FromURI(request)
	if err != nil {
		h.processError(writer, request, fmt.Errorf("%w: %s", ErrInvalidRequest, err))
		return
	}

	segments, err := h.userService.GetUserActiveSegments(request.Context(), input)
	if err != nil {
		h.processError(writer, request, err)
		return
	}

	render.Status(request, http.StatusOK)
	render.Render(writer, request, segments)
}

// GetUserSegmentsHistory godoc
// @Summary Get user's segments action history
// @Accept json
// @Produce json
// @Param user_id path uint true "User ID"
// @Param from	query string true "From Date"
// @Param to	query string true "To Date"
// @Success 200
// @Failure 400 {object} model.OutputError
// @Failure 404 {object} model.OutputError
// @Failure 500 {object} model.OutputError
// @Router /api/v1/segments/user/{user_id}/csv [get]
func (h HTTPHandlers) GetUserSegmentsHistory(writer http.ResponseWriter, request *http.Request) {
	input := &model.UserSegmentsHistoryInput{}
	log.Info("WE ARE IN GetUserSegmentsHistory")
	err := input.FromURI(request)
	if err != nil {
		h.processError(writer, request, fmt.Errorf("%w: %s", ErrInvalidRequest, err))
		return
	}

	//todo: go generateHist
	//todo: return file link to hist + status accepted

	history, err := h.userService.GetUserSegmentsHistory(request.Context(), input)
	if err != nil {
		h.processError(writer, request, err)
		return
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
