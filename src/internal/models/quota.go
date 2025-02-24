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

// Quota represents an individual installment in a monthly payment plan.
//
//	@Summary		Quota model
//	@Description	Contains details about a single installment, including its number, price, and associated month and year.
//	@Tags			Models
//	@Accept			json
//	@Produce		json
type Quota struct {
	Number int     `json:"number" example:"1"`       // Installment number
	Price  float64 `json:"price" example:"125.50"`   // Price of the installment
	Month  string  `json:"month" example:"February"` // Month when the installment is due
	Year   string  `json:"year" example:"2025"`      // Year when the installment is due
}
