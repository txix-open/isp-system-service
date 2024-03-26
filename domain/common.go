package domain

type (
	Identity struct {
		Id int `valid:"required~Required"`
	}

	DeleteResponse struct {
		Deleted int
	}
)
