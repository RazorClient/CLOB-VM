// CLOB/actions/cancel_order.go
package actions

import "CLOB/storage"

/
// This struct encapsulates the information needed to cancel an order in the order book.
// It includes the `OrderID` field, which specifies the ID of the order to be canceled.
type CancelOrderAction struct {
    OrderID string // The unique identifier of the order to cancel.
}

// The `Execute` method implements the logic to cancel an order in the order book.
// It performs the following steps:
// 1. Retrieve the order book from the VMContext.
// 2. Look up the order in the `OrderMap` using the provided `OrderID`.
// 3. If the order exists, invoke the `CancelOrder` method to remove it from the order book.
// 4. Return an error if the order does not exist.
//
// Parameters:
// - vm (VMContext): The virtual machine context that provides access to the order book.
//
// Returns:
// - error: An error if the order cannot be found or if there are issues during cancellation.
func (a *CancelOrderAction) Execute(vm VMContext) error {
    
    // The `VMContext` provides access to the shared state of the order book.
    // This allows the function to interact with the current orders.
    orderBook := vm.GetOrderBook()

    // The `OrderMap` is a mapping of order IDs to their corresponding order objects.
    // Check if the order exists in the map using the provided `OrderID`.
    order, exists := orderBook.OrderMap[a.OrderID]
    if !exists {
        // If the order does not exist, return a predefined error indicating that
        // the order was not found.
        return storage.ErrOrderNotFound
    }

    
    // If the order exists, call the `CancelOrder` method on the order book to
    // remove the order and perform any necessary cleanup.
    return orderBook.CancelOrder(order)
}
