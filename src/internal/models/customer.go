/*
* Payment Registration System - Customer Model
* --------------------------------------------
* This file defines the data model for a customer, representing a person who is a client
* of one or more banks. A customer can have multiple cards and be associated with different banks.
*
* Created: Oct. 19, 2024
* License: GNU General Public License v3.0
 */
package models

import (
	"time"
)

// Customer represents an individual who holds bank accounts and cards.
//
//	@Summary		Customer model
//	@Description	Contains personal details of a customer, including identification information, contact details, associated banks, and linked cards.
//	@Tags			Models
//	@Accept			json
//	@Produce		json
type Customer struct {
	CompleteName string    `json:"complete_name" example:"John Doe"`          // Full name of the customer
	Dni          string    `json:"dni" example:"12345678"`                    // National identification number (DNI)
	Cuit         string    `json:"cuit" example:"20-12345678-9"`              // Unique tax identification code (CUIT)
	Address      string    `json:"address" example:"456 Oak St, City"`        // Customer's residential address
	Telephone    string    `json:"telephone" example:"+54 11 9876-5432"`      // Contact number
	EntryDate    time.Time `json:"entry_date" example:"2022-03-15T00:00:00Z"` // Date the customer was registered
	BanksIds     []int     `json:"banks_ids"`                                 // List of bank IDs the customer is associated with
	Cards        []int     `json:"cards"`                                     // List of card IDs linked to the customer
}
