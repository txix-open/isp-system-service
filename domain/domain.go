package domain

import (
	"time"
)

type DomainCreateUpdateRequest struct {
	Id          int
	Name        string `validate:"required"`
	Description string
}

type Domain struct {
	Id          int
	Name        string
	Description string
	SystemId    int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type DomainWithService struct {
	Id          int
	Name        string
	Description string
	Services    []*ServiceWithApps
}
