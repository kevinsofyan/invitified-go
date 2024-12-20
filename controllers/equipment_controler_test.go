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

func TestPaymentController_CreatePayment(t *testing.T) {
	e := echo.New()
	mockPaymentRepo := new(repositories.MockPaymentRepository)
	mockRentalRepo := new(repositories.MockRentalRepository)
	mockUserRepo := new(repositories.MockUserRepository)

	ctrl := NewPaymentController(mockPaymentRepo, mockRentalRepo, mockUserRepo)

	tests := []struct {
		name       string
		payload    PaymentRequest
		setupAuth  func(c echo.Context)
		setupMocks func()
		wantCode   int
		wantMsg    string
	}{
		{
			name: "successful payment",
			payload: PaymentRequest{
				RentalID:      uuid.New().String(),
				PaymentMethod: "VIRTUAL_ACCOUNT",
				ChannelCode:   "BCA",
			},
			setupAuth: func(c echo.Context) {
				c.Set("userID", uuid.New().String())
			},
			setupMocks: func() {
				rental := &models.Rental{
					ID:        uuid.MustParse(uuid.New().String()),
					UserID:    uuid.MustParse(uuid.New().String()),
					TotalCost: 100000,
					Status:    models.RentalStatusPending,
				}

				mockRentalRepo.On("FindByID", mock.AnythingOfType("uuid.UUID")).Return(rental, nil)
				mockUserRepo.On("FindByID", mock.AnythingOfType("uuid.UUID")).Return(&models.User{
					ID:    rental.UserID,
					Email: "test@example.com",
				}, nil)
				mockPaymentRepo.On("Create", mock.AnythingOfType("*models.Payment")).Return(nil)
				mockRentalRepo.On("UpdateStatus", rental.ID, models.RentalStatusComplete).Return(nil)
			},
			wantCode: http.StatusCreated,
		},
		{
			name: "invalid rental id",
			payload: PaymentRequest{
				RentalID:      "invalid-uuid",
				PaymentMethod: "VIRTUAL_ACCOUNT",
				ChannelCode:   "BCA",
			},
			setupAuth: func(c echo.Context) {
				c.Set("userID", uuid.New().String())
			},
			setupMocks: func() {},
			wantCode:   http.StatusBadRequest,
			wantMsg:    "Invalid rental ID format",
		},
		{
			name: "unauthorized access",
			payload: PaymentRequest{
				RentalID:      uuid.New().String(),
				PaymentMethod: "VIRTUAL_ACCOUNT",
				ChannelCode:   "BCA",
			},
			setupAuth: func(c echo.Context) {
				c.Set("userID", uuid.New().String())
			},
			setupMocks: func() {
				rental := &models.Rental{
					ID:        uuid.MustParse(uuid.New().String()),
					UserID:    uuid.New(), // Different user ID
					TotalCost: 100000,
					Status:    models.RentalStatusPending,
				}
				mockRentalRepo.On("FindByID", mock.AnythingOfType("uuid.UUID")).Return(rental, nil)
				mockUserRepo.On("FindByID", mock.AnythingOfType("uuid.UUID")).Return(&models.User{}, nil)
			},
			wantCode: http.StatusForbidden,
			wantMsg:  "Not authorized to pay for this rental",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks
			mockPaymentRepo.ExpectedCalls = nil
			mockRentalRepo.ExpectedCalls = nil
			mockUserRepo.ExpectedCalls = nil

			tt.setupMocks()

			jsonBytes, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest(http.MethodPost, "/payments", bytes.NewBuffer(jsonBytes))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if tt.setupAuth != nil {
				tt.setupAuth(c)
			}

			err := ctrl.CreatePayment(c)
			assert.NoError(t, err)

			if tt.wantMsg != "" {
				var response map[string]string
				json.Unmarshal(rec.Body.Bytes(), &response)
				assert.Equal(t, tt.wantMsg, response["message"])
			}

			assert.Equal(t, tt.wantCode, rec.Code)

			// Verify mock expectations
			mockPaymentRepo.AssertExpectations(t)
			mockRentalRepo.AssertExpectations(t)
			mockUserRepo.AssertExpectations(t)
		})
	}
}
