package domain

type AccessListSetOneRequest struct {
	AppId  int    `valid:"required~Required"`
	Method string `valid:"required~Required"`
	Value  bool
}

type AccessListSetOneResponse struct {
	Count int
}

type AccessListSetListRequest struct {
	AppId     int `valid:"required~Required"`
	RemoveOld bool
	Methods   []MethodInfo
}

type MethodInfo struct {
	Method string
	Value  bool
}
