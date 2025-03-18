package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type SaasUser struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Email     string    `gorm:"type:varchar(255);unique_index"`
	FirstName string    `gorm:"type:varchar(255)"`
	LastName  string    `gorm:"type:varchar(255)"`
	Password  string    `gorm:"type:varchar(255)"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
