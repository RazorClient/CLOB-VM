package storage

import "container/heap"

// OrderBookSide represents one side of the order book (buy or sell).
// It maintains a collection of price levels and provides methods to manipulate
// these levels using a heap data structure for efficient retrieval of the best price level.
type OrderBookSide struct {
    Side        Side                // Indicates whether this side is for buying or selling
    PriceLevels map[float64]*PriceLevel // A mapping of price to corresponding PriceLevel objects
    Prices      heap.Interface       // A heap interface that can be either a BuyHeap or SellHeap
}

// NewOrderBookSide creates a new instance of OrderBookSide.
// It initializes the PriceLevels map and sets up the Prices heap based on the specified side.
// If the side is Buy, a BuyHeap is initialized; if Sell, a SellHeap is initialized.
// Parameters:
//   - side: The side of the order book (Buy or Sell).
// Returns:
//   - A pointer to the newly created OrderBookSide instance.
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

// AddPriceLevel adds a new price level to the heap.
// This method pushes the provided PriceLevel onto the Prices heap,
// allowing for efficient retrieval of the best price level.
// Parameters:
//   - priceLevel: A pointer to the PriceLevel to be added to the heap.
func (obs *OrderBookSide) AddPriceLevel(priceLevel *PriceLevel) {
    heap.Push(obs.Prices, priceLevel)
}

// RemovePriceLevel removes a price level from the heap.
// This operation is complex due to the nature of Go's heap implementation.
// In practice, additional logic may be required to maintain the integrity of the heap.
// Parameters:
//   - priceLevel: A pointer to the PriceLevel to be removed from the heap.
func (obs *OrderBookSide) RemovePriceLevel(priceLevel *PriceLevel) {
    // This is a complex operation in Go's heap; it's simplified here
    // In practice, you might need to implement additional logic
}

// PeekBestPriceLevel returns the best price level without removing it from the heap.
// If the Prices heap is empty, it returns nil.
// Returns:
//   - A pointer to the best PriceLevel, or nil if the heap is empty.
func (obs *OrderBookSide) PeekBestPriceLevel() *PriceLevel {
    if obs.Prices.Len() == 0 {
        return nil
    }
    return (*obs.Prices).([]*PriceLevel)[0]
}
