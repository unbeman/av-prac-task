package transport

import (
	"errors"
	"fmt"
	"net/http"

	logger "github.com/chi-middleware/logrus-logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	log "github.com/sirupsen/logrus"

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
	h.Route("/api/v1/", func(router chi.Router) {
		router.Route("/user", func(r chi.Router) {
			r.Get("/segments", h.GetActiveUserSegments)
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

func (h HTTPHandlers) CreateSegment(writer http.ResponseWriter, request *http.Request) {
	input := &model.SegmentInput{}
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

	render.Status(request, http.StatusOK) //StatusAccepted ?
}

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

func (h HTTPHandlers) GetActiveUserSegments(writer http.ResponseWriter, request *http.Request) {
	input := &model.User{}
	err := render.Bind(request, input)
	if err != nil {
		h.processError(writer, request, fmt.Errorf("%w: %v", ErrInvalidRequest, err))
		return
	}

	user, err := h.userService.GetActiveUserSegments(request.Context(), input)
	if err != nil {
		h.processError(writer, request, err)
		return
	}

	render.Status(request, http.StatusOK)
	render.Render(writer, request, user.Segments)
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
