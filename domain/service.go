package domain

import (
	"time"
)

type Service struct {
	Id          int
	Name        string
	Description string
	DomainId    int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ServiceCreateUpdateRequest struct {
	Id          int
	Name        string `validate:"required"`
	Description string
	DomainId    int `validate:"required"`
}

type ServiceWithApps struct {
	Id          int
	Name        string
	Description string
	Apps        []*ApplicationSimple
}
