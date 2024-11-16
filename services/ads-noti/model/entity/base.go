package entity

import "time"

const (
	NotificationDBName     string = "notifications"
	UserNotificationDBName string = "user_notifications"
)

type Base struct {
	CreatedAt *time.Time `json:"createdAt" gorm:"column:created_at;type:TIMESTAMP WITH TIME ZONE"`
	UpdatedAt *time.Time `json:"updatedAt" gorm:"column:updated_at;type:TIMESTAMP WITH TIME ZONE"`
	DeletedAt *time.Time `json:"deletedAt" gorm:"column:deleted_at;type:TIMESTAMP WITH TIME ZONE"`
	CreatedBy *string    `json:"createdBy" gorm:"column:created_by;type:varchar(100)"`
	UpdatedBy *string    `json:"updatedBy" gorm:"column:updated_by;type:varchar(100)"`
	DeletedBy *string    `json:"deletedBy" gorm:"column:deleted_by;type:varchar(100)"`
}
