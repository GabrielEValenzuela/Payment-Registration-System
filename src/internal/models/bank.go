/*
 * Payment Registration System - Bank Models
 * ----------------------------------------
 * This file defines the data models related to banks and their customer relationships
 * within the payment registration system.
 *
 * Created: Oct. 19, 2024
 * License: GNU General Public License v3.0
 */

package models

// Bank represents a financial institution.
//
//	@Summary		Bank model
//	@Description	Contains details about a bank, including its name, tax identification code (CUIT), address, contact information, and associated customers.
//	@Tags			Models
//	@Accept			json
//	@Produce		json
type Bank struct {
	Name      string     `json:"name" example:"Bank of Argentina"` // Bank name
	Cuit      string     `json:"cuit" example:"30-12345678-9"`     // Bank tax identification code (CUIT)
	Address   string     `json:"address" example:"Av. 9 de Julio"` // Bank address
	Telephone string     `json:"telephone" example:"0800-888-123"` // Bank contact number
	Members   []Customer `json:"customers_ids"`                    // List of customers associated with the bank
}

// BankCustomerCountDTO represents the number of customers associated with a bank.
//
//	@Summary		Bank customer count model
//	@Description	Provides a summary of the number of customers linked to a bank.
//	@Tags			Models
//	@Accept			json
//	@Produce		json
type BankCustomerCountDTO struct {
	BankCuit      string `json:"bank_cuit" example:"30-12345678-9"`     // Unique tax identification code (CUIT) of the bank
	BankName      string `json:"bank_name" example:"Bank of Argentina"` // Name of the bank
	CustomerCount int    `json:"customer_count" example:"2500"`         // Number of customers associated with the bank
}
