package models

import (
	"github.com/google/uuid"
	"github.com/gopher-saas/gopher-saas/apps/auth/models/enum"
)

type SystemUser struct {
	ID        uuid.UUID       `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Email     string          `gorm:"type:varchar(255);unique_index"`
	FirstName string          `gorm:"type:varchar(255)"`
	LastName  string          `gorm:"type:varchar(255)"`
	Password  string          `gorm:"type:varchar(255)"`
	Role      enum.SystemRole `gorm:"type:system_role;not null"`
}
