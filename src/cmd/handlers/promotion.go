/*
 * Payment Registration System - Promotion Handlers
 * ------------------------------------------------
 * This file defines the HTTP handlers for managing financing and discount promotions.
 *
 * Created: Oct. 19, 2024
 * License: GNU General Public License v3.0
 */

package handlers

import (
	"time"

	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/services"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

type PromotionHandler struct {
	promotion services.PromotionService
}

// NewPromotionHandler creates a new instance of PromotionHandler with the provided promotion service.
func NewPromotionHandler(promotion services.PromotionService) *PromotionHandler {
	return &PromotionHandler{
		promotion: promotion,
	}
}

// @Summary Get available promotions by store and date range
// @Description Retrieves available financing and discount promotions for a specific store within a date range.
// @Tags Promotion
// @Accept json
// @Produce json
// @Param cuit path string true "Store CUIT"
// @Param startDate path string true "Start date in RFC3339 format"
// @Param endDate path string true "End date in RFC3339 format"
// @Success 200 {array} map[string]interface{} "Available financing promotions"
// @Success 200 {array} map[string]interface{} "Available discount promotions"
// @Failure 400 {object} map[string]interface{} "Invalid request parameters"
// @Failure 500 {object} map[string]interface{} "Failed to retrieve promotions"
// @Router /v1/promotions/{cuit}/{startDate}/{endDate} [get]
func (h *PromotionHandler) GetAvailablePromotionsByStoreAndDateRange() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Log request
		logger.Info("GetAvailablePromotionsByStoreAndDateRange request from IP: %s", c.IP())

		// Get parameters from the path
		cuit := c.Params("cuit")
		startDateStr := c.Params("startDate")
		endDateStr := c.Params("endDate")

		// Convert dates from string to time.Time
		startDate, err := time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			logger.Warn("Invalid startDate format")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid startDate format. Expected RFC3339 format.",
			})
		}
		endDate, err := time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			logger.Warn("Invalid endDate format")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid endDate format. Expected RFC3339 format.",
			})
		}

		// Call the service to get the available promotions
		financingPromotions, discountPromotions, err := h.promotion.GetAvailablePromotionsByStoreAndDateRange(cuit, startDate, endDate)
		if err != nil {
			logger.Error("Failed to retrieve available promotions: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// Return success response with financing and discount promotions
		logger.Info("Available promotions retrieved successfully")
		return c.JSON(fiber.Map{
			"financing_promotions": financingPromotions,
			"discount_promotions":  discountPromotions,
		})
	}
}

// @Summary Get the most used promotion
// @Description Retrieves the most used promotion in the system.
// @Tags Promotion
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Most used promotion"
// @Failure 500 {object} map[string]interface{} "Failed to retrieve the most used promotion"
// @Router /v1/promotions/most-used [get]
func (h *PromotionHandler) GetMostUsedPromotion() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Log request
		logger.Info("GetMostUsedPromotion request from IP: %s", c.IP())

		// Call the service to get the most used promotion
		promotion, err := h.promotion.GetMostUsedPromotion()
		if err != nil {
			logger.Error("Failed to retrieve most used promotion: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// Return success response with the most used promotion
		logger.Info("Most used promotion retrieved successfully")
		return c.JSON(promotion)
	}
}
