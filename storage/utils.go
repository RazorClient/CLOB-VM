// CLOB/storage/utils.go
package storage

// Min returns the minimum of two float64 numbers
func Min(a, b float64) float64 {
    if a < b {
        return a
    }
    return b
}

// SettleTrade simulates the settlement of a trade between two orders
func SettleTrade(order1 *Order, order2 *Order, quantity float64) {
    // In a real system, this would update balances, notify parties, etc.
    // For now, we can log the trade
    println("Trade executed between", order1.ID, "and", order2.ID, "for", quantity, "units at price", order2.Price)
}

// GetPriceComparator returns a comparison function based on the side
func GetPriceComparator(side Side) func(float64, float64) bool {
    if side == Buy {
        return func(a, b float64) bool { return a <= b }
    }
    return func(a, b float64) bool { return a >= b }
}
