package notification

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	Id        uuid.UUID `json:"id" gorm:"column:id;type:varchar(40);primaryKey"`
	UserId    string    `json:"user_id" gorm:"column:user_id;type:varchar(40);not null"`
	Title     string    `json:"title" gorm:"column:title;type:varchar(255)"`
	Message   string    `json:"message" gorm:"column:message;type:text"`
	ReadAt    string    `json:"read_at" gorm:"column:read_at;type:timestamptz"`
	CreatedAt time.Time `json:"-" gorm:"column:created_at;type:timestamptz"`
	UpdatedAt time.Time `json:"-" gorm:"column:updated_at;type:timestamptz"`
}
