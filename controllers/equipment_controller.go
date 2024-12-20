package controllers

import (
	"invitified-go/models"
	"invitified-go/repositories"
	"invitified-go/utils"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// EquipmentController handles equipment-related requests
type EquipmentController struct {
	repo repositories.EquipmentRepository
}

// NewEquipmentController creates a new EquipmentController
func NewEquipmentController(repo repositories.EquipmentRepository) *EquipmentController {
	return &EquipmentController{repo}
}

// CreateCategory godoc
// @Summary Create a new equipment category
// @Description Create a new equipment category
// @Tags categories
// @Accept json
// @Produce json
// @Param category body models.EquipmentCategory true "Equipment Category"
// @Success 201 {object} models.EquipmentCategory
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /categories [post]
func (ctrl *EquipmentController) CreateCategory(c echo.Context) error {
	category := new(models.EquipmentCategory)
	if err := c.Bind(category); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
	}
	category.Slug = utils.ConvertToSlug(category.Name)
	if err := ctrl.repo.CreateCategory(category); err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, category)
}

// GetCategoryByID godoc
// @Summary Get a category by ID
// @Description Get a category by ID
// @Tags categories
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} models.EquipmentCategory
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /categories/{id} [get]
func (ctrl *EquipmentController) GetCategoryByID(c echo.Context) error {
	id := c.Param("id")
	categoryID, err := uuid.Parse(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
	}
	category, err := ctrl.repo.FindCategoryByID(categoryID)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.ErrorResponse{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, category)
}

// GetAllCategories godoc
// @Summary Get all categories
// @Description Get all categories
// @Tags categories
// @Produce json
// @Success 200 {array} models.EquipmentCategory
// @Failure 500 {object} models.ErrorResponse
// @Router /categories [get]
func (ctrl *EquipmentController) GetAllCategories(c echo.Context) error {
	categories, err := ctrl.repo.FindAllCategories()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, categories)
}

// UpdateCategory godoc
// @Summary Update a category
// @Description Update a category
// @Tags categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Param category body models.EquipmentCategory true "Equipment Category"
// @Success 200 {object} models.EquipmentCategory
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Param Authorization header string true "token" default(<token>)
// @Router /categories/{id} [put]
func (ctrl *EquipmentController) UpdateCategory(c echo.Context) error {
	id := c.Param("id")
	categoryID, err := uuid.Parse(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
	}
	category := new(models.EquipmentCategory)
	if err := c.Bind(category); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
	}
	category.ID = categoryID
	category.Slug = utils.ConvertToSlug(category.Name)
	if err := ctrl.repo.UpdateCategory(category); err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, category)
}

// DeleteCategory godoc
// @Summary Delete a category
// @Description Delete a category
// @Tags categories
// @Produce json
// @Param id path string true "Category ID"
// @Success 204
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Param Authorization header string true "token" default(<token>)
// @Router /categories/{id} [delete]
func (ctrl *EquipmentController) DeleteCategory(c echo.Context) error {
	id := c.Param("id")
	categoryID, err := uuid.Parse(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
	}
	if err := ctrl.repo.DeleteCategory(categoryID); err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

// CreateEquipment godoc
// @Summary Create a new equipment
// @Description Create a new equipment
// @Tags equipment
// @Accept json
// @Produce json
// @Param equipment body models.Equipment true "Equipment"
// @Success 201 {object} models.Equipment
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Param Authorization header string true "token" default(<token>)
// @Router /equipment [post]
func (ctrl *EquipmentController) CreateEquipment(c echo.Context) error {
	equipment := new(models.Equipment)
	if err := c.Bind(equipment); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
	}
	equipment.Slug = utils.ConvertToSlug(equipment.Name)

	userIDStr, ok := c.Get("userID").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{Message: "Invalid user ID"})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{Message: "Invalid user ID"})
	}

	equipment.CreatedBy = userID

	if err := ctrl.repo.CreateEquipment(equipment); err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, equipment)
}

// GetEquipmentBySlug godoc
// @Summary Get equipment by slug
// @Description Get equipment by slug
// @Tags equipment
// @Produce json
// @Param slug path string true "Equipment Slug"
// @Success 200 {object} models.Equipment
// @Failure 404 {object} models.ErrorResponse
// @Param Authorization header string true "token" default(<token>)
// @Router /equipment/{slug} [get]
func (ctrl *EquipmentController) GetEquipmentBySlug(c echo.Context) error {
	slug := c.Param("slug")
	equipment, err := ctrl.repo.FindEquipmentBySlug(slug)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.ErrorResponse{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, equipment)
}

// GetAllEquipment godoc
// @Summary Get all equipment
// @Description Get all equipment
// @Tags equipment
// @Produce json
// @Param category_id query string false "Category ID"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} models.ErrorResponse
// @Param Authorization header string true "token" default(<token>)
// @Router /equipment [get]
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
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: err.Error()})
	}

	utils.SetPagination(&pagination, total)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":       equipment,
		"pagination": pagination,
	})
}

// UpdateEquipment godoc
// @Summary Update equipment
// @Description Update equipment
// @Tags equipment
// @Accept json
// @Produce json
// @Param slug path string true "Equipment Slug"
// @Param equipment body models.Equipment true "Equipment"
// @Success 200 {object} models.Equipment
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Param Authorization header string true "token" default(<token>)
// @Router /equipment/{slug} [put]
func (ctrl *EquipmentController) UpdateEquipment(c echo.Context) error {
	slug := c.Param("slug")
	equipment, err := ctrl.repo.FindEquipmentBySlug(slug)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.ErrorResponse{Message: err.Error()})
	}
	if err := c.Bind(equipment); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
	}
	equipment.Slug = utils.ConvertToSlug(equipment.Name)
	if err := ctrl.repo.UpdateEquipment(equipment); err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, equipment)
}

// DeleteEquipment godoc
// @Summary Delete equipment
// @Description Delete equipment
// @Tags equipment
// @Produce json
// @Param slug path string true "Equipment Slug"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Param Authorization header string true "token" default(<token>)
// @Router /equipment/{slug} [delete]
func (ctrl *EquipmentController) DeleteEquipment(c echo.Context) error {
	slug := c.Param("slug")
	equipment, err := ctrl.repo.FindEquipmentBySlug(slug)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.ErrorResponse{Message: err.Error()})
	}
	if err := ctrl.repo.DeleteEquipment(equipment.ID); err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Successfully deleted equipment",
	})
}
