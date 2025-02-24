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

// GetAvailablePromotionsByStoreAndDateRange retrieves available promotions for a store within a specified date range.
//
//	@Summary		Get available promotions by store and date range
//	@Description	Retrieves the financing and discount promotions available for a store between the specified start and end dates.
//	@Tags			Promotion
//	@Accept			json
//	@Produce		json
//	@Param			cuit		path		string					true	"CUIT (Unique Tax Identification Code of the store)"
//	@Param			startDate	path		string					true	"Start date (RFC3339 format)"
//	@Param			endDate		path		string					true	"End date (RFC3339 format)"
//	@Success		200			{object}	map[string]interface{}	"Available promotions retrieved successfully"
//	@Failure		400			{object}	map[string]interface{}	"Invalid startDate or endDate format"
//	@Failure		500			{object}	map[string]interface{}	"Failed to retrieve available promotions"
//	@Router			/sql/promotions/available/{cuit}/{startDate}/{endDate} [get]
//	@Router			/no-sql/promotions/available/{cuit}/{startDate}/{endDate} [get]
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

		if financingPromotions == nil && discountPromotions == nil {
			return c.JSON(fiber.Map{
				"message": "Oops! Apparently, there are no data to show at the moment.",
			})
		} else {
			logger.Info("Available promotions retrieved successfully")
			return c.JSON(fiber.Map{
				"financing_promotions": financingPromotions,
				"discount_promotions":  discountPromotions,
			})
		}
	}
}

// GetMostUsedPromotion retrieves the most used promotion.
//
//	@Summary		Get most used promotion
//	@Description	Retrieves the promotion that has been used the most.
//	@Tags			Promotion
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	map[string]interface{}	"Most used promotion retrieved successfully"
//	@Failure		500	{object}	map[string]interface{}	"Failed to retrieve most used promotion"
//	@Router			/sql/promotions/most-used [get]
//	@Router			/no-sql/promotions/most-used [get]
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

		logger.Info("Most used promotion retrieved successfully")

		if promotion == nil {
			return c.JSON(fiber.Map{
				"message": "Oops! Apparently, there are no data to show at the moment.",
			})
		} else {
			return c.JSON(promotion)
		}
	}
}
