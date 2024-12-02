// CLOB/storage/order.go
package storage

import "time"

// Side represents the side of an order: Buy or Sell
type Side string

const (
    Buy  Side = "buy"
    Sell Side = "sell"
)

// OrderType represents the type of an order: Limit or Market
type OrderType string

const (
    Limit  OrderType = "limit"
    Market OrderType = "market"
)

// Order represents an individual order in the order book
type Order struct {
    ID        string
    Side      Side
    Price     float64       // 0 for market orders
    Quantity  float64
    Timestamp time.Time
    OrderType OrderType
    next      *Order        // For linked list 
    prev      *Order        // For linked list 
}
