package entity

import (
	"time"
)

type Service struct {
	Id          int
	Name        string
	Description *string
	DomainId    int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
