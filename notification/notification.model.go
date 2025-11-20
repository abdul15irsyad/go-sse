package notification

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	Id        uuid.UUID  `json:"id" gorm:"column:id;type:varchar(40);primaryKey"`
	UserId    uuid.UUID  `json:"userId" gorm:"column:user_id;type:varchar(40);not null"`
	Title     string     `json:"title" gorm:"column:title;type:varchar(255)"`
	Message   string     `json:"message" gorm:"column:message;type:text"`
	ReadAt    *time.Time `json:"readAt" gorm:"column:read_at;type:timestamptz"`
	CreatedAt time.Time  `json:"createdAt" gorm:"column:created_at;type:timestamptz"`
	UpdatedAt time.Time  `json:"updatedAt" gorm:"column:updated_at;type:timestamptz"`
}
