package controllers

import (
	"bytes"
	"encoding/json"
	"invitified-go/models"
	"invitified-go/repositories"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRentalController(t *testing.T) {
	e := echo.New()
	mockRentalRepo := new(repositories.MockRentalRepository)
	mockEquipmentRepo := new(repositories.MockEquipmentRepository)
	ctrl := NewRentalController(mockRentalRepo, mockEquipmentRepo)

	t.Run("CreateRental", func(t *testing.T) {
		tests := []struct {
			name       string
			payload    models.Rental
			setupAuth  func(c echo.Context)
			setupMocks func()
			wantCode   int
			wantErr    bool
		}{
			{
				name: "successful rental creation",
				payload: models.Rental{
					StartDate: time.Now(),
					EndDate:   time.Now().Add(24 * time.Hour),
					Items: []models.RentalItem{
						{
							EquipmentID: uuid.New(),
							Quantity:    2,
						},
					},
				},
				setupAuth: func(c echo.Context) {
					userID := uuid.New()
					c.Set("userID", userID.String())
				},
				setupMocks: func() {
					equipment := &models.Equipment{
						ID:          uuid.New(),
						Name:        "Test Equipment",
						RentalPrice: 100.0,
					}
					mockEquipmentRepo.On("FindEquipmentByID", mock.AnythingOfType("uuid.UUID")).Return(equipment, nil)
					mockRentalRepo.On("Create", mock.AnythingOfType("*models.Rental")).Return(nil)
				},
				wantCode: http.StatusCreated,
				wantErr:  false,
			},
			{
				name: "invalid date range",
				payload: models.Rental{
					StartDate: time.Now().Add(24 * time.Hour),
					EndDate:   time.Now(),
					Items: []models.RentalItem{
						{
							EquipmentID: uuid.New(),
							Quantity:    1,
						},
					},
				},
				setupAuth: func(c echo.Context) {
					userID := uuid.New()
					c.Set("userID", userID.String())
				},
				setupMocks: func() {},
				wantCode:   http.StatusBadRequest,
				wantErr:    true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				mockRentalRepo.ExpectedCalls = nil
				mockEquipmentRepo.ExpectedCalls = nil

				tt.setupMocks()

				jsonBytes, _ := json.Marshal(tt.payload)
				req := httptest.NewRequest(http.MethodPost, "/rentals", bytes.NewBuffer(jsonBytes))
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)

				if tt.setupAuth != nil {
					tt.setupAuth(c)
				}

				err := ctrl.CreateRental(c)
				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
				assert.Equal(t, tt.wantCode, rec.Code)

				mockRentalRepo.AssertExpectations(t)
				mockEquipmentRepo.AssertExpectations(t)
			})
		}
	})

	t.Run("GetRentalByID", func(t *testing.T) {
		rental := &models.Rental{
			ID:        uuid.New(),
			StartDate: time.Now(),
			EndDate:   time.Now().Add(24 * time.Hour),
			Items: []models.RentalItem{
				{
					EquipmentID: uuid.New(),
					Quantity:    1,
				},
			},
		}

		tests := []struct {
			name       string
			rentalID   string
			setupMocks func()
			wantCode   int
		}{
			{
				name:     "successful rental retrieval",
				rentalID: rental.ID.String(),
				setupMocks: func() {
					mockRentalRepo.On("FindByID", rental.ID).Return(rental, nil)
					mockEquipmentRepo.On("FindEquipmentByID", mock.AnythingOfType("uuid.UUID")).Return(&models.Equipment{
						ID:   rental.Items[0].EquipmentID,
						Name: "Test Equipment",
					}, nil)
				},
				wantCode: http.StatusOK,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				mockRentalRepo.ExpectedCalls = nil
				mockEquipmentRepo.ExpectedCalls = nil

				tt.setupMocks()

				req := httptest.NewRequest(http.MethodGet, "/rentals/"+tt.rentalID, nil)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				c.SetParamNames("id")
				c.SetParamValues(tt.rentalID)

				err := ctrl.GetRentalByID(c)
				assert.NoError(t, err)
				assert.Equal(t, tt.wantCode, rec.Code)

				mockRentalRepo.AssertExpectations(t)
				mockEquipmentRepo.AssertExpectations(t)
			})
		}
	})
}
