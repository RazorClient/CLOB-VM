// CLOB/vm/vm_interfaces.go
package vm

import (
    "CLOB/actions"
    "CLOB/storage"
)

// VM defines the interface for the virtual machine
type VM interface {
    ExecuteAction(action actions.Action) error
    GetOrderBook() *storage.OrderBook
}
