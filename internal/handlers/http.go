// Package handlers describes server router and handlers methods.
package handlers

import (
	"errors"
	"fmt"
	"net/http"

	logger "github.com/chi-middleware/logrus-logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	log "github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	_ "github.com/unbeman/av-prac-task/docs"
	"github.com/unbeman/av-prac-task/internal/database"
	"github.com/unbeman/av-prac-task/internal/model"
	"github.com/unbeman/av-prac-task/internal/services"
)

type HTTPHandler struct {
	*chi.Mux
	userService    *services.UserService
	segmentService *services.SegmentService
}

// GetHandler setups and returns HTTPHandler.
func GetHandler(userService *services.UserService, segmentService *services.SegmentService) (*HTTPHandler, error) {
	h := &HTTPHandler{
		Mux:            chi.NewMux(),
		userService:    userService,
		segmentService: segmentService,
	}

	h.Use(logger.Logger("router", log.New()))
	h.Get("/swagger/*", httpSwagger.Handler())
	h.Route("/api/v1/", func(router chi.Router) {
		router.Route("/segments/user", func(r chi.Router) {
			r.Get("/{user_id}", h.GetActiveUserSegments)
			r.Get("/{user_id}/csv", h.GenerateUserSegmentsHistory)
			r.Post("/{user_id}", h.UpdateUserSegments)
			r.Get("/history/{filename}", h.GetUserSegmentsHistoryFile)
		})

		router.Route("/segment", func(r chi.Router) {
			r.Post("/", h.CreateSegment)
			r.Delete("/{slug}", h.DeleteSegment)
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
// @Router /segment [post]
func (h HTTPHandler) CreateSegment(writer http.ResponseWriter, request *http.Request) {
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
// @Produce json
// @Param slug path string true "slug"
// @Success 200
// @Failure 400 {object} model.OutputError
// @Failure 404 {object} model.OutputError
// @Failure 500 {object} model.OutputError
// @Router /segment/{slug} [delete]
func (h HTTPHandler) DeleteSegment(writer http.ResponseWriter, request *http.Request) {
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
// @Router /segments/user/{user_id} [post]
func (h HTTPHandler) UpdateUserSegments(writer http.ResponseWriter, request *http.Request) {
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
// @Produce json
// @Param user_id path uint true "User ID"
// @Success 200 {object} model.Slugs
// @Failure 400 {object} model.OutputError
// @Failure 404 {object} model.OutputError
// @Failure 500 {object} model.OutputError
// @Router /segments/user/{user_id} [get]
func (h HTTPHandler) GetActiveUserSegments(writer http.ResponseWriter, request *http.Request) {
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

// GenerateUserSegmentsHistory godoc
// @Summary Get user's segments history link to download
// @Produce json
// @Param user_id path uint true "User ID"
// @Param from	query string true "From Date" Format(date) Example("2023-08-01")
// @Param to	query string true "To Date" Format(date) Example("2023-09-01")
// @Success 200
// @Failure 400 {object} model.OutputError
// @Failure 404 {object} model.OutputError
// @Failure 500 {object} model.OutputError
// @Router /segments/user/{user_id}/csv [get]
func (h HTTPHandler) GenerateUserSegmentsHistory(writer http.ResponseWriter, request *http.Request) {
	input := &model.UserSegmentsHistoryInput{}

	err := input.FromURI(request)
	if err != nil {
		h.processError(writer, request, fmt.Errorf("%w: %s", ErrInvalidRequest, err))
		return
	}

	//todo: go generateHist
	//todo: return file link to hist + status accepted

	filename, err := h.userService.GenerateUserSegmentsHistoryFile(request.Context(), input)
	if err != nil {
		h.processError(writer, request, err)
		return
	}

	render.Status(request, http.StatusOK)
	render.Render(writer, request, model.UserSegmentsHistoryOutput{Link: fmt.Sprintf("%s/history/%s", request.Host, filename)})
}

// GetUserSegmentsHistoryFile godoc
// @Summary Get user's segments history csv file
// @Produce text/csv
// @Param filename path string true "file name"
// @Success 200
// @Failure 400 {object} model.OutputError
// @Failure 404 {object} model.OutputError
// @Failure 500 {object} model.OutputError
// @Router /segments/user/history/{filename} [get]
func (h HTTPHandler) GetUserSegmentsHistoryFile(writer http.ResponseWriter, request *http.Request) {
	filename := chi.URLParam(request, "filename")

	log.Infof("requesting file: %s", filename)

	filePath, err := h.userService.DownloadUserSegmentsHistory(filename)
	if err != nil {
		h.processError(writer, request, err)
	}

	http.ServeFile(writer, request, filePath)
}

// processError handle app errors.
func (h HTTPHandler) processError(w http.ResponseWriter, r *http.Request, err error) {
	var httpCode int
	switch {
	case errors.Is(err, ErrInvalidRequest):
		httpCode = http.StatusBadRequest
	case errors.Is(err, database.ErrAlreadyExists):
		httpCode = http.StatusConflict
	case errors.Is(err, database.ErrNotFound):
		httpCode = http.StatusNotFound
	case errors.Is(err, services.ErrFileNotFound):
		httpCode = http.StatusNotFound
	default:
		httpCode = http.StatusInternalServerError
	}
	oErr := model.OutputError{HttpCode: httpCode, Message: err.Error()}

	render.Render(w, r, oErr)
}
