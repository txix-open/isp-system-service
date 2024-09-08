package domain

type AccessListSetOneRequest struct {
	AppId  int    `validate:"required"`
	Method string `validate:"required"`
	Value  bool
}

type AccessListSetOneResponse struct {
	Count int
}

type AccessListSetListRequest struct {
	AppId     int `validate:"required"`
	RemoveOld bool
	Methods   []MethodInfo
}

type MethodInfo struct {
	Method string
	Value  bool
}
