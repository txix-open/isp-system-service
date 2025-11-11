package domain

type AuthenticateRequest struct {
	Token string `validate:"required"`
}

type AuthenticateResponse struct {
	Authenticated bool
	ErrorReason   string
	AuthData      *AuthData
}

type AuthData struct {
	AppName       string
	SystemId      int
	DomainId      int
	ServiceId     int
	ApplicationId int
}

type AuthorizeRequest struct {
	ApplicationId int    `validate:"required"`
	Endpoint      string `validate:"required"`
}

type AuthorizeResponse struct {
	Authorized bool
}
