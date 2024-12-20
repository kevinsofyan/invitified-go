package controllers

import (
	"bytes"
	"encoding/json"
	"invitified-go/models"
	"invitified-go/repositories"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestEquipmentController_Categories(t *testing.T) {
	e := echo.New()
	mockRepo := new(repositories.MockEquipmentRepository)
	ctrl := NewEquipmentController(mockRepo)

	t.Run("CreateCategory", func(t *testing.T) {
		tests := []struct {
			name       string
			payload    models.EquipmentCategory
			setupMocks func()
			wantCode   int
		}{
			{
				name: "successful category creation",
				payload: models.EquipmentCategory{
					Name:        "Outdoor Equipment",
					Description: "Equipment for outdoor activities",
				},
				setupMocks: func() {
					mockRepo.On("CreateCategory", mock.AnythingOfType("*models.EquipmentCategory")).Return(nil)
				},
				wantCode: http.StatusCreated,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				mockRepo.ExpectedCalls = nil
				tt.setupMocks()

				jsonBytes, _ := json.Marshal(tt.payload)
				req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewBuffer(jsonBytes))
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)

				assert.NoError(t, ctrl.CreateCategory(c))
				assert.Equal(t, tt.wantCode, rec.Code)
			})
		}
	})
}

func TestEquipmentController_Equipment(t *testing.T) {
	e := echo.New()
	mockRepo := new(repositories.MockEquipmentRepository)
	ctrl := NewEquipmentController(mockRepo)

	t.Run("CreateEquipment", func(t *testing.T) {
		tests := []struct {
			name       string
			payload    models.Equipment
			setupMocks func()
			wantCode   int
		}{
			{
				name: "successful equipment creation",
				payload: models.Equipment{
					Name:          "Camping Tent",
					StockQuantity: 5,
					RentalPrice:   50.00,
					CategoryID:    uuid.New(),
					IsAvailable:   true,
				},
				setupMocks: func() {
					mockRepo.On("CreateEquipment", mock.AnythingOfType("*models.Equipment")).Return(nil)
				},
				wantCode: http.StatusCreated,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				mockRepo.ExpectedCalls = nil
				tt.setupMocks()

				jsonBytes, _ := json.Marshal(tt.payload)
				req := httptest.NewRequest(http.MethodPost, "/equipment", bytes.NewBuffer(jsonBytes))
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				c.Set("userID", uuid.New().String())

				assert.NoError(t, ctrl.CreateEquipment(c))
				assert.Equal(t, tt.wantCode, rec.Code)
			})
		}
	})
}
