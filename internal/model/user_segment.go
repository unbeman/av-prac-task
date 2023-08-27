package model

import (
	"net/http"
	"time"
)

type UserSegment struct {
	Base
	UserID    uint
	User      User `gorm:"foreignKey:UserID;references:ID"`
	SegmentID uint
	Segment   Segment `gorm:"foreignKey:SegmentID;references:ID"`
	TTL       time.Time
}

type UserSegmentsInput struct {
	UserID           uint   `json:"user_id"`
	SegmentsToAdd    []Slug `json:"segments_to_add"`
	SegmentsToDelete []Slug `json:"segments_to_delete"`
}

func (u *UserSegmentsInput) Bind(r *http.Request) error {
	for _, s := range u.SegmentsToAdd {
		if err := s.Bind(r); err != nil {
			return err
		}
	}
	for _, s := range u.SegmentsToDelete {
		if err := s.Bind(r); err != nil {
			return err
		}
	}
	return nil
}
