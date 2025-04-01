package entity

import (
	"time"
)

type AppGroup struct {
	Id          int
	Name        string
	Description *string
	DomainId    int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
