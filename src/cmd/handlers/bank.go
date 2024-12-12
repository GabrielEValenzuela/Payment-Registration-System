package handlers

import (
	"time"

	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/services"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

type BankHandler struct {
	bank services.BankService
}

// NewBankHandler creates a new instance of BankHandler with the provided bank service.
func NewBankHandler(bank services.BankService) *BankHandler {
	return &BankHandler{
		bank: bank,
	}
}

// @Summary Add a financing promotion to a bank
// @Description Adds a new financing promotion to the bank using the provided details in the request body.
// @Tags Bank
// @Accept json
// @Produce json
// @Success 201 {object} map[string]interface{} "Financing promotion added successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 500 {object} map[string]interface{} "Failed to add financing promotion to bank"
// @Router /v1/sql/promotions/financing [post]
func (h *BankHandler) AddFinancingPromotionToBank() fiber.Handler {
	return func(c *fiber.Ctx) error {

		// Log request
		logger.Info("AddFinancingPromotionToBank request from IP: %s", c.IP())

		// Parse request body to FinancingEntity
		var promotion models.Financing
		if err := c.BodyParser(&promotion); err != nil {
			logger.Warn("Invalid request body")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		// Add promotion to the bank using the service
		err := h.bank.AddFinancingPromotionToBank(promotion)
		if err != nil {
			logger.Error("Failed to add financing promotion to bank: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// Return success response
		logger.Info("Financing promotion added successfully")
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "Financing promotion added successfully",
			"data":    promotion,
		})
	}
}

// @Summary Extend the validity of a financing promotion
// @Description Extends the validity of a financing promotion using the provided details in the request body.
// @Tags Bank
// @Accept json
// @Produce json
// @Param code path string true "Promotion code"
// @Success 200 {object} map[string]interface{} "Financing promotion validity extended successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 404 {object} map[string]interface{} "Promotion not found"
// @Failure 500 {object} map[string]interface{} "Failed to extend financing promotion validity"
// @Router /promotions/financing/{code}/extend [patch]
func (h *BankHandler) ExtendFinancingPromotionValidity() fiber.Handler {
	return func(c *fiber.Ctx) error {

		// Log request
		logger.Info("ExtendFinancingPromotionValidity request from IP: %s", c.IP())

		// Get promotion code from the path
		code := c.Params("code")
		if code == "" {
			logger.Warn("Promotion code is missing in the request")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Promotion code is required",
			})
		}

		// Parse request body to get the new date
		var requestBody struct {
			NewDate string `json:"new_date"`
		}

		if err := c.BodyParser(&requestBody); err != nil {
			logger.Warn("Invalid request body")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		// Parse the provided date string into a time.Time object
		newDate, err := time.Parse(time.RFC3339, requestBody.NewDate)
		if err != nil {
			logger.Warn("Invalid date format")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid date format. Expected RFC3339 format.",
			})
		}

		// Call the service to extend the promotion validity
		err = h.bank.ExtendFinancingPromotionValidity(code, newDate)
		if err != nil {
			logger.Error("Failed to extend financing promotion validity %s due %s", code, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// Return success response
		logger.Info("Financing promotion validity extended successfully %s", code)
		return c.JSON(fiber.Map{
			"message":  "Financing promotion validity extended successfully",
			"code":     code,
			"new_date": newDate.Format(time.RFC3339),
		})
	}
}

// @Summary Extend the validity of a discount promotion
// @Description Extends the validity of a discount promotion using the provided details in the request body.
// @Tags Bank
// @Accept json
// @Produce json
// @Param code path string true "Promotion code"
// @Success 200 {object} map[string]interface{} "Discount promotion validity extended successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 404 {object} map[string]interface{} "Promotion not found"
// @Failure 500 {object} map[string]interface{} "Failed to extend discount promotion validity"
// @Router /promotions/discount/{code}/extend [patch]
func (h *BankHandler) ExtendDiscountPromotionValidity() fiber.Handler {
	return func(c *fiber.Ctx) error {

		// Log request
		logger.Info("ExtendDiscountPromotionValidity request from IP: %s", c.IP())

		// Get promotion code from the path
		code := c.Params("code")
		if code == "" {
			logger.Warn("Promotion code is missing in the request")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Promotion code is required",
			})
		}

		// Parse request body to get the new date
		var requestBody struct {
			NewDate string `json:"new_date"`
		}

		if err := c.BodyParser(&requestBody); err != nil {
			logger.Warn("Invalid request body")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		// Parse the provided date string into a time.Time object
		newDate, err := time.Parse(time.RFC3339, requestBody.NewDate)
		if err != nil {
			logger.Warn("Invalid date format")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid date format. Expected RFC3339 format.",
			})
		}

		// Call the service to extend the promotion validity
		err = h.bank.ExtendDiscountPromotionValidity(code, newDate)
		if err != nil {
			logger.Error("Failed to extend discount promotion validity %s due %s", code, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// Return success response
		logger.Info("Discount promotion validity extended successfully %s", code)
		return c.JSON(fiber.Map{
			"message":  "Discount promotion validity extended successfully",
			"code":     code,
			"new_date": newDate.Format(time.RFC3339),
		})
	}
}

// @Summary Delete a financing promotion
// @Description Deletes a financing promotion from the bank using the provided promotion code.
// @Tags Bank
// @Accept json
// @Produce json
// @Param code path string true "Promotion code"
// @Success 200 {object} map[string]interface{} "Financing promotion deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 404 {object} map[string]interface{} "Promotion not found"
// @Failure 500 {object} map[string]interface{} "Failed to delete financing promotion"
// @Router /promotions/financing/{code} [delete]
func (h *BankHandler) DeleteFinancingPromotion() fiber.Handler {
	return func(c *fiber.Ctx) error {

		// Log request
		logger.Info("DeleteFinancingPromotion request from IP: %s", c.IP())

		// Get promotion code from the path
		code := c.Params("code")
		if code == "" {
			logger.Warn("Promotion code is missing in the request")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Promotion code is required",
			})
		}

		// Call the service to delete the promotion
		err := h.bank.DeleteFinancingPromotion(code)
		if err != nil {
			logger.Error("Failed to delete financing promotion %s due %s", code, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// Return success response
		logger.Info("Financing promotion deleted successfully %s", code)
		return c.JSON(fiber.Map{
			"message": "Financing promotion deleted successfully",
			"code":    code,
		})
	}
}

// @Summary Delete a discount promotion
// @Description Deletes a discount promotion from the bank using the provided promotion code.
// @Tags Bank
// @Accept json
// @Produce json
// @Param code path string true "Promotion code"
// @Success 200 {object} map[string]interface{} "Discount promotion deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 404 {object} map[string]interface{} "Promotion not found"
// @Failure 500 {object} map[string]interface{} "Failed to delete discount promotion"
// @Router /promotions/discount/{code} [delete]
func (h *BankHandler) DeleteDiscountPromotion() fiber.Handler {
	return func(c *fiber.Ctx) error {

		// Log request
		logger.Info("DeleteDiscountPromotion request from IP: %s", c.IP())

		// Get promotion code from the path
		code := c.Params("code")
		if code == "" {
			logger.Warn("Promotion code is missing in the request")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Promotion code is required",
			})
		}

		// Call the service to delete the promotion
		err := h.bank.DeleteDiscountPromotion(code)
		if err != nil {
			logger.Error("Failed to delete discount promotion %s due %s", code, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// Return success response
		logger.Info("Discount promotion deleted successfully %s", code)
		return c.JSON(fiber.Map{
			"message": "Discount promotion deleted successfully",
			"code":    code,
		})
	}
}

// @Summary Get the number of customers for each bank
// @Description Returns the number of customers for each bank in the system.
// @Tags Bank
// @Accept json
// @Produce json
// @Failure 500 {object} map[string]interface{} "Failed to get bank customer counts"
// @Router /customers/count [get]
func (h *BankHandler) GetBankCustomerCounts() fiber.Handler {
	return func(c *fiber.Ctx) error {

		// Log request
		logger.Info("GetBankCustomerCounts request from IP: %s", c.IP())

		// Call the service to get the customer counts
		customerCounts, err := h.bank.GetBankCustomerCounts()
		if err != nil {
			logger.Error("Failed to get bank customer counts: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// Return success response
		logger.Info("Bank customer counts retrieved successfully")
		return c.JSON(customerCounts)
	}
}
