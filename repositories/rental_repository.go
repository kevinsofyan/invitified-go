package repositories

import (
	"errors"
	"invitified-go/models"
	"os"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RentalRepository interface {
	Create(rental *models.Rental) error
	FindByID(id uuid.UUID) (*models.Rental, error)
	FindAll() ([]models.Rental, error)
	Update(rental *models.Rental) error
	Delete(id uuid.UUID) error
	CheckOverlap(equipmentID uuid.UUID, startDate, endDate time.Time) (bool, error)
}

type rentalRepository struct {
	db *gorm.DB
}

func NewRentalRepository(db *gorm.DB) RentalRepository {
	return &rentalRepository{db}
}

func (r *rentalRepository) Create(rental *models.Rental) error {
	// Check for overlapping rentals for each item
	for _, item := range rental.Items {
		overlap, err := r.CheckOverlap(item.EquipmentID, rental.StartDate, rental.EndDate)
		if err != nil {
			return err
		}
		if overlap {
			return errors.New("equipment is already rented for the specified time period")
		}
	}

	// Create rental and rental items
	if err := r.db.Create(rental).Error; err != nil {
		return err
	}
	for i := range rental.Items {
		rental.Items[i].RentalID = rental.ID
		rental.Items[i].ID = uuid.New() // Ensure unique rental_item_id
		if err := r.db.Create(&rental.Items[i]).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *rentalRepository) FindByID(id uuid.UUID) (*models.Rental, error) {
	var rental models.Rental
	err := r.db.Preload("Items").First(&rental, "rental_id = ?", id).Error
	return &rental, err
}

func (r *rentalRepository) FindAll() ([]models.Rental, error) {
	var rentals []models.Rental
	err := r.db.Preload("Items").Find(&rentals).Error
	return rentals, err
}

func (r *rentalRepository) Update(rental *models.Rental) error {
	return r.db.Save(rental).Error
}

func (r *rentalRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Rental{}, "rental_id = ?", id).Error
}

func (r *rentalRepository) CheckOverlap(equipmentID uuid.UUID, startDate, endDate time.Time) (bool, error) {
	var count int64
	schema := os.Getenv("DB_SCHEMA")
	err := r.db.Model(&models.Rental{}).
		Joins("JOIN \""+schema+"\".rental_items ON rentals.rental_id = rental_items.rental_id").
		Where("rental_items.equipment_id = ? AND rentals.status = 'PAID' AND rentals.end_date > ? AND rentals.start_date < ?", equipmentID, startDate, endDate).
		Count(&count).Error
	return count > 0, err
}
