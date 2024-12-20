package repositories

import (
	"invitified-go/models"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// Mock User Repository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Update(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}
func (m *MockUserRepository) FindRoleByID(id uuid.UUID) (*models.Role, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Role), args.Error(1)
}
func (m *MockUserRepository) FindByUsername(username string) (*models.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(id uuid.UUID) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindRoleByName(name string) (*models.Role, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Role), args.Error(1)
}

func (m *MockUserRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

// Mock Token Repository
type MockTokenRepository struct {
	mock.Mock
}

func (m *MockTokenRepository) SaveToken(token *models.Tokens) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *MockTokenRepository) FindToken(token string) (*models.Tokens, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Tokens), args.Error(1)
}
func (m *MockTokenRepository) InvalidateToken(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockTokenRepository) FindValidToken(userID uuid.UUID) (*models.Tokens, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Tokens), args.Error(1)
}
func (m *MockTokenRepository) FindByID(id uuid.UUID) (*models.Tokens, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Tokens), args.Error(1)
}
func (m *MockTokenRepository) UpdateToken(token *models.Tokens) error {
	args := m.Called(token)
	return args.Error(0)
}

// Mock Equipment Repository
type MockEquipmentRepository struct {
	mock.Mock
}

func (m *MockEquipmentRepository) CreateCategory(category *models.EquipmentCategory) error {
	args := m.Called(category)
	return args.Error(0)
}

func (m *MockEquipmentRepository) FindCategoryByID(id uuid.UUID) (*models.EquipmentCategory, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.EquipmentCategory), args.Error(1)
}

func (m *MockEquipmentRepository) FindAllCategories() ([]models.EquipmentCategory, error) {
	args := m.Called()
	return args.Get(0).([]models.EquipmentCategory), args.Error(1)
}

func (m *MockEquipmentRepository) UpdateCategory(category *models.EquipmentCategory) error {
	args := m.Called(category)
	return args.Error(0)
}

func (m *MockEquipmentRepository) DeleteCategory(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockEquipmentRepository) CreateEquipment(equipment *models.Equipment) error {
	args := m.Called(equipment)
	return args.Error(0)
}

func (m *MockEquipmentRepository) FindEquipmentByID(id uuid.UUID) (*models.Equipment, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Equipment), args.Error(1)
}

func (m *MockEquipmentRepository) FindEquipmentBySlug(slug string) (*models.Equipment, error) {
	args := m.Called(slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Equipment), args.Error(1)
}

func (m *MockEquipmentRepository) FindAllEquipmentWithPagination(limit, offset int) ([]models.Equipment, int64, error) {
	args := m.Called(limit, offset)
	return args.Get(0).([]models.Equipment), args.Get(1).(int64), args.Error(2)
}

func (m *MockEquipmentRepository) FindEquipmentByCategoryIDWithPagination(categoryID uuid.UUID, limit, offset int) ([]models.Equipment, int64, error) {
	args := m.Called(categoryID, limit, offset)
	return args.Get(0).([]models.Equipment), args.Get(1).(int64), args.Error(2)
}

func (m *MockEquipmentRepository) UpdateEquipment(equipment *models.Equipment) error {
	args := m.Called(equipment)
	return args.Error(0)
}

func (m *MockEquipmentRepository) DeleteEquipment(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}
func (m *MockEquipmentRepository) FindAllEquipment() ([]models.Equipment, error) {
	args := m.Called()
	return args.Get(0).([]models.Equipment), args.Error(1)
}

func (m *MockEquipmentRepository) FindEquipmentByCategoryID(categoryID uuid.UUID) ([]models.Equipment, error) {
	args := m.Called(categoryID)
	return args.Get(0).([]models.Equipment), args.Error(1)
}

func (m *MockEquipmentRepository) FindUserByID(id uuid.UUID) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

// Mock Rental Repository
type MockRentalRepository struct {
	mock.Mock
}

func (m *MockRentalRepository) Create(rental *models.Rental) error {
	args := m.Called(rental)
	return args.Error(0)
}

func (m *MockRentalRepository) FindByID(id uuid.UUID) (*models.Rental, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Rental), args.Error(1)
}

func (m *MockRentalRepository) FindAll() ([]models.Rental, error) {
	args := m.Called()
	return args.Get(0).([]models.Rental), args.Error(1)
}

func (m *MockRentalRepository) CheckOverlap(equipmentID uuid.UUID, startDate, endDate time.Time) (bool, error) {
	args := m.Called(equipmentID, startDate, endDate)
	return args.Get(0).(bool), args.Error(1)
}

func (m *MockRentalRepository) Update(rental *models.Rental) error {
	args := m.Called(rental)
	return args.Error(0)
}

func (m *MockRentalRepository) UpdateStatus(id uuid.UUID, status string) error {
	args := m.Called(id, status)
	return args.Error(0)
}

func (m *MockRentalRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

// Mock Payment Repository
type MockPaymentRepository struct {
	mock.Mock
}

func (m *MockPaymentRepository) Create(payment *models.Payment) error {
	args := m.Called(payment)
	return args.Error(0)
}

func (m *MockPaymentRepository) FindByID(id uuid.UUID) (*models.Payment, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Payment), args.Error(1)
}

func (m *MockPaymentRepository) FindByRentalID(rentalID uuid.UUID) (*models.Payment, error) {
	args := m.Called(rentalID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Payment), args.Error(1)
}

func (m *MockPaymentRepository) Update(payment *models.Payment) error {
	args := m.Called(payment)
	return args.Error(0)
}
func (m *MockPaymentRepository) FindByExternalID(externalID string) (*models.Payment, error) {
	args := m.Called(externalID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Payment), args.Error(1)
}
