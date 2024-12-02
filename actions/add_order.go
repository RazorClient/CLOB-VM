// CLOB/actions/add_order.go

package actions

import (
	"fmt"

	"CLOB/genesis"
	"CLOB/storage"
)

type AddOrderAction struct {
	Order *storage.Order
}

func (a *AddOrderAction) Execute(vm *vm.MatchingEngineVM) error {
	orderBook := vm.OrderBook
	rules := vm.Rules

	// Example usage of Rules
	maxBlockTxs := rules.GetMaxBlockTxs()
	if orderBook.CurrentBlockTxs >= maxBlockTxs {
		return fmt.Errorf("cannot add order: max block transactions (%d) reached", maxBlockTxs)
	}

	// Proceed to add the order
	if err := orderBook.AddOrder(orderBook.Context, a.Order); err != nil {
		return fmt.Errorf("failed to add order: %w", err)
	}

	return nil
}
