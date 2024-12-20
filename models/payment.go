package models

import (
	"time"

	"github.com/google/uuid"
)

type Payment struct {
	ID                   uuid.UUID  `json:"id" gorm:"column:payment_id;type:uuid;primary_key;default:gen_random_uuid()"`
	RentalID             uuid.UUID  `json:"rental_id" gorm:"type:uuid;not null"`
	UserID               uuid.UUID  `json:"user_id" gorm:"type:uuid;not null"`
	Amount               float64    `json:"amount" gorm:"not null"`
	PointsUsed           int        `json:"points_used" gorm:"default:0"`
	PointsEarned         int        `json:"points_earned" gorm:"default:0"`
	PaymentMethod        string     `json:"payment_method" gorm:"type:varchar(50);not null"`
	PaymentStatus        string     `json:"payment_status" gorm:"type:varchar(20);default:'PENDING'"`
	XenditInvoiceID      string     `json:"xendit_invoice_id" gorm:"type:varchar(100)"`
	XenditPaymentURL     string     `json:"xendit_payment_url" gorm:"type:varchar(255)"`
	XenditPaymentChannel string     `json:"xendit_payment_channel" gorm:"type:varchar(50)"`
	XenditPaidAmount     float64    `json:"xendit_paid_amount"`
	PaidAt               *time.Time `json:"paid_at"`
	ExpiredAt            *time.Time `json:"expired_at"`
	CreatedAt            time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt            time.Time  `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}
