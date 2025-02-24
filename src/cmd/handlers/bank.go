/*
 * Payment Registration System - Bank Handlers
 * -------------------------------------------
 *
 * This file defines the HTTP handlers for bank-related operations in the system.
 *
 * Created: Dec. 11, 2024
 * License: GNU General Public License v3.0
 */

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

// AddFinancingPromotionToBank adds a financing promotion to a bank.
//
//	@Summary		Add a financing promotion to a bank
//	@Description	Adds a new financing promotion using the request body data.
//	@Tags			Bank
//	@Accept			json
//	@Produce		json
//	@Param			request	body		models.Financing		true	"Financing promotion details"
//	@Success		201		{object}	map[string]interface{}	"Financing promotion added successfully"
//	@Failure		400		{object}	map[string]interface{}	"Invalid request body"
//	@Failure		500		{object}	map[string]interface{}	"Failed to add promotion"
//	@Router			/sql/promotions/add-promotion [post]
//	@Router			/no-sql/promotions/add-promotion [post]
func (h *BankHandler) AddFinancingPromotionToBank() fiber.Handler {
	return func(c *fiber.Ctx) error {

		// Log request
		logger.Info("AddFinancingPromotionToBank request from IP: %s", c.IP())

		var promotion models.Financing
		if err := c.BodyParser(&promotion); err != nil {
			logger.Warn("Invalid request body: %v", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Invalid request body",
				"message": err.Error(),
			})
		}

		if err := h.bank.AddFinancingPromotionToBank(promotion); err != nil {
			logger.Error("Failed to add financing promotion: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Failed to add promotion",
				"message": err.Error(),
			})
		}

		logger.Info("Financing promotion added successfully")
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "Financing promotion added successfully",
			"data":    promotion,
		})
	}
}

// ExtendFinancingPromotionValidity extends the validity of an existing financing promotion.
//
//	@Summary		Extend financing promotion validity
//	@Description	Updates the expiration date of a financing promotion identified by its code.
//	@Tags			Bank
//	@Accept			json
//	@Produce		json
//	@Param			code	path		string							true	"Promotion Code"
//	@Param			request	body		models.ExtendPromotionRequest	true	"New expiration date (RFC3339 format)"
//	@Success		200		{object}	map[string]interface{}			"Financing promotion validity extended successfully"
//	@Failure		400		{object}	map[string]interface{}			"Invalid request body or missing promotion code"
//	@Failure		500		{object}	map[string]interface{}			"Failed to extend promotion validity"
//	@Router			/sql/promotions/financing/{code} [patch]
//	@Router			/no-sql/promotions/financing/{code} [patch]
func (h *BankHandler) ExtendFinancingPromotionValidity() fiber.Handler {
	return func(c *fiber.Ctx) error {

		logger.Info("ExtendFinancingPromotionValidity request from IP: %s", c.IP())

		code := c.Params("code")
		if code == "" {
			logger.Warn("Promotion code is missing in the request")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Promotion code is required",
			})
		}

		var requestBody models.ExtendPromotionRequest

		if err := c.BodyParser(&requestBody); err != nil {
			logger.Warn("Invalid request body")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		newDate, err := time.Parse(time.RFC3339, requestBody.NewDate)
		if err != nil {
			logger.Warn("Invalid date format")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid date format. Expected RFC3339 format.",
			})
		}

		err = h.bank.ExtendFinancingPromotionValidity(code, newDate)
		if err != nil {
			logger.Error("Failed to extend financing promotion validity %s due %s", code, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		logger.Info("Financing promotion validity extended successfully %s", code)
		return c.JSON(fiber.Map{
			"message":  "Financing promotion validity extended successfully",
			"code":     code,
			"new_date": newDate.Format(time.RFC3339),
		})
	}
}

// ExtendDiscountPromotionValidity extends the validity of an existing discount promotion.
//
//	@Summary		Extend discount promotion validity
//	@Description	Updates the expiration date of a discount promotion identified by its code.
//	@Tags			Bank
//	@Accept			json
//	@Produce		json
//	@Param			code	path		string							true	"Promotion Code"
//	@Param			request	body		models.ExtendPromotionRequest	true	"New expiration date (RFC3339 format)"
//	@Success		200		{object}	map[string]interface{}			"Discount promotion validity extended successfully"
//	@Failure		400		{object}	map[string]interface{}			"Invalid request body or missing promotion code"
//	@Failure		500		{object}	map[string]interface{}			"Failed to extend promotion validity"
//	@Router			/sql/promotions/discount/{code} [patch]
//	@Router			/no-sql/promotions/discount/{code} [patch]
func (h *BankHandler) ExtendDiscountPromotionValidity() fiber.Handler {
	return func(c *fiber.Ctx) error {

		logger.Info("ExtendDiscountPromotionValidity request from IP: %s", c.IP())

		code := c.Params("code")
		if code == "" {
			logger.Warn("Promotion code is missing in the request")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Promotion code is required",
			})
		}

		var requestBody models.ExtendPromotionRequest

		if err := c.BodyParser(&requestBody); err != nil {
			logger.Warn("Invalid request body")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		newDate, err := time.Parse(time.RFC3339, requestBody.NewDate)
		if err != nil {
			logger.Warn("Invalid date format")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid date format. Expected RFC3339 format.",
			})
		}

		err = h.bank.ExtendDiscountPromotionValidity(code, newDate)
		if err != nil {
			logger.Error("Failed to extend discount promotion validity %s due %s", code, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		logger.Info("Discount promotion validity extended successfully %s", code)
		return c.JSON(fiber.Map{
			"message":  "Discount promotion validity extended successfully",
			"code":     code,
			"new_date": newDate.Format(time.RFC3339),
		})
	}
}

// DeleteFinancingPromotion removes an existing financing promotion.
//
//	@Summary		Delete financing promotion
//	@Description	Deletes a financing promotion identified by its code.
//	@Tags			Bank
//	@Accept			json
//	@Produce		json
//	@Param			code	path		string					true	"Promotion Code"
//	@Success		200		{object}	map[string]interface{}	"Financing promotion deleted successfully"
//	@Failure		400		{object}	map[string]interface{}	"Missing promotion code"
//	@Failure		500		{object}	map[string]interface{}	"Failed to delete promotion"
//	@Router			/sql/promotions/financing/{code} [delete]
//	@Router			/no-sql/promotions/financing/{code} [delete]
func (h *BankHandler) DeleteFinancingPromotion() fiber.Handler {
	return func(c *fiber.Ctx) error {

		logger.Info("DeleteFinancingPromotion request from IP: %s", c.IP())

		code := c.Params("code")
		if code == "" {
			logger.Warn("Promotion code is missing in the request")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Promotion code is required",
			})
		}

		err := h.bank.DeleteFinancingPromotion(code)
		if err != nil {
			logger.Error("Failed to delete financing promotion %s due %s", code, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		logger.Info("Financing promotion deleted successfully %s", code)
		return c.JSON(fiber.Map{
			"message": "Financing promotion deleted successfully",
			"code":    code,
		})
	}
}

// DeleteDiscountPromotion removes an existing discount promotion.
//
//	@Summary		Delete discount promotion
//	@Description	Deletes a discount promotion identified by its code.
//	@Tags			Bank
//	@Accept			json
//	@Produce		json
//	@Param			code	path		string					true	"Promotion Code"
//	@Success		200		{object}	map[string]interface{}	"Discount promotion deleted successfully"
//	@Failure		400		{object}	map[string]interface{}	"Missing promotion code"
//	@Failure		500		{object}	map[string]interface{}	"Failed to delete promotion"
//	@Router			/sql/promotions/discount/{code} [delete]
//	@Router			/no-sql/promotions/discount/{code} [delete]
func (h *BankHandler) DeleteDiscountPromotion() fiber.Handler {
	return func(c *fiber.Ctx) error {

		logger.Info("DeleteDiscountPromotion request from IP: %s", c.IP())

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

// GetBankCustomerCounts retrieves the total number of customers per bank.
//
//	@Summary		Get bank customer counts
//	@Description	Retrieves the number of customers associated with each bank.
//	@Tags			Bank
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	map[string]interface{}	"Bank customer counts retrieved successfully"
//	@Failure		500	{object}	map[string]interface{}	"Failed to get bank customer counts"
//	@Router			/sql/banks/customers/count [get]
//	@Router			/no-sql/banks/customers/count [get]
func (h *BankHandler) GetBankCustomerCounts() fiber.Handler {
	return func(c *fiber.Ctx) error {

		logger.Info("GetBankCustomerCounts request from IP: %s", c.IP())

		customerCounts, err := h.bank.GetBankCustomerCounts()
		if err != nil {
			logger.Error("Failed to get bank customer counts: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		logger.Info("Bank customer counts retrieved successfully")

		if customerCounts == nil {
			return c.JSON(fiber.Map{
				"message": "Oops! Apparently, there are no data to show at the moment.",
			})
		} else {
			return c.JSON(customerCounts)
		}
	}
}
