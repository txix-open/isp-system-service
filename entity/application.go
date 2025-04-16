package entity

import (
	"database/sql"
	"time"
)

type Application struct {
	Id                 int
	Name               string
	Description        sql.NullString
	ApplicationGroupId int
	Type               string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
