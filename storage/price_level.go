// CLOB/storage/price_level.go
package storage

// PriceLevel represents a level in the order book at a specific price
type PriceLevel struct {
    Price  float64
    Orders *OrderQueue
}
