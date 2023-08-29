package model

import "net/http"

type User struct {
	Base
	Name     string   `json:"name"`
	Segments Segments `json:"segments,omitempty" gorm:"many2many:user_segments;"`
}

type UserInput struct {
	UserID uint `json:"user_id"`
}

func (u *UserInput) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (u *UserInput) Bind(r *http.Request) error {
	return nil
}
