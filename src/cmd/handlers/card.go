/*
 * Payment Registration System - Card Handlers
 * -------------------------------------------
 * This file defines the HTTP handlers for managing credit and debit card operations.
 *
 * Created: Oct. 19, 2024
 * License: GNU General Public License v3.0
 */

package handlers

import (
	"strconv"

	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/services"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

type CardHandler struct {
	card services.CardService
}

// NewCardHandler creates a new instance of CardHandler with the provided card service.
func NewCardHandler(card services.CardService) *CardHandler {
	return &CardHandler{
		card: card,
	}
}

// GetPaymentSummary retrieves the payment summary for a specific card and period.
//
//	@Summary		Get payment summary
//	@Description	Retrieves the payment summary for a given card number, month, and year.
//	@Tags			Card
//	@Accept			json
//	@Produce		json
//	@Param			cardNumber	path		string					true	"Card Number"
//	@Param			month		path		int						true	"Month (1-12)"
//	@Param			year		path		int						true	"Year (e.g., 2025)"
//	@Success		200			{object}	map[string]interface{}	"Payment summary retrieved successfully"
//	@Failure		400			{object}	map[string]interface{}	"Invalid month or year parameter"
//	@Failure		500			{object}	map[string]interface{}	"Failed to retrieve payment summary"
//	@Router			/sql/cards/payment-summary/{cardNumber}/{month}/{year} [get]
//	@Router			/no-sql/cards/payment-summary/{cardNumber}/{month}/{year} [get]
func (h *CardHandler) GetPaymentSummary() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Log request
		logger.Info("GetPaymentSummary request from IP: %s", c.IP())

		// Get parameters from the path
		cardNumber := c.Params("cardNumber")
		monthStr := c.Params("month")
		yearStr := c.Params("year")

		// Convert month and year to int
		month, err := strconv.Atoi(monthStr)
		if err != nil {
			logger.Warn("Invalid month parameter")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid month parameter",
			})
		}
		year, err := strconv.Atoi(yearStr)
		if err != nil {
			logger.Warn("Invalid year parameter")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid year parameter",
			})
		}

		// Call the service to get the payment summary
		paymentSummary, err := h.card.GetPaymentSummary(cardNumber, month, year)
		if err != nil {
			logger.Error("Failed to retrieve payment summary: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		logger.Info("Payment summary retrieved successfully")

		if paymentSummary == nil {
			return c.JSON(fiber.Map{
				"message": "Oops! Apparently, there are no data to show at the moment.",
			})
		} else {
			return c.JSON(paymentSummary)
		}
	}
}

// GetCardsExpiringInNext30Days retrieves cards expiring within the next 30 days.
//
//	@Summary		Get cards expiring in the next 30 days
//	@Description	Retrieves a list of cards that will expire within the next 30 days from the given date.
//	@Tags			Card
//	@Accept			json
//	@Produce		json
//	@Param			day		path		int						true	"Day (1-31)"
//	@Param			month	path		int						true	"Month (1-12)"
//	@Param			year	path		int						true	"Year (e.g., 2025)"
//	@Success		200		{object}	map[string]interface{}	"List of cards expiring in the next 30 days"
//	@Failure		400		{object}	map[string]interface{}	"Invalid day, month, or year parameter"
//	@Failure		500		{object}	map[string]interface{}	"Failed to retrieve expiring cards"
//	@Router			/sql/cards/expiring-next-30-days/{day}/{month}/{year} [get]
//	@Router			/no-sql/cards/expiring-next-30-days/{day}/{month}/{year} [get]
func (h *CardHandler) GetCardsExpiringInNext30Days() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Log request
		logger.Info("GetCardsExpiringInNext30Days request from IP: %s", c.IP())

		// Get parameters from the path
		dayStr := c.Params("day")
		monthStr := c.Params("month")
		yearStr := c.Params("year")

		// Convert day, month, and year to int
		day, err := strconv.Atoi(dayStr)
		if err != nil {
			logger.Warn("Invalid day parameter")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid day parameter",
			})
		}
		month, err := strconv.Atoi(monthStr)
		if err != nil {
			logger.Warn("Invalid month parameter")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid month parameter",
			})
		}
		year, err := strconv.Atoi(yearStr)
		if err != nil {
			logger.Warn("Invalid year parameter")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid year parameter",
			})
		}

		// Call the service to get cards expiring in the next 30 days
		cards, err := h.card.GetCardsExpiringInNext30Days(day, month, year)
		if err != nil {
			logger.Error("Failed to retrieve expiring cards: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		logger.Info("Cards expiring in the next 30 days retrieved successfully")

		if cards == nil {
			return c.JSON(fiber.Map{
				"message": "Oops! Apparently, there are no data to show at the moment.",
			})
		} else {
			return c.JSON(cards)
		}
	}
}

// GetPurchaseMonthly retrieves the monthly purchase details.
//
//	@Summary		Get monthly purchase details
//	@Description	Retrieves the purchase details for a given CUIT, final amount, and payment voucher.
//	@Tags			Card
//	@Accept			json
//	@Produce		json
//	@Param			cuit			path		string					true	"CUIT (Unique Tax Identification Code)"
//	@Param			finalAmount		path		float64					true	"Final purchase amount"
//	@Param			paymentVoucher	path		string					true	"Payment voucher identifier"
//	@Success		200				{object}	map[string]interface{}	"Monthly purchase details retrieved successfully"
//	@Failure		400				{object}	map[string]interface{}	"Invalid finalAmount parameter"
//	@Failure		500				{object}	map[string]interface{}	"Failed to retrieve monthly purchase details"
//	@Router			/sql/cards/purchase-monthly/{cuit}/{finalAmount}/{paymentVoucher} [get]
//	@Router			/no-sql/cards/purchase-monthly/{cuit}/{finalAmount}/{paymentVoucher} [get]
func (h *CardHandler) GetPurchaseMonthly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Log request
		logger.Info("GetPurchaseMonthly request from IP: %s", c.IP())

		// Get parameters from the path
		cuit := c.Params("cuit")
		finalAmountStr := c.Params("finalAmount")
		paymentVoucher := c.Params("paymentVoucher")

		// Convert finalAmount to float64
		finalAmount, err := strconv.ParseFloat(finalAmountStr, 64)
		if err != nil {
			logger.Warn("Invalid finalAmount parameter")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid finalAmount parameter",
			})
		}

		// Call the service to get the monthly purchase details
		purchase, err := h.card.GetPurchaseMonthly(cuit, finalAmount, paymentVoucher)
		if err != nil {
			logger.Error("Failed to retrieve monthly purchase details: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		logger.Info("Monthly purchase details retrieved successfully")

		if purchase == nil {
			return c.JSON(fiber.Map{
				"message": "Oops! Apparently, there are no data to show at the moment.",
			})
		} else {
			return c.JSON(purchase)
		}
	}
}

// GetTop10CardsByPurchases retrieves the top 10 cards ranked by purchase volume.
//
//	@Summary		Get top 10 cards by purchases
//	@Description	Retrieves the top 10 cards with the highest purchase volume.
//	@Tags			Card
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	map[string]interface{}	"Top 10 cards by purchases retrieved successfully"
//	@Failure		500	{object}	map[string]interface{}	"Failed to retrieve top 10 cards"
//	@Router			/sql/cards/top [get]
//	@Router			/no-sql/cards/top [get]
func (h *CardHandler) GetTop10CardsByPurchases() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Log request
		logger.Info("GetTop10CardsByPurchases request from IP: %s", c.IP())

		// Call the service to get the top 10 cards by purchases
		cards, err := h.card.GetTop10CardsByPurchases()
		if err != nil {
			logger.Error("Failed to retrieve top 10 cards: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		logger.Info("Top 10 cards by purchases retrieved successfully")

		if cards == nil {
			return c.JSON(fiber.Map{
				"message": "Oops! Apparently, there are no data to show at the moment.",
			})
		} else {
			return c.JSON(cards)
		}
	}
}
