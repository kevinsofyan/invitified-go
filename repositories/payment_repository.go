package repositories

import (
	"invitified-go/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentRepository interface {
	Create(payment *models.Payment) error
	FindByID(id uuid.UUID) (*models.Payment, error)
	FindByExternalID(externalID string) (*models.Payment, error)
	Update(payment *models.Payment) error
}

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db}
}

func (r *paymentRepository) Create(payment *models.Payment) error {
	return r.db.Create(payment).Error
}

func (r *paymentRepository) FindByID(id uuid.UUID) (*models.Payment, error) {
	var payment models.Payment
	if err := r.db.First(&payment, "payment_id = ?", id).Error; err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) FindByExternalID(externalID string) (*models.Payment, error) {
	var payment models.Payment
	if err := r.db.First(&payment, "xendit_invoice_id = ?", externalID).Error; err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) Update(payment *models.Payment) error {
	return r.db.Save(payment).Error
}
