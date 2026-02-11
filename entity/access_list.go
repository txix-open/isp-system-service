package entity

type AccessList struct {
	AppId      int
	HttpMethod string
	Method     string
	Value      bool
}

type Method struct {
	HttpMethod string
	Method     string
}
