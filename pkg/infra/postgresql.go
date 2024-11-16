package infra

import (
	"context"
	"github.com/spf13/cast"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
	"yes4all/ads-noti-api/pkg/config"
	logger_custom "yes4all/ads-noti-api/pkg/logger"
	"yes4all/ads-noti-api/services/ads-noti/model/entity"
)

// PGDB is wrapper of pg.DB
type PGDB struct {
	gormDB *gorm.DB
}

var pgSingleton *PGDB

func InitPostgresql() {
	dbConfiguration := config.PostgresConfig()
	// logger_custom := logger_custom.NewLogger() //nolint
	// https://github.com/go-gorm/postgres
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dbConfiguration.DBUrl,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		// logger_custom.Fatalf("Can't instance db connection error: %v", err)
		return
	}

	db.AutoMigrate(
		&entity.Notification{},
		&entity.UserNotification{},
	)

	sqlDB, err := db.DB()
	if err != nil {
		// logger_custom.Fatalf("Database setting have problem when set connection pool error: %v", err)
		return
	}

	if len(dbConfiguration.SetMaxIdleConns) > 0 {
		sqlDB.SetMaxIdleConns(cast.ToInt(dbConfiguration.SetMaxIdleConns))
	}

	if len(dbConfiguration.SetMaxOpenConns) > 0 {
		sqlDB.SetMaxOpenConns(cast.ToInt(dbConfiguration.SetMaxOpenConns))
	}

	if len(dbConfiguration.SetConnMaxLifetime) > 0 {
		sqlDB.SetConnMaxLifetime(
			cast.ToDuration(dbConfiguration.SetConnMaxLifetime) * time.Minute)
	}

	pgSingleton = &PGDB{db}
}

func ClosePostgresql() error {
	db, err := pgSingleton.gormDB.DB()
	db.Close()
	if err != nil {
		logger_custom.NewLogger().Fatalf("Close connection db error: %v", err)
	}
	return nil
}

func GetDB() *gorm.DB {
	if pgSingleton == nil {
		logger_custom.NewLogger().Fatal("Connection to database Postgres is not setup")
	}
	return pgSingleton.gormDB
}

// BeginTransaction start an Transaction, require defer ReleaseTransaction instantly
func BeginTransaction() (*gorm.DB, error) {
	tx := pgSingleton.gormDB.Begin()
	return tx, nil
}

func CommitTransaction(ctx context.Context, tx *gorm.DB) {
	tx.Commit()
}

func RollbackTransaction(ctx context.Context, tx *gorm.DB) {
	tx.Rollback()
}

func ReleaseTransaction(tx *gorm.DB, err error) {
	if err != nil {
		tx.Rollback()
	}
	tx.Commit()
}
