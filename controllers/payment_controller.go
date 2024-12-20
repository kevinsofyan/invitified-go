package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"invitified-go/models"
	"invitified-go/repositories"
	"invitified-go/utils"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type PaymentController struct {
	paymentRepo repositories.PaymentRepository
	rentalRepo  repositories.RentalRepository
	userRepo    repositories.UserRepository
}

type PaymentRequest struct {
	RentalID      string `json:"rental_id" validate:"required"`
	PaymentMethod string `json:"payment_method" validate:"required,oneof=QR_CODE VIRTUAL_ACCOUNT EWALLET"`
	ChannelCode   string `json:"channel_code" validate:"required"`
}

func NewPaymentController(pr repositories.PaymentRepository, rr repositories.RentalRepository, ur repositories.UserRepository) *PaymentController {
	return &PaymentController{pr, rr, ur}
}

func (ctrl *PaymentController) CreatePayment(c echo.Context) error {
	var req PaymentRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid request format",
		})
	}

	// Parse rental ID
	rentalID, err := uuid.Parse(req.RentalID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid rental ID format",
		})
	}

	// Find rental
	rental, err := ctrl.rentalRepo.FindByID(rentalID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"message": "Rental not found",
		})
	}

	// Get user ID from token
	userIDStr, ok := c.Get("userID").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"message": "Invalid user ID",
		})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"message": "Invalid user ID format",
		})
	}

	user, err := ctrl.userRepo.FindByID(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Failed to find user",
		})
	}

	// Verify rental ownership
	if rental.UserID != userID {
		return c.JSON(http.StatusForbidden, map[string]string{
			"message": "Not authorized to pay for this rental",
		})
	}

	// Create virtual account
	vaRequest := map[string]interface{}{
		"external_id": "VA_fixed-" + time.Now().Format("20060102150405"),
		"bank_code":   req.ChannelCode,
		"name":        "Customer Payment",
	}

	vaResponse, err := createVirtualAccount(vaRequest)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Failed to create virtual account",
		})
	}

	// Log the entire response for debugging
	log.Println("Virtual Account Response:", vaResponse)

	externalID, ok := vaResponse["external_id"].(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Invalid response from virtual account creation: missing external_id",
		})
	}

	accountNumber, ok := vaResponse["account_number"].(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Invalid response from virtual account creation: missing account_number",
		})
	}

	// Simulate payment
	paymentResponse, err := simulatePayment(externalID, rental.TotalCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Failed to simulate payment",
		})
	}

	status, ok := paymentResponse["status"].(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Invalid response from payment simulation: missing status",
		})
	}

	// Create payment record
	payment := &models.Payment{
		ID:                   uuid.New(),
		RentalID:             rental.ID,
		UserID:               userID,
		Amount:               rental.TotalCost,
		PaymentMethod:        req.PaymentMethod,
		PaymentStatus:        status,
		XenditInvoiceID:      vaResponse["id"].(string),
		XenditPaymentURL:     accountNumber,
		XenditPaymentChannel: req.ChannelCode,
	}

	if err := ctrl.paymentRepo.Create(payment); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Failed to save payment",
		})
	}

	// Update rental status to COMPLETED if payment is successful
	if status == "COMPLETED" {
		if err := ctrl.rentalRepo.UpdateStatus(rental.ID, models.RentalStatusComplete); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"message": "Failed to update rental status",
			})
		}
		rental.Status = models.RentalStatusComplete
		subject := "Payment Completed"
		htmlBody := utils.GetOrderConfirmationEmail(rental.ID.String(), fmt.Sprintf("%.2f", rental.TotalCost))
		if err := utils.SendHTMLEmail(user.Email, subject, htmlBody); err != nil {
			log.Println("Failed to send email:", err)
		}
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"payment": payment,
		"rental":  rental,
	})
}

func createVirtualAccount(request map[string]interface{}) (map[string]interface{}, error) {
	url := "https://api.xendit.co/callback_virtual_accounts"
	apiKey := os.Getenv("XENDIT_SECRET_KEY")

	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(apiKey, "")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func simulatePayment(externalID string, amount float64) (map[string]interface{}, error) {
	url := fmt.Sprintf("https://api.xendit.co/callback_virtual_accounts/external_id=%s/simulate_payment", externalID)
	apiKey := os.Getenv("XENDIT_SECRET_KEY")

	request := map[string]interface{}{
		"amount": amount,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(apiKey, "")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}
