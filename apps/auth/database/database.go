package database

import (
	"github.com/gopher-saas/gopher-saas/apps/auth/config"
	"github.com/gopher-saas/gopher-saas/apps/auth/models"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"moul.io/zapgorm2"
)

func NewDB(cfg config.PostgresConfig, logger *zap.Logger) (*gorm.DB, error) {
	dsn := "host=" + cfg.DatabaseHost + " user=" + cfg.DatabaseUser + " password=" + cfg.DatabasePassword + " dbname=" + cfg.DatabaseName + " port=" + cfg.DatabasePort
	ormLogger := zapgorm2.New(logger)
	ormLogger.SetAsDefault()
	ormLogger.IgnoreRecordNotFoundError = true
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: ormLogger,
	})
	if err != nil {
		return nil, err
	}
	logger.Info("Database connection established")
	err = db.AutoMigrate(&models.SaasUser{}, &models.SystemUser{})
	if err != nil {
		logger.Error("Error migrating models", zap.Error(err))
		return nil, err
	}

	return db, nil
}
