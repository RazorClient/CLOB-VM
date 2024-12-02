// CLOB/storage/state.go
package storage

import (
    "container/heap"
    "errors"
)

// Errors
var (
    ErrOrderNotFound = errors.New("order not found")
)

// OrderBook represents the entire order book
type OrderBook struct {
    Bids     *OrderBookSide
    Asks     *OrderBookSide
    OrderMap map[string]*Order // Maps Order ID to Order
}

// NewOrderBook creates a new OrderBook
func NewOrderBook() *OrderBook {
    return &OrderBook{
        Bids:     NewOrderBookSide(Buy),
        Asks:     NewOrderBookSide(Sell),
        OrderMap: make(map[string]*Order),
    }
}

// GetOppositeSide returns the opposite side of the given side
func (ob *OrderBook) GetOppositeSide(side Side) *OrderBookSide {
    if side == Buy {
        return ob.Asks
    }
    return ob.Bids
}

// AddLimitOrder adds a limit order to the appropriate side
func (ob *OrderBook) AddLimitOrder(order *Order) error {
    side := ob.GetSide(order.Side)
    priceLevel, exists := side.PriceLevels[order.Price]
    if !exists {
        priceLevel = &PriceLevel{
            Price:  order.Price,
            Orders: NewOrderQueue(),
        }
        side.PriceLevels[order.Price] = priceLevel
        side.AddPriceLevel(priceLevel)
    }
    priceLevel.Orders.Enqueue(order)
    return nil
}

// CancelOrder removes an order from the order book
func (ob *OrderBook) CancelOrder(order *Order) error {
    side := ob.GetSide(order.Side)
    priceLevel, exists := side.PriceLevels[order.Price]
    if !exists {
        return ErrOrderNotFound
    }
    // Remove the order from the queue
    priceLevel.Orders.Remove(order)
    delete(ob.OrderMap, order.ID)

    // If the price level is empty, remove it
    if priceLevel.Orders.Size == 0 {
        side.RemovePriceLevel(priceLevel)
        delete(side.PriceLevels, order.Price)
    }
    return nil
}

// GetSide returns the OrderBookSide for the given side
func (ob *OrderBook) GetSide(side Side) *OrderBookSide {
    if side == Buy {
        return ob.Bids
    }
    return ob.Asks
}
