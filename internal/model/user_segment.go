package model

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

type UserSegment struct {
	UserID    uint
	User      User `gorm:"foreignKey:UserID;references:ID"`
	SegmentID uint
	Segment   Segment `gorm:"foreignKey:SegmentID;references:ID"`
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt `sql:"index"`
}

type UserSegmentsInput struct {
	UserID           uint64 `json:"-" swaggerignore:"true"`
	SegmentsToAdd    []Slug `json:"segments_to_add"`
	SegmentsToDelete []Slug `json:"segments_to_delete"`
}

func (u *UserSegmentsInput) Bind(r *http.Request) error {
	idParam := chi.URLParam(r, "user_id")
	userID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return ErrInvalidUserID
	}

	u.UserID = userID

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
	UserID   uint64    `json:"user_id" swaggerignore:"true"`
	FromDate time.Time `json:"from_date"` //todo: mb use timestamp
	ToDate   time.Time `json:"to_date"`
}

func (u *UserSegmentsHistoryInput) FromURI(r *http.Request) error {
	idParam := chi.URLParam(r, "user_id")
	userID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return ErrInvalidUserID
	}

	fromDate := r.URL.Query().Get("from")
	toDate := r.URL.Query().Get("to")

	u.UserID = userID
	u.FromDate, err = time.Parse(time.DateOnly, fromDate)
	if err != nil {
		return ErrInvalidDateFormat
	}
	u.ToDate, err = time.Parse(time.DateOnly, toDate)
	if err != nil {
		return ErrInvalidDateFormat
	}
	return nil
}

type UserSegmentsHistoryOutput struct {
	Link string `json:"link"`
}

func (u UserSegmentsHistoryOutput) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
