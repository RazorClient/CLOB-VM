// CLOB/actions/action_interfaces.go
package actions

import "CLOB/storage"

// Action defines the interface for all actions
type Action interface {
    Execute(vm VMContext) error
}

// VMContext provides access to the VM's state for actions
type VMContext interface {
    GetOrderBook() *storage.OrderBook
}
