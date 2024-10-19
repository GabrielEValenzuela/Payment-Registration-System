package handlers

import (
	"net/http"
	"strconv"

	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/customer"
	"github.com/gofiber/fiber/v2"
)

type CustomerHandler struct {
	customerService customer.Service
}

func NewCustomerHandler(customerService customer.Service) *CustomerHandler {
	return &CustomerHandler{
		customerService: customerService,
	}
}

func (custhandler *CustomerHandler) GetCustomerById() fiber.Handler {
	return func(c *fiber.Ctx) error {
		idParam := c.Params("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid ID",
			})
		}
		customer, err := custhandler.customerService.GetCustomerById(id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal Server Error",
			})
		}
		return c.Status(fiber.StatusOK).JSON(customer)
	}
}

func (custhandler *CustomerHandler) GetAllCustomers() fiber.Handler {
	return func(c *fiber.Ctx) error {
		customers, err := custhandler.customerService.GetAllCustomers()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal Server Error",
			})
		}
		return c.Status(fiber.StatusOK).JSON(customers)
	}
}
