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

func TestRegisterUser(t *testing.T) {
	// Setup
	e := echo.New()
	mockUserRepo := new(repositories.MockUserRepository)
	mockTokenRepo := new(repositories.MockTokenRepository)
	ctrl := NewUserController(mockUserRepo, mockTokenRepo)

	// Test cases
	tests := []struct {
		name       string
		payload    models.UserRequest
		setupMocks func()
		wantCode   int
		wantMsg    string
	}{
		{
			name: "successful registration",
			payload: models.UserRequest{
				Username:      "kevinsofyan",
				Email:         "kevinsofyan.13@gmail.com",
				Password:      "password123",
				FullName:      "kevinsofyan",
				ContactNumber: "1234567890",
				RoleName:      "ADMIN",
			},
			setupMocks: func() {
				mockUserRepo.On("FindByEmail", "kevinsofyan.13@gmail.com").Return(nil, nil)
				mockUserRepo.On("FindRoleByName", "ADMIN").Return(&models.Role{
					ID:   uuid.New(),
					Name: "ADMIN",
				}, nil)
				mockUserRepo.On("Create", mock.AnythingOfType("*models.User")).Return(nil)
			},
			wantCode: http.StatusOK,
			wantMsg:  "User registered successfully",
		},
		{
			name: "email exists",
			payload: models.UserRequest{
				Username:      "kevinsofyan",
				Email:         "kevinsofyan.13@gmail.com",
				Password:      "password123",
				FullName:      "kevinsofyan",
				ContactNumber: "1234567890",
				RoleName:      "ADMIN",
			},
			setupMocks: func() {
				mockUserRepo.On("FindByEmail", "kevinsofyan.13@gmail.com").Return(&models.User{}, nil)
			},
			wantCode: http.StatusBadRequest,
			wantMsg:  "Email already exists",
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			jsonBytes, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest(http.MethodPost, "/users/register", bytes.NewBuffer(jsonBytes))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := ctrl.RegisterUser(c)
			assert.NoError(t, err)

			var response map[string]interface{}
			json.Unmarshal(rec.Body.Bytes(), &response)

			assert.Equal(t, tt.wantCode, rec.Code)
			if tt.wantCode == http.StatusOK {
				assert.Equal(t, tt.wantMsg, response["message"])
				assert.NotNil(t, response["user"])
			} else {
				assert.Equal(t, tt.wantMsg, response["message"])
			}
		})
	}
}
