package model

import (
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

type User struct {
	ID        uint64    `json:"id" gorm:"primary_key"`
	Segments  []Segment `json:"segments,omitempty" gorm:"many2many:user_segments;"`
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt `sql:"index"`
}

type UserInput struct {
	UserID uint64 `json:"user_id"`
}

func (u *UserInput) FromURI(r *http.Request) error {
	idParam := chi.URLParam(r, "user_id")
	userID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return ErrInvalidUserID
	}
	u.UserID = userID
	return nil
}
