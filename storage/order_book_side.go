// CLOB/storage/order_book_side.go
package storage

import "container/heap"

// OrderBookSide represents one side of the order book (buy or sell)
type OrderBookSide struct {
    Side        Side
    PriceLevels map[float64]*PriceLevel // Price to PriceLevel
    Prices      heap.Interface          // Either *BuyHeap or *SellHeap
}

// NewOrderBookSide creates a new OrderBookSide
func NewOrderBookSide(side Side) *OrderBookSide {
    obSide := &OrderBookSide{
        Side:        side,
        PriceLevels: make(map[float64]*PriceLevel),
    }
    if side == Buy {
        bh := &BuyHeap{}
        heap.Init(bh)
        obSide.Prices = bh
    } else {
        sh := &SellHeap{}
        heap.Init(sh)
        obSide.Prices = sh
    }
    return obSide
}

// AddPriceLevel adds a price level to the heap
func (obs *OrderBookSide) AddPriceLevel(priceLevel *PriceLevel) {
    heap.Push(obs.Prices, priceLevel)
}

// RemovePriceLevel removes a price level from the heap
func (obs *OrderBookSide) RemovePriceLevel(priceLevel *PriceLevel) {
    // This is a complex operation in Go's heap; it's simplified here
    // In practice, you might need to implement additional logic
}

// PeekBestPriceLevel returns the best price level without removing it
func (obs *OrderBookSide) PeekBestPriceLevel() *PriceLevel {
    if obs.Prices.Len() == 0 {
        return nil
    }
    return (*obs.Prices).([]*PriceLevel)[0]
}
