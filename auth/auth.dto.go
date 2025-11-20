package auth

type (
	LoginDTO struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	RegisterDTO struct {
		Name            string `json:"name" validate:"required"`
		Username        string `json:"username" validate:"required"`
		Password        string `json:"password" validate:"required"`
		ConfirmPassword string `json:"confirmPassword" validate:"required"`
	}
)
