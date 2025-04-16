package entity

import (
	"database/sql"
	"time"
)

type AppGroup struct {
	Id          int
	Name        string
	Description sql.NullString
	DomainId    int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
