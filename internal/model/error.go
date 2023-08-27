package model

import (
	"errors"
	"github.com/go-chi/render"
	"net/http"
)

var (
	ErrInvalidSlug      = errors.New("invalid slug")
	ErrInvalidSelection = errors.New("invalid user selection value")
	ErrInvalidUserID    = errors.New("invalid userID")
)

type OutputError struct {
	HttpCode int    `json:"-"`
	Message  string `json:"message"`
}

func (o OutputError) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, o.HttpCode)
	return nil
}
