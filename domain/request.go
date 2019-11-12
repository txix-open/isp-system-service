package domain

type (
	Identity struct {
		Id int32 `json:"id" valid:"required~Required"`
	}

	RevokeTokensRequest struct {
		AppId  int32 `valid:"required~Required"`
		Tokens []string
	}

	CreateTokenRequest struct {
		AppId        int32 `valid:"required~Required"`
		ExpireTimeMs int64
	}
)
