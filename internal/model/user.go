package model

import "net/http"

type User struct {
	Base
	Name     string   `json:"name"`
	Segments Segments `json:"segments,omitempty" gorm:"many2many:user_segments;"`
}

func (u *User) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (u *User) Bind(r *http.Request) error {
	return nil
}
