package entity

import "time"

const (
	RootAdminApplicationId = 1
)

const (
	AppTypeSystem = "SYSTEM"
	AppTypeMobile = "MOBILE"
)

type System struct {
	TableName   string `sql:"system_service.system" json:"-"`
	Id          int32
	Name        string `valid:"required~Required"`
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Domain struct {
	TableName   string `sql:"system_service.domain" json:"-"`
	Id          int32
	Name        string `valid:"required~Required"`
	Description string
	SystemId    int32
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Service struct {
	TableName   string `sql:"system_service.service" json:"-"`
	Id          int32
	Name        string `valid:"required~Required"`
	Description string
	DomainId    int32 `valid:"required~Required"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Application struct {
	TableName   string `sql:"system_service.application" json:"-"`
	Id          int32
	Name        string `valid:"required~Required"`
	Description string
	Type        string `valid:"required~Required,in(SYSTEM|MOBILE)"`
	ServiceId   int32  `valid:"required~Required"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Token struct {
	TableName  string `sql:"system_service.token" json:"-"`
	Token      string `valid:"required~Required" sql:"pk:token"`
	AppId      int32  `valid:"required~Required"`
	ExpireTime int64
	CreatedAt  time.Time
}

type AccessList struct {
	TableName string `sql:"system_service.access_list" json:"-"`
	AppId     int32  `sql:",pk"`
	Method    string `sql:",pk"`
	Value     bool
}
