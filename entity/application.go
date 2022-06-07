package entity

import (
	"time"
)

type Application struct {
	Id          int
	Name        string
	Description *string
	ServiceId   int
	Type        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
