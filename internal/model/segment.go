package model

import (
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"time"
)

type Slug string

func (s *Slug) Bind(r *http.Request) error {
	if *s == "" {
		return ErrInvalidSlug
	}

	*s = Slug(strings.ToUpper(string(*s)))
	return nil
}

type Slugs []Slug

func (s Slugs) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type Segment struct {
	ID        uint64 `json:"id" gorm:"primary_key"`
	Slug      Slug   `json:"slug" gorm:"uniqueIndex"`
	Users     []User `json:"users,omitempty" gorm:"many2many:user_segments;"`
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt `sql:"index"`
}

func (s *Segment) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type CreateSegmentInput struct {
	Slug      Slug     `json:"slug" example:"AVITO_VOICE_MESSAGES"`
	Selection *float64 `json:"selection,omitempty" example:"0.2"`
}

func (s *CreateSegmentInput) Bind(r *http.Request) error {
	if err := s.Slug.Bind(r); err != nil { //no idea why it doesn't work without explicit call
		return err
	}
	if s.Selection != nil && (*s.Selection > 1.0 || *s.Selection < 0.0) {
		return ErrInvalidSelection
	}
	return nil
}

type SegmentInput struct {
	Slug Slug `json:"slug" example:"AVITO_VOICE_MESSAGES"`
}

func (s *SegmentInput) FromURI(r *http.Request) error {
	s.Slug = Slug(chi.URLParam(r, "slug"))
	if s.Slug == "" {
		return ErrInvalidSlug
	}
	return nil
}
