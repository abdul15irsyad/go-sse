package notification

import (
	"errors"
	"fmt"
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

func GetPaginatedNotifications(userId uuid.UUID, status string) ([]Notification, int64, error) {
	notifications := []Notification{}
	query := util.DB.Model(&Notification{})
	if status == "unread" {
		query = query.Where("read_at IS NULL")
	} else if status == "read" {
		query = query.Where("read_at IS NOT NULL")
	}
	if err := query.Where("user_id = ?", userId.String()).Order("created_at DESC").Find(&notifications).Error; err != nil {
		return nil, 0, err
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	return notifications, count, nil
}

func GetCountNotifications(userId uuid.UUID, status string) (int64, error) {
	query := util.DB.Model(&Notification{})
	if status == "unread" {
		query = query.Where("read_at IS NULL")
	} else if status == "read" {
		query = query.Where("read_at IS NOT NULL")
	}
	query = query.Where("user_id = ?", userId.String()).Order("created_at DESC")
	var count int64
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func ReadNotification(id uuid.UUID, userId uuid.UUID) error {
	var notification Notification
	if err := util.DB.Model(&Notification{}).Where("id = ?", id).Where("user_id = ?", userId).First(&notification).Error; err != nil {
		return err
	}

	if notification.ReadAt != nil {
		return errors.New("notification already read")
	}

	util.DB.Model(&notification).Update("read_at", time.Now())
	fmt.Println("notification read")

	return nil
}
