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

// @Summary Get the payment summary for a card
// @Description Retrieves the payment summary for the specified card number, month, and year.
// @Tags Card
// @Accept json
// @Produce json
// @Param cardNumber path string true "Card number"
// @Param month path int true "Month"
// @Param year path int true "Year"
// @Success 200 {object} map[string]interface{} "Payment summary retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request parameters"
// @Failure 500 {object} map[string]interface{} "Failed to retrieve payment summary"
// @Router /v1/cards/summary/{cardNumber}/{month}/{year} [get]
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

		// Return success response
		logger.Info("Payment summary retrieved successfully")
		return c.JSON(paymentSummary)
	}
}

// @Summary Get cards expiring in the next 30 days
// @Description Retrieves the cards that will expire in the next 30 days based on the current day, month, and year.
// @Tags Card
// @Accept json
// @Produce json
// @Param day path int true "Day"
// @Param month path int true "Month"
// @Param year path int true "Year"
// @Success 200 {array} map[string]interface{} "Cards expiring in the next 30 days"
// @Failure 500 {object} map[string]interface{} "Failed to retrieve expiring cards"
// @Router /v1/cards/expiring/{day}/{month}/{year} [get]
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

		// Return success response
		logger.Info("Cards expiring in the next 30 days retrieved successfully")
		return c.JSON(cards)
	}
}

// @Summary Get monthly purchase details for a card
// @Description Retrieves the monthly purchase details for the specified card based on CUIT and payment voucher.
// @Tags Card
// @Accept json
// @Produce json
// @Param cuit path string true "CUIT"
// @Param finalAmount path float64 true "Final amount"
// @Param paymentVoucher path string true "Payment voucher"
// @Success 200 {object} map[string]interface{} "Monthly purchase details retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request parameters"
// @Failure 500 {object} map[string]interface{} "Failed to retrieve monthly purchase details"
// @Router /v1/cards/purchase/monthly [get]
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

		// Return success response
		logger.Info("Monthly purchase details retrieved successfully")
		return c.JSON(purchase)
	}
}

// @Summary Get top 10 cards by purchases
// @Description Retrieves the top 10 cards by purchases.
// @Tags Card
// @Accept json
// @Produce json
// @Success 200 {array} map[string]interface{} "Top 10 cards by purchases"
// @Failure 500 {object} map[string]interface{} "Failed to retrieve top 10 cards"
// @Router /v1/cards/top [get]
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

		// Return success response
		logger.Info("Top 10 cards by purchases retrieved successfully")
		return c.JSON(cards)
	}
}
