/*
 * Payment Registration System - Store DTO
 * ----------------------------------------
 * This file defines the data transfer object (DTO) for a store, containing the store's name and tax identification code.
 *
 * Authors: marventu94, GabrielEValenzuela
 * Created: Oct. 19, 2024
 * License: GNU General Public License v3.0
 */
package models

// StoreDTO represents a store with its basic details.
//
//	@Summary		StoreDTO model
//	@Description	Contains essential information about a store, including its name and tax identification code (CUIT).
//	@Tags			Models
//	@Accept			json
//	@Produce		json
type StoreDTO struct {
	Name string `json:"name" example:"Tech Store"`    // Store name
	Cuit string `json:"cuit" example:"30-98765432-1"` // Store tax identification code (CUIT)
}
