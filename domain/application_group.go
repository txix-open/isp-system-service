package domain

import (
	"time"
)

type ApplicationGroup struct {
	Id          int
	Name        string
	Description string
	DomainId    int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ApplicationGroupCreateUpdateRequest struct {
	Id          int
	Name        string `valid:"required~Required"`
	Description string
	DomainId    int `valid:"required~Required"`
}

type ApplicationGroupWithApps struct {
	Id          int
	Name        string
	Description string
	Apps        []*ApplicationSimple
}
