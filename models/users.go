package models

import (
	"time"

	"github.com/google/uuid"
)

type Role struct {
	ID          uuid.UUID `json:"id" gorm:"column:role_id;type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string    `json:"name" gorm:"column:role_name;type:user_role;unique;not null"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
}

type User struct {
	ID            uuid.UUID `json:"id" gorm:"column:user_id;type:uuid;primary_key;default:gen_random_uuid()"`
	Username      string    `json:"username" gorm:"unique;not null"`
	Email         string    `json:"email" gorm:"unique;not null"`
	Password      string    `json:"-" gorm:"not null"`
	FullName      string    `json:"full_name" gorm:"not null"`
	ContactNumber string    `json:"contact_number"`
	RoleID        uuid.UUID `json:"-" gorm:"column:role_id"`
	RoleName      string    `json:"role_name" gorm:"-"`
	CreatedAt     time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
}
