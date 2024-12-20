package repositories

import (
	"invitified-go/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EquipmentRepository interface {
	CreateCategory(category *models.EquipmentCategory) error
	FindCategoryByID(id uuid.UUID) (*models.EquipmentCategory, error)
	FindAllCategories() ([]models.EquipmentCategory, error)
	UpdateCategory(category *models.EquipmentCategory) error
	DeleteCategory(id uuid.UUID) error

	CreateEquipment(equipment *models.Equipment) error
	FindEquipmentByID(id uuid.UUID) (*models.Equipment, error)
	FindEquipmentBySlug(slug string) (*models.Equipment, error)
	FindAllEquipment() ([]models.Equipment, error)
	FindEquipmentByCategoryID(categoryID uuid.UUID) ([]models.Equipment, error)
	FindAllEquipmentWithPagination(limit, offset int) ([]models.Equipment, int64, error)
	FindEquipmentByCategoryIDWithPagination(categoryID uuid.UUID, limit, offset int) ([]models.Equipment, int64, error)
	UpdateEquipment(equipment *models.Equipment) error
	DeleteEquipment(id uuid.UUID) error

	FindUserByID(id uuid.UUID) (*models.User, error)
}

type equipmentRepository struct {
	db *gorm.DB
}

func NewEquipmentRepository(db *gorm.DB) EquipmentRepository {
	return &equipmentRepository{db}
}

func (r *equipmentRepository) CreateCategory(category *models.EquipmentCategory) error {
	return r.db.Create(category).Error
}

func (r *equipmentRepository) FindCategoryByID(id uuid.UUID) (*models.EquipmentCategory, error) {
	var category models.EquipmentCategory
	err := r.db.First(&category, "category_id = ?", id).Error
	return &category, err
}

func (r *equipmentRepository) FindAllCategories() ([]models.EquipmentCategory, error) {
	var categories []models.EquipmentCategory
	err := r.db.Find(&categories).Error
	return categories, err
}

func (r *equipmentRepository) UpdateCategory(category *models.EquipmentCategory) error {
	return r.db.Save(category).Error
}

func (r *equipmentRepository) DeleteCategory(id uuid.UUID) error {
	return r.db.Delete(&models.EquipmentCategory{}, "category_id = ?", id).Error
}

func (r *equipmentRepository) CreateEquipment(equipment *models.Equipment) error {
	return r.db.Create(equipment).Error
}

func (r *equipmentRepository) FindEquipmentByID(id uuid.UUID) (*models.Equipment, error) {
	var equipment models.Equipment
	err := r.db.First(&equipment, "equipment_id = ?", id).Error
	return &equipment, err
}

func (r *equipmentRepository) FindEquipmentBySlug(slug string) (*models.Equipment, error) {
	var equipment models.Equipment
	err := r.db.First(&equipment, "slug = ?", slug).Error
	return &equipment, err
}

func (r *equipmentRepository) FindAllEquipment() ([]models.Equipment, error) {
	var equipment []models.Equipment
	err := r.db.Find(&equipment).Error
	return equipment, err
}

func (r *equipmentRepository) FindEquipmentByCategoryID(categoryID uuid.UUID) ([]models.Equipment, error) {
	var equipment []models.Equipment
	err := r.db.Where("category_id = ?", categoryID).Find(&equipment).Error
	return equipment, err
}

func (r *equipmentRepository) FindAllEquipmentWithPagination(limit, offset int) ([]models.Equipment, int64, error) {
	var equipment []models.Equipment
	var total int64
	err := r.db.Model(&models.Equipment{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = r.db.Limit(limit).Offset(offset).Find(&equipment).Error
	return equipment, total, err
}

func (r *equipmentRepository) FindEquipmentByCategoryIDWithPagination(categoryID uuid.UUID, limit, offset int) ([]models.Equipment, int64, error) {
	var equipment []models.Equipment
	var total int64
	err := r.db.Model(&models.Equipment{}).Where("category_id = ?", categoryID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = r.db.Where("category_id = ?", categoryID).Limit(limit).Offset(offset).Find(&equipment).Error
	return equipment, total, err
}

func (r *equipmentRepository) UpdateEquipment(equipment *models.Equipment) error {
	return r.db.Save(equipment).Error
}

func (r *equipmentRepository) DeleteEquipment(id uuid.UUID) error {
	return r.db.Delete(&models.Equipment{}, "equipment_id = ?", id).Error
}

func (r *equipmentRepository) FindUserByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, "user_id = ?", id).Error
	return &user, err
}
