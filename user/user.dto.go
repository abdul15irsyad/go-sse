package user

type (
	CreateUserDTO struct {
		Name     string
		Username *string
		Password *string
	}

	GetUserDto struct {
		Id string `validate:"required,uuid"`
	}

	GetUsersDto struct {
		Page   int     `form:"page" validate:"required,number,gte=1"`
		Limit  int     `form:"limit" validate:"required,number"`
		Search *string `form:"search" validate:"omitempty"`
	}
)
