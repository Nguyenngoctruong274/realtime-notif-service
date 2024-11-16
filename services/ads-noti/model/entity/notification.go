package entity

type Notification struct {
	ID        uint   `gorm:"primaryKey"`
	ProfileID string `gorm:"column:profile_id;type:varchar(50)"`
	ObjectID  string `gorm:"column:object_id;type:varchar(50)"`
	Type      string `gorm:"column:type;type:varchar(50)"`

	UserNotifications []UserNotification `gorm:"foreignKey:NotificationID"`
	Base
}
