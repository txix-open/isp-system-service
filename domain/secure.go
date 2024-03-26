package domain

type AuthenticateRequest struct {
	Token string `valid:"required~Required"`
}

type AuthenticateResponse struct {
	Authenticated bool
	ErrorReason   string
	AuthData      *AuthData
}

type AuthData struct {
	SystemId      int
	DomainId      int
	ServiceId     int
	ApplicationId int
}

type AuthorizeRequest struct {
	ApplicationId int    `valid:"required~Required"`
	Endpoint      string `valid:"required~Required"`
}

type AuthorizeResponse struct {
	Authorized bool
}
