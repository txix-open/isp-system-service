package entity

import (
	"time"
)

type ApplicationGroup struct {
	Id          int
	Name        string
	Description *string
	DomainId    int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
