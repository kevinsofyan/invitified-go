package controllers

import (
	"invitified-go/models"
	"invitified-go/repositories"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type RentalController struct {
	repo          repositories.RentalRepository
	equipmentRepo repositories.EquipmentRepository
}

func NewRentalController(repo repositories.RentalRepository, equipmentRepo repositories.EquipmentRepository) *RentalController {
	return &RentalController{repo, equipmentRepo}
}

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
