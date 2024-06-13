package domain

import (
	"time"
)

type DomainCreateUpdateRequest struct {
	Id          int
	Name        string `valid:"required~Required"`
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

type DomainWithApplicationGroup struct {
	Id               int
	Name             string
	Description      string
	ApplicationGroup []*ApplicationGroupWithApps
}
