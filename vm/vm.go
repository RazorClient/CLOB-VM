// CLOB/vm/vm.go

package vm

import (
	"context"
	"fmt"

	"CLOB/actions"
	"CLOB/genesis"
	"CLOB/storage"

	"github.com/ava-labs/hypersdk/chain"
	"github.com/ava-labs/hypersdk/vm"
)

type MatchingEngineVM struct {
	OrderBook *storage.OrderBook
	Rules     *genesis.Rules
}

// NewMatchingEngineVM creates a new instance of the VM with genesis configuration
func NewMatchingEngineVM(genesisConfig []byte) (*MatchingEngineVM, error) {
	// Initialize Genesis
	genesisInstance, err := genesis.New(genesisConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize genesis: %w", err)
	}

	// Initialize Rules
	rules := genesisInstance.Rules(0) // Pass appropriate parameter if needed

	// Initialize OrderBook
	orderBook := storage.NewOrderBook()

	// Load Genesis into OrderBook
	if err := genesisInstance.Load(context.Background(), orderBook); err != nil {
		return nil, fmt.Errorf("failed to load genesis into order book: %w", err)
	}

	return &MatchingEngineVM{
		OrderBook: orderBook,
		Rules:     rules,
	}, nil
}

// ExecuteAction executes a given action on the VM
func (vm *MatchingEngineVM) ExecuteAction(action actions.Action) error {
	if err := action.Execute(vm); err != nil {
		// Wrap or handle the error as needed
		return fmt.Errorf("failed to execute action: %w", err)
	}
	return nil
}

// GetOrderBook returns the VM's order book
func (vm *MatchingEngineVM) GetOrderBook() *storage.OrderBook {
	return vm.OrderBook
}

// GetRules returns the VM's rules
func (vm *MatchingEngineVM) GetRules() *genesis.Rules {
	return vm.Rules
}


func New(options ...vm.Option) (*vm.VM, error) {
	options = append(options, With()) // Add MorpheusVM API
	return defaultvm.New(
		consts.Version,
		genesis.DefaultGenesisFactory{},
		&storage.BalanceHandler{},
		metadata.NewDefaultManager(),
		ActionParser,
		AuthParser,
		OutputParser,
		auth.Engines(),
		options...,
	)
}