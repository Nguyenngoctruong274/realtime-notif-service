package repository

import (
	"context"
	"gorm.io/gorm"
	"yes4all/ads-noti-api/pkg/infra"
	"yes4all/ads-noti-api/pkg/logger"
	"yes4all/ads-noti-api/pkg/xcontext"
	"yes4all/ads-noti-api/services/ads-noti/model/entity"
)

type NotificationRepo interface {
	GetByID(ctx context.Context, objectID string) (result entity.Notification, err error)
	GetList(ctx context.Context) (items []entity.Notification, err error)
	GetListByUser(ctx context.Context, email string) (items []entity.UserNotification, err error)
	Create(ctx context.Context, item *entity.Notification) (err error)
	UpdateUserNotification(ctx context.Context, item []entity.UserNotification) (err error)
}

type notificationRepo struct {
	db *gorm.DB
}

func NewNotificationRepo() NotificationRepo {
	return &notificationRepo{
		db: infra.GetDB(),
	}
}
func (r *notificationRepo) tableNotification() *gorm.DB {
	return r.db.Table(entity.NotificationDBName)
}
func (r *notificationRepo) tableUserNotification() *gorm.DB {
	return r.db.Table(entity.UserNotificationDBName)
}

func (r *notificationRepo) GetList(ctx context.Context) (resp []entity.Notification, err error) {
	entry := logger.NewLogger().WithKeyword(ctx, "GetList")
	profileID, _ := xcontext.GetProfileID(ctx)

	db := r.tableNotification()
	if len(profileID) > 0 {
		db = db.Where("profile_id in ?", profileID)
	}

	err = db.
		Preload("UserNotification").
		Find(&resp).Error
	entry.WithOutputField("resp", resp)
	if err != nil {
		entry.WithError(err).Error()
		return
	}
	entry.Info()
	return
}

func (r *notificationRepo) GetListByUser(ctx context.Context, email string) (resp []entity.UserNotification, err error) {
	entry := logger.NewLogger().WithKeyword(ctx, "GetListByUser")
	profileID, _ := xcontext.GetProfileID(ctx)

	db := r.tableUserNotification()
	db = db.Where("email_to = ?", email)
	if len(profileID) > 0 {
		db = db.Where("profile_id in ?", profileID)
	}

	err = db.
		Preload("Notification").Order("created_at DESC").Find(&resp).Error
	entry.WithOutputField("resp", resp)
	if err != nil {
		entry.WithError(err).Error()
		return
	}
	entry.Info()
	return
}

func (r *notificationRepo) GetByID(ctx context.Context, objectID string) (res entity.Notification, err error) {
	entry := logger.NewLogger().WithKeyword(ctx, "GetByID").WithOutputField("objectID", objectID)

	err = r.tableNotification().Where("object_id = ?", objectID).First(&res).Error
	if err != nil {
		entry.WithError(err).Error()
		return
	}

	entry.Info()
	return
}

func (r *notificationRepo) Create(ctx context.Context, item *entity.Notification) (err error) {
	entry := logger.NewLogger().WithKeyword(ctx, "Create").WithOutputField("item", item)
	err = r.tableNotification().Create(&item).Error
	if err != nil {
		entry.WithError(err).Error()
		return
	}
	entry.Info()
	return
}

func (r *notificationRepo) UpdateUserNotification(ctx context.Context, items []entity.UserNotification) (err error) {
	entry := logger.NewLogger().WithKeyword(ctx, "Updates").
		WithField("items", items)

	if len(items) == 0 {
		entry.Info()
		return
	}
	for _, v := range items {
		err = r.tableUserNotification().
			Select("is_marked").
			Updates(v).
			Error
		if err != nil {
			entry.WithError(err).Error()
			return
		}
	}
	entry.Info()
	return
}
