package handlers

import (
	"net/http"
	"strconv"

	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/customer"
	"github.com/gofiber/fiber/v2"
)

type CustomerHandler struct {
	customerSerivce customer.Service
}

func NewCustomerHandler(customerService customer.Service) *CustomerHandler {
	return &CustomerHandler{
		customerSerivce: customerService,
	}
}

func (h *CustomerHandler) GetCustomerById() fiber.Handler {
	return func(c *fiber.Ctx) error {
		idParam := c.Params("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid ID",
			})
		}
		customer, err := h.customerSerivce.GetCustomerById(id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal Server Error",
			})
		}
		return c.Status(fiber.StatusOK).JSON(customer)
	}
}
