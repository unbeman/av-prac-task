package model

import (
	"fmt"
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

	uniqueAddSegs := make(map[Slug]bool)
	for _, segment := range u.SegmentsToAdd {
		if _, ok := uniqueAddSegs[segment]; !ok {
			uniqueAddSegs[segment] = true
		} else {
			return fmt.Errorf("%w: dublicating the segment to add", ErrInvalidSlug)
		}

	}

	uniqueDelSegs := make(map[Slug]bool)
	for _, segment := range u.SegmentsToDelete {
		if _, ok := uniqueDelSegs[segment]; !ok {
			uniqueDelSegs[segment] = true
		} else {
			return fmt.Errorf("%w: dublicating the segment to delete", ErrInvalidSlug)
		}
		if _, ok := uniqueAddSegs[segment]; ok {
			return fmt.Errorf("%w: segment to add intersects segment for delete", ErrInvalidSlug)
		}
	}

	return nil
}

type UserSegmentsHistoryInput struct {
	UserID   uint      `json:"user_id"`
	FromDate time.Time `json:"from_date"` //todo: mb use timestamp
	ToDate   time.Time `json:"to_date"`
}

func (u *UserSegmentsHistoryInput) Bind(r *http.Request) error {
	//todo: check dates with format
	return nil
}
