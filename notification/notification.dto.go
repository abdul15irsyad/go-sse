package notification

type (
	GetNotificationsDto struct {
		Page  int `form:"page" validate:"required,number,gte=1"`
		Limit int `form:"limit" validate:"required,number"`
	}
)
