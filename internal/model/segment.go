package model

import (
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

// Segment describes segment model.
type Segment struct {
	ID        uint64 `json:"id" gorm:"primary_key"`
	Slug      Slug   `json:"slug" gorm:"uniqueIndex"`
	Users     []User `json:"users,omitempty" gorm:"many2many:user_segments;"`
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt `sql:"index"`
}

// Render implements render.Render interface method.
func (s *Segment) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// Slug type.
type Slug string

// Bind implements render.Binder interface method.
func (s *Slug) Bind(r *http.Request) error {
	if *s == "" {
		return ErrInvalidSlug
	}

	*s = Slug(strings.ToUpper(string(*s)))
	return nil
}

type Slugs []Slug

// Render implements render.Render interface method.
func (s Slugs) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// CreateSegmentInput describes json input for segment creation.
type CreateSegmentInput struct {
	Slug      Slug     `json:"slug" example:"AVITO_VOICE_MESSAGES"`
	Selection *float64 `json:"selection,omitempty" example:"0.2"`
}

// Bind implements render.Binder interface method.
func (s *CreateSegmentInput) Bind(r *http.Request) error {
	if err := s.Slug.Bind(r); err != nil { //no idea why it doesn't work without explicit call
		return err
	}
	if s.Selection != nil && (*s.Selection > 1.0 || *s.Selection < 0.0) {
		return ErrInvalidSelection
	}
	return nil
}

// SegmentInput describes path input to get/delete segment.
type SegmentInput struct {
	Slug Slug `example:"AVITO_VOICE_MESSAGES"`
}

// FromURI gets and checks segment input from url.
func (s *SegmentInput) FromURI(r *http.Request) error {
	s.Slug = Slug(chi.URLParam(r, "slug"))
	return s.Slug.Bind(r)
}
