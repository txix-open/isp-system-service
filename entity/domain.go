package entity

import (
	"time"
)

type Domain struct {
	Id          int
	Name        string
	Description *string
	SystemId    int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
