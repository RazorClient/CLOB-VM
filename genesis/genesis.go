// CLOB/genesis/genesis.go

package genesis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"CLOB/storage"
	"CLOB/utils"

	"github.com/rafael-abuawad/tokenvm/consts" // Adjust the import path as needed
)

// Ensure Genesis implements any required interfaces (if applicable)
// var _ SomeInterface = (*Genesis)(nil)

// CustomInitialOrder represents an initial order to be loaded into the order book
type CustomInitialOrder struct {
	ID        string  `json:"id"`
	Side      string  `json:"side"`      // "buy" or "sell"
	Price     float64 `json:"price"`     // 0 for market orders
	Quantity  float64 `json:"quantity"`
	Timestamp string  `json:"timestamp"` // ISO8601 format
	OrderType string  `json:"order_type"` // "limit" or "market"
}

// Genesis defines the structure for initializing the order book
type Genesis struct {
	// Configuration Parameters
	MaxBlockTxs   int `json:"max_block_txs"`
	MaxBlockUnits int `json:"max_block_units"`

	// Initial Orders
	InitialOrders []CustomInitialOrder `json:"initial_orders"`
}

// Default returns a Genesis instance with default configurations
func Default() *Genesis {
	return &Genesis{
		MaxBlockTxs:   1000,
		MaxBlockUnits: 1000000,
		InitialOrders: []CustomInitialOrder{
			{
				ID:        "init_buy_1",
				Side:      "buy",
				Price:     100.0,
				Quantity:  50,
				Timestamp: time.Now().UTC().Format(time.RFC3339),
				OrderType: "limit",
			},
			{
				ID:        "init_sell_1",
				Side:      "sell",
				Price:     101.0,
				Quantity:  50,
				Timestamp: time.Now().UTC().Format(time.RFC3339),
				OrderType: "limit",
			},
			// Add more initial orders as needed
		},
	}
}

// New creates a new Genesis instance from JSON configuration
func New(configBytes []byte) (*Genesis, error) {
	genesis := Default()
	if len(configBytes) > 0 {
		if err := json.Unmarshal(configBytes, genesis); err != nil {
			return nil, fmt.Errorf("failed to unmarshal genesis config: %w", err)
		}
	}

	// Validate Genesis Configuration
	if genesis.MaxBlockTxs <= 0 {
		return nil, fmt.Errorf("%w: MaxBlockTxs must be positive", ErrInvalidGenesisConfig)
	}
	if genesis.MaxBlockUnits <= 0 {
		return nil, fmt.Errorf("%w: MaxBlockUnits must be positive", ErrInvalidGenesisConfig)
	}

	// Validate Initial Orders
	orderIDs := make(map[string]struct{})
	for _, order := range genesis.InitialOrders {
		// Validate order side
		if order.Side != "buy" && order.Side != "sell" {
			return nil, fmt.Errorf("%w: invalid side '%s' for order ID '%s'", ErrInvalidGenesisConfig, order.Side, order.ID)
		}

		// Validate order type
		if order.OrderType != "limit" && order.OrderType != "market" {
			return nil, fmt.Errorf("%w: invalid order type '%s' for order ID '%s'", ErrInvalidGenesisConfig, order.OrderType, order.ID)
		}

		// Check for duplicate order IDs
		if _, exists := orderIDs[order.ID]; exists {
			return nil, fmt.Errorf("%w: duplicate order ID '%s'", ErrDuplicateInitialOrder, order.ID)
		}
		orderIDs[order.ID] = struct{}{}
	}

	return genesis, nil
}

// Load initializes the order book with the genesis configuration
func (g *Genesis) Load(ctx context.Context, db storage.Database) error {
	// Set up any required configurations in the storage layer
	// For example, setting block parameters if applicable
	// This depends on how your storage layer utilizes these parameters

	// Initialize initial orders
	for _, order := range g.InitialOrders {
		parsedTimestamp, err := time.Parse(time.RFC3339, order.Timestamp)
		if err != nil {
			return fmt.Errorf("invalid timestamp for order ID '%s': %w", order.ID, err)
		}

		// Create Order struct
		storageOrder := &storage.Order{
			ID:        order.ID,
			Side:      storage.Side(order.Side),
			Price:     order.Price,
			Quantity:  order.Quantity,
			Timestamp: parsedTimestamp,
			OrderType: storage.OrderType(order.OrderType),
		}

		// Add order to the order book
		if err := db.AddOrder(ctx, storageOrder); err != nil {
			return fmt.Errorf("failed to add initial order ID '%s': %w", order.ID, err)
		}
	}

	return nil
}
