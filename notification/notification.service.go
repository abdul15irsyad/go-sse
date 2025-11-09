package notification

import (
	"go-sse/util"

	"github.com/google/uuid"
)

func GetPaginatedNotifications(userId uuid.UUID, dto GetNotificationsDto) ([]Notification, int64, error) {
	notifications := []Notification{}
	offset := (dto.Page - 1) * dto.Limit
	query := util.DB.Model(&Notification{})
	if err := query.Where("user_id = ?", userId.String()).Limit(dto.Limit).Offset(offset).Order("created_at DESC").Find(&notifications).Error; err != nil {
		return nil, 0, err
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	return notifications, count, nil
}
