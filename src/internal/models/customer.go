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

/* Customer
 * ----------------------------------------
 * Represents a customer who is a client of one or more banks.
 *
 * Fields:
 * - CompleteName (string): The full name of the customer.
 * - Dni (string): The national identification number of the customer.
 * - Cuit (string): The unique tax identification code for the customer.
 * - Address (string): The physical address of the customer.
 * - Telephone (string): The contact number for the customer.
 * - EntryDate (time.Time): The date when the customer was registered.
 * - BanksIds ([]int): A list of bank IDs associated with the customer.
 * - Cards ([]int): A list of card IDs associated with the customer.
 */
type Customer struct {
	CompleteName string    `json:"complete_name" validate:"required"`
	Dni          string    `json:"dni" validate:"required"`
	Cuit         string    `json:"cuit" validate:"required"`
	Address      string    `json:"address" validate:"required"`
	Telephone    string    `json:"telephone" validate:"required"`
	EntryDate    time.Time `json:"entry_date" validate:"required"`
	BanksIds     []int     `json:"banks_ids" validate:"dive,required"`
	Cards        []int     `json:"cards" validate:"dive,required"`
}
