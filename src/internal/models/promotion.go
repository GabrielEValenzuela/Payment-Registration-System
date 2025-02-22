/*
 * Payment Registration System - Promotion Model
 * ---------------------------------------------
 * This file defines the data model for a promotion, representing a special offer provided by a bank to customers.
 * It applies to specific stores and is valid for a certain period of time.
 *
 * Created: Oct. 19, 2024
 * License: GNU General Public License v3.0
 */

package models

/*
 * Promotion
 * ----------------------------------------
 * Represents a promotion offered by a bank to customers for a specific store.
 *
 * Fields:
 * - Code (string): The unique identifier of the promotion.
 * - PromotionTitle (string): The title or name of the promotion.
 * - NameStore (string): The name of the store where the promotion is valid.
 * - CuitStore (string): The unique tax identification code for the store.
 * - ValidityStartDate (CustomTime): The start date when the promotion is valid.
 * - ValidityEndDate (CustomTime): The end date when the promotion expires.
 * - Comments (string): Additional comments or details about the promotion.
 * - Bank (Bank): The bank offering the promotion.
 */
type Promotion struct {
	Code              string     `json:"code"`
	PromotionTitle    string     `json:"promotion_title"`
	NameStore         string     `json:"name_store"`
	CuitStore         string     `json:"cuit_store"`
	ValidityStartDate CustomTime `json:"validity_start_date"`
	ValidityEndDate   CustomTime `json:"validity_end_date"`
	Comments          string     `json:"comments"`
	Bank              Bank       `json:"bank"`
}

/*
 * Discount
 * ----------------------------------------
 * Represents a promotion that offers a discount on purchases at a specific store.
 *
 * Fields:
 * - Promotion (Promotion): The base promotion details.
 * - DiscountPercentage (float64): The percentage of the discount offered.
 * - PriceCap (float64): The maximum price limit for the discount to apply.
 * - OnlyCash (bool): Indicates if the discount is only applicable for cash payments.
 */
type Discount struct {
	Promotion
	DiscountPercentage float64 `json:"discount_percentage"`
	PriceCap           float64 `json:"price_cap"`
	OnlyCash           bool    `json:"only_cash"`
}

/*
 * Financing
 * ----------------------------------------
 * Represents a promotion that offers financing options for purchases at a specific store.
 *
 * Fields:
 * - Promotion (Promotion): The base promotion details.
 * - NumberOfQuotas (int): The number of installments for the financing.
* - Interest (float64): The interest rate for the financing.
*/
type Financing struct {
	Promotion
	NumberOfQuotas int     `json:"number_of_quotas"`
	Interest       float64 `json:"interest"`
}
