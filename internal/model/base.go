package model

import (
	"time"
)

type Base struct {
	ID        uint `json:"id" gorm:"primary_key"`
	CreatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}
