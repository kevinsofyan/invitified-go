package controllers

import (
	"invitified-go/models"
	"invitified-go/repositories"
	"invitified-go/utils"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	repo      repositories.UserRepository
	tokenRepo repositories.TokenRepository
}

func NewUserController(repo repositories.UserRepository, tokenRepo repositories.TokenRepository) *UserController {
	return &UserController{repo, tokenRepo}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// ErrorResponse represents the error response structure
type ErrorResponse struct {
	Message string `json:"message"`
}

// RegisterUser registers a new user
func (ctrl *UserController) RegisterUser(c echo.Context) error {
	type UserRequest struct {
		Username      string `json:"username"`
		Email         string `json:"email"`
		Password      string `json:"password"`
		FullName      string `json:"full_name"`
		ContactNumber string `json:"contact_number"`
		RoleName      string `json:"role_name"`
	}

	var req UserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid input"})
	}

	// Check if the email already exists
	existingUser, err := ctrl.repo.FindByEmail(req.Email)
	if err == nil && existingUser != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Email already exists"})
	}

	// Find role by role_name
	role, err := ctrl.repo.FindRoleByName(req.RoleName)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid role")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Failed to hash password"})
	}

	user := &models.User{
		Username:      req.Username,
		Email:         req.Email,
		Password:      string(hashedPassword),
		FullName:      req.FullName,
		ContactNumber: req.ContactNumber,
		RoleID:        role.ID,
		RoleName:      role.Name,
	}

	if err := ctrl.repo.Create(user); err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Failed to register user"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "User registered successfully",
		"user":    user,
	})
}

func (ctrl *UserController) LoginUser(c echo.Context) error {
	loginRequest := new(LoginRequest)
	if err := c.Bind(loginRequest); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid input"})
	}

	user, err := ctrl.repo.FindByEmail(loginRequest.Email)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, ErrorResponse{Message: "Invalid email or password"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, ErrorResponse{Message: "Invalid email or password"})
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Failed to generate token"})
	}

	// Check if a token already exists for the user
	existingToken, err := ctrl.tokenRepo.FindValidToken(user.ID)
	if err == nil && existingToken != nil {
		// Update the existing token
		existingToken.Token = token
		existingToken.ExpiresAt = time.Now().Add(time.Hour * 24) // Token valid for 24 hours
		existingToken.IsValid = true
		if err := ctrl.tokenRepo.UpdateToken(existingToken); err != nil {
			return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Failed to update token"})
		}
	} else {
		// Create a new token
		tokenModel := &models.Tokens{
			UserID:    user.ID,
			Token:     token,
			ExpiresAt: time.Now().Add(time.Hour * 24), // Token valid for 24 hours
			IsValid:   true,
			CreatedAt: time.Now(),
		}
		if err := ctrl.tokenRepo.SaveToken(tokenModel); err != nil {
			return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Failed to save token"})
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Login successful",
		"token":   token,
	})
}

// GetUserProfile gets the profile of the logged-in user
func (ctrl *UserController) GetUserProfile(c echo.Context) error {
	userID, ok := c.Get("userID").(uuid.UUID)
	if !ok {
		return c.JSON(http.StatusUnauthorized, ErrorResponse{Message: "Invalid user ID"})
	}

	user, err := ctrl.repo.FindByID(userID)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, ErrorResponse{Message: "User not found"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"user": user,
	})
}

func (ctrl *UserController) DeleteUser(c echo.Context) error {
	loggedInUserID, ok := c.Get("userID").(uuid.UUID)
	if !ok {
		return c.JSON(http.StatusUnauthorized, ErrorResponse{Message: "Invalid user ID"})
	}

	loggedInUser, err := ctrl.repo.FindByID(loggedInUserID)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, ErrorResponse{Message: "User not found"})
	}

	if loggedInUser.RoleName != "ADMIN" {
		return c.JSON(http.StatusForbidden, ErrorResponse{Message: "Only admins can delete users"})
	}

	userID := c.Param("id")
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid user ID"})
	}

	if err := ctrl.repo.Delete(userUUID); err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Failed to delete user"})
	}

	return c.NoContent(http.StatusNoContent)
}
