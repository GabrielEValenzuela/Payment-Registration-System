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

/*
 * Bank
 * ----------------------------------------
 * Represents a financial institution that holds customers and issues cards.
 *
 * Fields:
 * - Name (string): The name of the bank.
 * - Cuit (string): The unique tax identification code for the bank.
 * - Address (string): The physical address of the bank.
 * - Telephone (string): The contact number for the bank.
 * - Members ([]Customer): A list of associated customers, represented by their IDs.
 */
type Bank struct {
	Name      string     `json:"name"`          // Bank name
	Cuit      string     `json:"cuit"`          // Bank tax identification code (CUIT)
	Address   string     `json:"address"`       // Bank address
	Telephone string     `json:"telephone"`     // Bank contact number
	Members   []Customer `json:"customers_ids"` // List of customers associated with the bank
}

/*
 * BankCustomerCountDTO
 * ----------------------------------------
 * Represents a data transfer object (DTO) for aggregating customer counts per bank.
 *
 * Fields:
 * - BankCuit (string): The unique tax identification code for the bank.
 * - BankName (string): The name of the bank.
 * - CustomerCount (int): The total number of customers associated with the bank.
 */
type BankCustomerCountDTO struct {
	BankCuit      string `json:"bank_cuit"`      // Unique tax identification code (CUIT) of the bank
	BankName      string `json:"bank_name"`      // Name of the bank
	CustomerCount int    `json:"customer_count"` // Number of customers associated with the bank
}
