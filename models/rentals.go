package models

import (
	"time"

	"github.com/google/uuid"
)

type Rental struct {
	ID        uuid.UUID    `json:"id" gorm:"column:rental_id;type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID    `json:"user_id" gorm:"type:uuid;not null"`
	StartDate time.Time    `json:"start_date" gorm:"not null"`
	EndDate   time.Time    `json:"end_date" gorm:"not null"`
	TotalCost float64      `json:"total_cost" gorm:"not null"`
	Status    string       `json:"status" gorm:"type:varchar(20);default:'PENDING'"`
	Items     []RentalItem `json:"items" gorm:"foreignKey:RentalID"`
}

type RentalItem struct {
	ID            uuid.UUID `json:"id" gorm:"column:rental_item_id;type:uuid;primary_key;default:gen_random_uuid()"`
	RentalID      uuid.UUID `json:"rental_id" gorm:"type:uuid;not null"`
	EquipmentID   uuid.UUID `json:"equipment_id" gorm:"type:uuid;not null"`
	Quantity      int       `json:"quantity" gorm:"not null"`
	EquipmentName string    `json:"equipment_name" gorm:"-"`
}

const (
	RentalStatusPending  = "PENDING"
	RentalStatusPaid     = "PAID"
	RentalStatusComplete = "COMPLETED"
)
