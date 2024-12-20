package controllers

import (
	"invitified-go/models"
	"invitified-go/repositories"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// RentalController handles rental-related requests
type RentalController struct {
	repo          repositories.RentalRepository
	equipmentRepo repositories.EquipmentRepository
}

// NewRentalController creates a new RentalController
func NewRentalController(repo repositories.RentalRepository, equipmentRepo repositories.EquipmentRepository) *RentalController {
	return &RentalController{repo, equipmentRepo}
}

// CreateRental godoc
// @Summary Create a new rental
// @Description Create a new rental
// @Tags rentals
// @Accept json
// @Produce json
// @Param rental body models.Rental true "Rental"
// @Success 201 {object} models.Rental
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Security ApiKeyAuth
// @Router /rentals [post]
func (ctrl *RentalController) CreateRental(c echo.Context) error {
	rental := new(models.Rental)
	if err := c.Bind(rental); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	userIDStr, ok := c.Get("userID").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid user ID"})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid user ID"})
	}

	rental.UserID = userID

	// Calculate total cost and set equipment names in one loop
	var totalCost float64
	equipmentMap := make(map[uuid.UUID]*models.Equipment)

	for i := range rental.Items {
		equipment, err := ctrl.equipmentRepo.FindEquipmentByID(rental.Items[i].EquipmentID)
		if err != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"message": "Equipment not found"})
		}
		equipmentMap[equipment.ID] = equipment

		duration := rental.EndDate.Sub(rental.StartDate).Hours() / 24
		itemCost := float64(rental.Items[i].Quantity) * equipment.RentalPrice * duration
		totalCost += itemCost

		rental.Items[i].EquipmentName = equipment.Name
	}
	rental.TotalCost = totalCost

	// Create rental and rental items
	if err := ctrl.repo.Create(rental); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, rental)
}

// GetRentalByID godoc
// @Summary Get a rental by ID
// @Description Get a rental by ID
// @Tags rentals
// @Produce json
// @Param id path string true "Rental ID"
// @Success 200 {object} models.Rental
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Security ApiKeyAuth
// @Router /rentals/{id} [get]
func (ctrl *RentalController) GetRentalByID(c echo.Context) error {
	id := c.Param("id")
	rentalID, err := uuid.Parse(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	rental, err := ctrl.repo.FindByID(rentalID)
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}

	for i := range rental.Items {
		equipment, err := ctrl.equipmentRepo.FindEquipmentByID(rental.Items[i].EquipmentID)
		if err != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"message": "Equipment not found"})
		}
		rental.Items[i].EquipmentName = equipment.Name
	}

	return c.JSON(http.StatusOK, rental)
}

// GetAllRentals godoc
// @Summary Get all rentals
// @Description Get all rentals
// @Tags rentals
// @Produce json
// @Success 200 {array} models.Rental
// @Failure 500 {object} models.ErrorResponse
// @Security ApiKeyAuth
// @Router /rentals [get]
func (ctrl *RentalController) GetAllRentals(c echo.Context) error {
	rentals, err := ctrl.repo.FindAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	// Fetch equipment names for each rental item
	for i := range rentals {
		for j := range rentals[i].Items {
			equipment, err := ctrl.equipmentRepo.FindEquipmentByID(rentals[i].Items[j].EquipmentID)
			if err != nil {
				return c.JSON(http.StatusNotFound, map[string]string{"message": "Equipment not found"})
			}
			rentals[i].Items[j].EquipmentName = equipment.Name
		}
	}

	return c.JSON(http.StatusOK, rentals)
}

// UpdateRental godoc
// @Summary Update a rental
// @Description Update a rental
// @Tags rentals
// @Accept json
// @Produce json
// @Param id path string true "Rental ID"
// @Param rental body models.Rental true "Rental"
// @Success 200 {object} models.Rental
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Security ApiKeyAuth
// @Router /rentals/{id} [put]
func (ctrl *RentalController) UpdateRental(c echo.Context) error {
	id := c.Param("id")
	rentalID, err := uuid.Parse(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	rental, err := ctrl.repo.FindByID(rentalID)
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}
	if err := c.Bind(rental); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err := ctrl.repo.Update(rental); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, rental)
}

// DeleteRental godoc
// @Summary Delete a rental
// @Description Delete a rental
// @Tags rentals
// @Produce json
// @Param id path string true "Rental ID"
// @Success 204
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Security ApiKeyAuth
// @Router /rentals/{id} [delete]
func (ctrl *RentalController) DeleteRental(c echo.Context) error {
	id := c.Param("id")
	rentalID, err := uuid.Parse(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err := ctrl.repo.Delete(rentalID); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusNoContent)
}
