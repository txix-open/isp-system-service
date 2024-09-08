package domain

type (
	Identity struct {
		Id int `validate:"required"`
	}

	DeleteResponse struct {
		Deleted int
	}
)
