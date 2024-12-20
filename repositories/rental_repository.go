package repositories

import (
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
	UpdateStatus(id uuid.UUID, status string) error
}

type rentalRepository struct {
	db *gorm.DB
}

func NewRentalRepository(db *gorm.DB) RentalRepository {
	return &rentalRepository{db}
}

func (r *rentalRepository) Create(rental *models.Rental) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Generate IDs
		rental.ID = uuid.New()
		rental.Status = "PENDING"

		// Store equipment names before creating items
		equipmentNames := make(map[uuid.UUID]string)
		for _, item := range rental.Items {
			var equipment models.Equipment
			if err := tx.First(&equipment, "equipment_id = ?", item.EquipmentID).Error; err != nil {
				return err
			}
			equipmentNames[item.EquipmentID] = equipment.Name
		}

		// Create rental
		if err := tx.Create(&models.Rental{
			ID:        rental.ID,
			UserID:    rental.UserID,
			StartDate: rental.StartDate,
			EndDate:   rental.EndDate,
			TotalCost: rental.TotalCost,
			Status:    rental.Status,
		}).Error; err != nil {
			return err
		}

		// Create rental items
		items := make([]models.RentalItem, len(rental.Items))
		for i, item := range rental.Items {
			items[i] = models.RentalItem{
				ID:            uuid.New(),
				RentalID:      rental.ID,
				EquipmentID:   item.EquipmentID,
				Quantity:      item.Quantity,
				EquipmentName: equipmentNames[item.EquipmentID],
			}
		}

		// Batch create items
		if err := tx.Create(&items).Error; err != nil {
			return err
		}

		// Update rental object with created items
		rental.Items = items
		return nil
	})
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

func (r *rentalRepository) UpdateStatus(id uuid.UUID, status string) error {
	return r.db.Model(&models.Rental{}).Where("rental_id = ?", id).Update("status", status).Error
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
