package notification

import (
	"go-sse/util"
	"time"

	"github.com/google/uuid"
)

func CreateNotification(dto CreateNotificationDTO) (Notification, error) {
	newUuid, _ := uuid.NewRandom()

	newNotification := Notification{
		Id:        newUuid,
		UserId:    dto.UserId,
		Title:     dto.Title,
		Message:   dto.Message,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := util.DB.Model(&Notification{}).Create(newNotification).Error; err != nil {
		return Notification{}, err
	}

	return newNotification, nil
}

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

func ReadNotification(id uuid.UUID) error {
	var notification Notification
	if err := util.DB.Model(&Notification{}).Where("id = ?", id).First(&notification).Error; err != nil {
		return err
	}

	if notification.ReadAt != nil {
		util.DB.Model(&notification).Update("read_at", time.Now())
	}

	return nil
}
