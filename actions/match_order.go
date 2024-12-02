// CLOB/actions/match_order.go
package actions

import (
    "container/heap"
    "errors"
    "CLOB/storage"
)

// MatchMarketOrder processes a market order
func MatchMarketOrder(orderBook *storage.OrderBook, order *storage.Order) error {
    oppositeSide := orderBook.GetOppositeSide(order.Side)
    remainingQty := order.Quantity

    for remainingQty > 0 && oppositeSide.Prices.Len() > 0 {
        bestPriceLevel := heap.Pop(oppositeSide.Prices).(*storage.PriceLevel)
        ordersQueue := bestPriceLevel.Orders

        for ordersQueue.Size > 0 && remainingQty > 0 {
            headOrder := ordersQueue.Dequeue()
            tradeQty := storage.Min(remainingQty, headOrder.Quantity)

            // Settle the trade
            storage.SettleTrade(order, headOrder, tradeQty)

            remainingQty -= tradeQty
            headOrder.Quantity -= tradeQty

            if headOrder.Quantity == 0 {
                delete(orderBook.OrderMap, headOrder.ID)
            }

            if remainingQty == 0 {
                break
            }
        }

        if ordersQueue.Size == 0 {
            delete(oppositeSide.PriceLevels, bestPriceLevel.Price)
        } else {
            heap.Push(oppositeSide.Prices, bestPriceLevel)
        }
    }

    if remainingQty == 0 {
        delete(orderBook.OrderMap, order.ID)
        return nil
    } else {
        return errors.New("market order could not be fully matched")
    }
}

// MatchLimitOrder processes a limit order
func MatchLimitOrder(orderBook *storage.OrderBook, order *storage.Order) error {
    oppositeSide := orderBook.GetOppositeSide(order.Side)
    compare := storage.GetPriceComparator(order.Side)
    remainingQty := order.Quantity

    for remainingQty > 0 && oppositeSide.Prices.Len() > 0 {
        bestPriceLevel := oppositeSide.PeekBestPriceLevel()

        if compare(bestPriceLevel.Price, order.Price) {
            heap.Pop(oppositeSide.Prices)
            ordersQueue := bestPriceLevel.Orders

            for ordersQueue.Size > 0 && remainingQty > 0 {
                headOrder := ordersQueue.Dequeue()
                tradeQty := storage.Min(remainingQty, headOrder.Quantity)

                // Settle the trade
                storage.SettleTrade(order, headOrder, tradeQty)

                remainingQty -= tradeQty
                headOrder.Quantity -= tradeQty

                if headOrder.Quantity == 0 {
                    delete(orderBook.OrderMap, headOrder.ID)
                }

                if remainingQty == 0 {
                    break
                }
            }

            if ordersQueue.Size == 0 {
                delete(oppositeSide.PriceLevels, bestPriceLevel.Price)
            } else {
                heap.Push(oppositeSide.Prices, bestPriceLevel)
            }
        } else {
            break
        }
    }

    if remainingQty > 0 {
        order.Quantity = remainingQty
        return orderBook.AddLimitOrder(order)
    } else {
        delete(orderBook.OrderMap, order.ID)
    }
    return nil
}
