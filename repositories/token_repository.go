package repositories

import (
	"invitified-go/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TokenRepository interface {
	SaveToken(token *models.Tokens) error
	UpdateToken(token *models.Tokens) error
	FindToken(token string) (*models.Tokens, error)
	FindByID(id uuid.UUID) (*models.Tokens, error)
	FindValidToken(userID uuid.UUID) (*models.Tokens, error)
	InvalidateToken(id uuid.UUID) error
}

type tokenRepository struct {
	db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) TokenRepository {
	return &tokenRepository{db}
}

func (r *tokenRepository) SaveToken(token *models.Tokens) error {
	return r.db.Create(token).Error
}

func (r *tokenRepository) UpdateToken(token *models.Tokens) error {
	return r.db.Save(token).Error
}

func (r *tokenRepository) FindToken(token string) (*models.Tokens, error) {
	var tokenModel models.Tokens
	err := r.db.Where("token = ? AND is_valid = ? AND expires_at > NOW()", token, true).First(&tokenModel).Error
	return &tokenModel, err
}

func (r *tokenRepository) FindByID(id uuid.UUID) (*models.Tokens, error) {
	var token models.Tokens
	err := r.db.First(&token, "token_id = ?", id).Error
	return &token, err
}

func (r *tokenRepository) FindValidToken(userID uuid.UUID) (*models.Tokens, error) {
	var token models.Tokens
	err := r.db.Where("user_id = ? AND is_valid = ? AND expires_at > NOW()", userID, true).First(&token).Error
	return &token, err
}

func (r *tokenRepository) InvalidateToken(id uuid.UUID) error {
	return r.db.Model(&models.Tokens{}).Where("token_id = ?", id).Update("is_valid", false).Error
}
