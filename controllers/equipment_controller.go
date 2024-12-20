package controllers

import (
	"invitified-go/models"
	"invitified-go/repositories"
	"invitified-go/utils"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type EquipmentController struct {
	repo repositories.EquipmentRepository
}

func NewEquipmentController(repo repositories.EquipmentRepository) *EquipmentController {
	return &EquipmentController{repo}
}

func (ctrl *EquipmentController) CreateCategory(c echo.Context) error {
	category := new(models.EquipmentCategory)
	if err := c.Bind(category); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	category.Slug = utils.ConvertToSlug(category.Name)
	if err := ctrl.repo.CreateCategory(category); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusCreated, category)
}

func (ctrl *EquipmentController) GetCategoryByID(c echo.Context) error {
	id := c.Param("id")
	categoryID, err := uuid.Parse(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	category, err := ctrl.repo.FindCategoryByID(categoryID)
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}
	return c.JSON(http.StatusOK, category)
}

func (ctrl *EquipmentController) GetAllCategories(c echo.Context) error {
	categories, err := ctrl.repo.FindAllCategories()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, categories)
}

func (ctrl *EquipmentController) UpdateCategory(c echo.Context) error {
	id := c.Param("id")
	categoryID, err := uuid.Parse(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	category := new(models.EquipmentCategory)
	if err := c.Bind(category); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	category.ID = categoryID
	category.Slug = utils.ConvertToSlug(category.Name)
	if err := ctrl.repo.UpdateCategory(category); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, category)
}

func (ctrl *EquipmentController) DeleteCategory(c echo.Context) error {
	id := c.Param("id")
	categoryID, err := uuid.Parse(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err := ctrl.repo.DeleteCategory(categoryID); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusNoContent)
}

// Equipment Handlers

func (ctrl *EquipmentController) CreateEquipment(c echo.Context) error {
	equipment := new(models.Equipment)
	if err := c.Bind(equipment); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	equipment.Slug = utils.ConvertToSlug(equipment.Name)

	userIDStr, ok := c.Get("userID").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid user ID"})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid user ID"})
	}

	equipment.CreatedBy = userID

	if err := ctrl.repo.CreateEquipment(equipment); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusCreated, equipment)
}

func (ctrl *EquipmentController) GetEquipmentBySlug(c echo.Context) error {
	slug := c.Param("slug")
	equipment, err := ctrl.repo.FindEquipmentBySlug(slug)
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}
	return c.JSON(http.StatusOK, equipment)
}

func (ctrl *EquipmentController) GetAllEquipment(c echo.Context) error {
	categoryID := c.QueryParam("category_id")
	pagination := utils.GetPagination(c)

	var equipment []models.Equipment
	var total int64
	var err error
	if categoryID != "" {
		equipment, total, err = ctrl.repo.FindEquipmentByCategoryIDWithPagination(uuid.MustParse(categoryID), pagination.Limit, pagination.Offset)
	} else {
		equipment, total, err = ctrl.repo.FindAllEquipmentWithPagination(pagination.Limit, pagination.Offset)
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	utils.SetPagination(&pagination, total)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":       equipment,
		"pagination": pagination,
	})
}

func (ctrl *EquipmentController) UpdateEquipment(c echo.Context) error {
	slug := c.Param("slug")
	equipment, err := ctrl.repo.FindEquipmentBySlug(slug)
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}
	if err := c.Bind(equipment); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	equipment.Slug = utils.ConvertToSlug(equipment.Name)
	if err := ctrl.repo.UpdateEquipment(equipment); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, equipment)
}

func (ctrl *EquipmentController) DeleteEquipment(c echo.Context) error {
	slug := c.Param("slug")
	equipment, err := ctrl.repo.FindEquipmentBySlug(slug)
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}
	if err := ctrl.repo.DeleteEquipment(equipment.ID); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Successfully deleted equipment",
	})
}
