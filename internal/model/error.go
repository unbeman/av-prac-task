package model

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
)

// Error that are used for parsing of incoming parameters.
// Used in Bind and FroURI input model methods.
var (
	ErrInvalidSlug         = errors.New("invalid slug")
	ErrInvalidSelection    = errors.New("invalid user selection value")
	ErrInvalidUserID       = errors.New("invalid userID")
	ErrInvalidDateFormat   = errors.New("invalid date format")
	ErrInvalidDateInterval = errors.New("invalid date interval")
)

// OutputError describes json response for error.
type OutputError struct {
	HttpCode int    `json:"-"`
	Message  string `json:"message" example:"error message"`
}

// Render implements render.Render interface method.
func (o OutputError) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, o.HttpCode)
	return nil
}
