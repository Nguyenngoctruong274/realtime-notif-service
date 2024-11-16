package entity

type UserNotification struct {
	ID             uint `gorm:"primaryKey"`
	NotificationID uint
	EmailBy        string       `gorm:"column:email_by;type:varchar(50)"`
	EmailTo        string       `gorm:"column:email_to;type:varchar(50)"`
	IsMarked       bool         `gorm:"column:is_marked;type:bool"`
	ActionType     string       `gorm:"column:action_type;type:varchar(50)"`
	Priority       string       `gorm:"column:priority;type:varchar(50)"`
	Content        string       `gorm:"column:content;type:text"`
	ImageUrl       string       `gorm:"column:image_url;type:text"`
	Notification   Notification `gorm:"foreignKey:NotificationID"`
	Description    string       `gorm:"column:description;type:text"`

	Base
}
