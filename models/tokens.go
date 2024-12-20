package models

import (
	"time"

	"github.com/google/uuid"
)

type Tokens struct {
	ID        uuid.UUID `json:"id" gorm:"column:token_id;type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid"`
	Token     string    `json:"token" gorm:"not null"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
	IsValid   bool      `json:"is_valid" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	User      User      `json:"user" gorm:"foreignKey:UserID;references:ID"`
}
