/*
 * Payment Registration System - Store Handlers
 * ------------------------------------------------
 * This file defines the HTTP handlers for managing store information.
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

type StoreHandler struct {
	store services.StoreService
}

// NewStoreHandler creates a new instance of StoreHandler with the provided store service.
func NewStoreHandler(store services.StoreService) *StoreHandler {
	return &StoreHandler{
		store: store,
	}
}

// @Summary Get the store with the highest revenue by month
// @Description Retrieves the store with the highest revenue for a specific month and year.
// @Tags Store
// @Accept json
// @Produce json
// @Param month path int true "Month"
// @Param year path int true "Year"
// @Success 200 {object} map[string]interface{} "Store with the highest revenue"
// @Failure 400 {object} map[string]interface{} "Invalid request parameters"
// @Failure 500 {object} map[string]interface{} "Failed to retrieve store information"
// @Router /v1/stores/highest-revenue/{month}/{year} [get]
func (h *StoreHandler) GetStoreWithHighestRevenueByMonth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Log request
		logger.Info("GetStoreWithHighestRevenueByMonth request from IP: %s", c.IP())

		// Get parameters from the path
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

		// Call the service to get the store with the highest revenue
		store, err := h.store.GetStoreWithHighestRevenueByMonth(month, year)
		if err != nil {
			logger.Error("Failed to retrieve store with highest revenue: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// Return success response
		logger.Info("Store with highest revenue retrieved successfully")
		return c.JSON(store)
	}
}
