package models

import (
	"time"

	"github.com/google/uuid"
)

type EquipmentCategory struct {
	ID          uuid.UUID `json:"id" gorm:"column:category_id;type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string    `json:"name" gorm:"unique;not null"`
	Slug        string    `json:"slug" gorm:"unique;not null"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
}

type Equipment struct {
	ID            uuid.UUID `json:"id" gorm:"column:equipment_id;type:uuid;primary_key;default:gen_random_uuid()"`
	Name          string    `json:"name" gorm:"not null"`
	Slug          string    `json:"slug" gorm:"unique;not null"`
	StockQuantity int       `json:"stock_quantity" gorm:"not null;default:0"`
	RentalPrice   float64   `json:"rental_price" gorm:"not null"`
	CategoryID    uuid.UUID `json:"category_id" gorm:"type:uuid"`
	IsAvailable   bool      `json:"is_available" gorm:"default:true"`
	CreatedBy     uuid.UUID `json:"created_by" gorm:"type:uuid"`
	CreatedAt     time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
}
