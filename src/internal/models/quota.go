/*
* Payment Registration System - Quota Model
* -----------------------------------------
* This file defines the data model for a quota, representing an installment for a purchase that is being paid in monthly payments.
*
*
* Created: Oct. 19, 2024
* License: GNU General Public License v3.0
 */
package models

/*
 * Quota
 * ----------------------------------------
 * Represents an installment for a purchase that is being paid in monthly payments.
 *
 * Fields:
 * - Number (int): The installment number.
 * - Price (float64): The price of the installment.
 * - Month (string): The month of the installment.
 * - Year (string): The year of the installment.
 */
type Quota struct {
	Number int     `json:"number"`
	Price  float64 `json:"price"`
	Month  string  `json:"month"`
	Year   string  `json:"year"`
}
