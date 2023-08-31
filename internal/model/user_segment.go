package model

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

// UserSegment describes users to segments relation model.
type UserSegment struct {
	UserID    uint
	User      User `gorm:"foreignKey:UserID;references:ID"`
	SegmentID uint
	Segment   Segment `gorm:"foreignKey:SegmentID;references:ID"`
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt `sql:"index"`
}

// UserSegmentsInput describes input params for updating user's segments.
type UserSegmentsInput struct {
	UserID           uint64 `json:"-" swaggerignore:"true"`
	SegmentsToAdd    []Slug `json:"segments_to_add" example:"PROTECTED_PHONE_NUMBER,VOICE_MSG"`
	SegmentsToDelete []Slug `json:"segments_to_delete" example:"PROMO_5"`
}

// Bind implements render.Binder interface method.
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

// UserSegmentsHistoryInput describes input path/query params
// for generating user's segments history.
type UserSegmentsHistoryInput struct {
	UserID   uint64
	FromDate time.Time
	ToDate   time.Time
}

// FromURI gets and checks input params from request.
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
	if !u.FromDate.Before(u.ToDate) {
		return ErrInvalidDateInterval
	}
	return nil
}

// UserSegmentsHistoryOutput describes json response of gen history response.
type UserSegmentsHistoryOutput struct {
	Link string `json:"link"`
}

// Render implements render.Render interface method.
func (u UserSegmentsHistoryOutput) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
