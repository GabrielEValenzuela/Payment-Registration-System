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

// Promotion represents a promotional offer associated with a bank and store.
//
//	@Summary		Promotion model
//	@Description	Contains details about a promotion, including its title, associated store, validity period, and comments.
//	@Tags			Models
//	@Accept			json
//	@Produce		json
type Promotion struct {
	Code              string `json:"code" example:"PROMO2025"`                           // Unique promotion code
	PromotionTitle    string `json:"promotion_title" example:"Holiday Special"`          // Title of the promotion
	NameStore         string `json:"name_store" example:"Tech Store"`                    // Name of the store offering the promotion
	CuitStore         string `json:"cuit_store" example:"30-98765432-1"`                 // Unique tax identification code (CUIT) of the store
	ValidityStartDate string `json:"validity_start_date" example:"2025-01-01T00:00:00Z"` // Change to string for Swagger compatibility
	ValidityEndDate   string `json:"validity_end_date" example:"2026-01-01T00:00:00Z"`   // Change to string for Swagger compatibility
	Comments          string `json:"comments" example:"Limited-time offer!"`             // Additional comments about the promotion
	Bank              Bank   `json:"bank"`                                               // Associated bank details
}

// Discount represents a promotion that offers a discount on purchases.
//
//	@Summary		Discount model
//	@Description	Contains details about a discount promotion, including the discount percentage, price cap, and payment type restrictions.
//	@Tags			Models
//	@Accept			json
//	@Produce		json
type Discount struct {
	Promotion
	DiscountPercentage float64 `json:"discount_percentage" example:"10.5"` // Discount percentage applied to purchases
	PriceCap           float64 `json:"price_cap" example:"5000.00"`        // Maximum price limit for the discount
	OnlyCash           bool    `json:"only_cash" example:"true"`           // Indicates if the discount is cash-only
}

// Financing represents a promotion that provides financing options.
//
//	@Summary		Financing model
//	@Description	Contains details about a financing promotion, including the number of installment payments and interest rate.
//	@Tags			Models
//	@Accept			json
//	@Produce		json
type Financing struct {
	Promotion
	NumberOfQuotas int     `json:"number_of_quotas" example:"12"` // Number of installment payments available
	Interest       float64 `json:"interest" example:"5.5"`        // Interest rate applied to the financing
}

// ExtendPromotionRequest represents a request to extend a promotion's validity period.
//
//	@Summary		Extend promotion request model
//	@Description	Used to update the expiration date of a promotion.
//	@Tags			Models
//	@Accept			json
//	@Produce		json
type ExtendPromotionRequest struct {
	NewDate string `json:"new_date" example:"2026-01-01T00:00:00Z"` // New expiration date in RFC3339 format
}
