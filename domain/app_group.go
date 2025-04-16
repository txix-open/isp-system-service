package domain

import "time"

type AppGroup struct {
	Id          int
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CreateAppGroupRequest struct {
	Name        string `validate:"required"`
	Description string
}

type UpdateAppGroupRequest struct {
	Id          int    `validate:"required"`
	Name        string `validate:"required"`
	Description string
}

type IdListRequest struct {
	IdList []int `validate:"required,min=1"`
}
