package notification

import "github.com/google/uuid"

type (
	CreateNotificationDTO struct {
		UserId  uuid.UUID
		Title   string
		Message string
	}

	GetNotificationsDto struct {
		Page  int `form:"page" validate:"required,number,gte=1"`
		Limit int `form:"limit" validate:"required,number"`
	}
)
