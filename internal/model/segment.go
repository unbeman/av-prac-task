package model

import (
	"net/http"
	"strings"
)

type Slug string

func (s *Slug) Bind(r *http.Request) error {
	if *s == "" { //todo: add slug validator
		return ErrInvalidSlug
	}

	*s = Slug(strings.ToUpper(string(*s)))
	return nil
}

type Segment struct {
	Base
	Slug      Slug     `json:"slug" gorm:"uniqueIndex"`
	Selection *float64 `json:"selection,omitempty" ` // 0 < user selection <= 1
	Users     []*User  `json:"users,omitempty" gorm:"many2many:user_segments;"`
}

func (s *Segment) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type Segments []*Segment

func (s Segments) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type SegmentInput struct {
	Slug      Slug     `json:"slug"`
	Selection *float64 `json:"selection,omitempty"`
}

func (s *SegmentInput) Bind(r *http.Request) error {
	if err := s.Slug.Bind(r); err != nil { //no idea why it doesn't work without explicit call
		return err
	}
	if s.Selection != nil && (*s.Selection > 1.0 || *s.Selection < 0.0) {
		return ErrInvalidSelection
	}
	return nil
}
