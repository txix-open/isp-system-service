package domain

type AccessListSetOneRequest struct {
	AppId      int `validate:"required"`
	HttpMethod string
	Method     string `validate:"required"`
	Value      bool
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
	HttpMethod string
	Method     string
	Value      bool
}

type AccessListDeleteListRequest struct {
	AppId   int      `validate:"required"`
	Methods []string `validate:"required,min=1"`
}

type AccessListDeleteV2ListRequest struct {
	AppId   int      `validate:"required"`
	Methods []Method `validate:"required,min=1"`
}

type Method struct {
	HttpMethod string
	Method     string `validate:"required"`
}
