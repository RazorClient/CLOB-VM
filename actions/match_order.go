// CLOB/actions/match_order.go
package actions

import (
    "container/heap"
    "errors"
    "CLOB/storage"
)

// MatchMarketOrder processes a market order
func MatchMarketOrder(orderBook *storage.OrderBook, order *storage.Order) error {
    // Get the opposite side of the order (buy/sell)
    oppositeSide := orderBook.GetOppositeSide(order.Side)
    remainingQty := order.Quantity

    // Loop until the order is fully matched or no orders left on the opposite side
    for remainingQty > 0 && oppositeSide.Prices.Len() > 0 {
        // Get the best price level from the opposite side
        bestPriceLevel := heap.Pop(oppositeSide.Prices).(*storage.PriceLevel)
        ordersQueue := bestPriceLevel.Orders

        // Process orders in the queue until the market order is matched or queue is empty
        for ordersQueue.Size > 0 && remainingQty > 0 {
            headOrder := ordersQueue.Dequeue() // Get the next order in the queue
            tradeQty := storage.Min(remainingQty, headOrder.Quantity) // Determine trade quantity

            // Settle the trade
            storage.SettleTrade(order, headOrder, tradeQty)

            remainingQty -= tradeQty // Decrease remaining quantity
            headOrder.Quantity -= tradeQty // Decrease head order quantity

            // Remove head order if fully matched
            if headOrder.Quantity == 0 {
                delete(orderBook.OrderMap, headOrder.ID)
            }

            // Break if the market order is fully matched
            if remainingQty == 0 {
                break
            }
        }

        // Remove price level if no orders left, otherwise push it back
        if ordersQueue.Size == 0 {
            delete(oppositeSide.PriceLevels, bestPriceLevel.Price)
        } else {
            heap.Push(oppositeSide.Prices, bestPriceLevel)
        }
    }

    // Return error if the market order could not be fully matched
    if remainingQty == 0 {
        delete(orderBook.OrderMap, order.ID)
        return nil
    } else {
        return errors.New("market order could not be fully matched")
    }
}

// MatchLimitOrder processes a limit order
func MatchLimitOrder(orderBook *storage.OrderBook, order *storage.Order) error {
    // Get the opposite side of the order and the price comparator
    oppositeSide := orderBook.GetOppositeSide(order.Side)
    compare := storage.GetPriceComparator(order.Side)
    remainingQty := order.Quantity

    // Loop until the order is fully matched or no orders left on the opposite side
    for remainingQty > 0 && oppositeSide.Prices.Len() > 0 {
        bestPriceLevel := oppositeSide.PeekBestPriceLevel() // Peek the best price level

        // Check if the best price level meets the limit order's price criteria
        if compare(bestPriceLevel.Price, order.Price) {
            heap.Pop(oppositeSide.Prices) // Remove the best price level
            ordersQueue := bestPriceLevel.Orders

            // Process orders in the queue until the limit order is matched or queue is empty
            for ordersQueue.Size > 0 && remainingQty > 0 {
                headOrder := ordersQueue.Dequeue() // Get the next order in the queue
                tradeQty := storage.Min(remainingQty, headOrder.Quantity) // Determine trade quantity

                // Settle the trade
                storage.SettleTrade(order, headOrder, tradeQty)

                remainingQty -= tradeQty // Decrease remaining quantity
                headOrder.Quantity -= tradeQty // Decrease head order quantity

                // Remove head order if fully matched
                if headOrder.Quantity == 0 {
                    delete(orderBook.OrderMap, headOrder.ID)
                }

                // Break if the limit order is fully matched
                if remainingQty == 0 {
                    break
                }
            }

            // Remove price level if no orders left, otherwise push it back
            if ordersQueue.Size == 0 {
                delete(oppositeSide.PriceLevels, bestPriceLevel.Price)
            } else {
                heap.Push(oppositeSide.Prices, bestPriceLevel)
            }
        } else {
            break // Exit if the best price does not meet the limit order's criteria
        }
    }

    // If the limit order is not fully matched, update its quantity and add it back to the order book
    if remainingQty > 0 {
        order.Quantity = remainingQty
        return orderBook.AddLimitOrder(order)
    } else {
        delete(orderBook.OrderMap, order.ID) // Remove the order if fully matched
    }
    return nil
}
