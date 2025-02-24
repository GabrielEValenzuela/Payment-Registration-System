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

/*
 * StoreDTO
 * ----------------------------------------
 * Represents a data transfer object (DTO) for a store, containing the store's name and tax identification code.
 *
 * Fields:
 * - Name (string): The name of the store.
 * - Cuit (string): The unique tax identification code for the store.
 */
type StoreDTO struct {
	Name string `json:"name"`
	Cuit string `json:"cuit"`
}
