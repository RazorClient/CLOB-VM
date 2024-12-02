// CLOB/genesis/rules.go

package genesis

import (
	"github.com/ava-labs/hypersdk/chain"
	// Import other necessary packages if needed
)

var _ chain.Rules = (*Rules)(nil)

// Rules defines the interface for accessing genesis configuration parameters.
// It wraps the Genesis struct and provides getter methods.
type Rules struct {
	g *Genesis
}

// GetMaxBlockTxs returns the maximum number of transactions allowed per block.
func (r *Rules) GetMaxBlockTxs() int {
	return r.g.MaxBlockTxs
}

// GetMaxBlockUnits returns the maximum number of units allowed per block.
func (r *Rules) GetMaxBlockUnits() uint64 {
	return r.g.MaxBlockUnits
}

// GetBaseUnits returns the base units used in transactions.
func (r *Rules) GetBaseUnits() uint64 {
	return r.g.BaseUnits
}

// GetValidityWindow returns the validity window for transactions.
func (r *Rules) GetValidityWindow() int64 {
	return r.g.ValidityWindow
}

// GetMinUnitPrice returns the minimum unit price for orders.
func (r *Rules) GetMinUnitPrice() uint64 {
	return r.g.MinUnitPrice
}

// GetUnitPriceChangeDenominator returns the denominator for unit price changes.
func (r *Rules) GetUnitPriceChangeDenominator() uint64 {
	return r.g.UnitPriceChangeDenominator
}

// GetWindowTargetUnits returns the target units per window.
func (r *Rules) GetWindowTargetUnits() uint64 {
	return r.g.WindowTargetUnits
}

// GetMinBlockCost returns the minimum cost per block.
func (r *Rules) GetMinBlockCost() uint64 {
	return r.g.MinBlockCost
}

// GetBlockCostChangeDenominator returns the denominator for block cost changes.
func (r *Rules) GetBlockCostChangeDenominator() uint64 {
	return r.g.BlockCostChangeDenominator
}

// GetWindowTargetBlocks returns the target number of blocks per window.
func (r *Rules) GetWindowTargetBlocks() uint64 {
	return r.g.WindowTargetBlocks
}

// GetHRP returns the Human-Readable Part for addresses.
func (r *Rules) GetHRP() string {
	return r.g.HRP
}

// FetchCustom allows fetching custom configuration parameters.
// Modify this method based on your project's requirements.
func (r *Rules) FetchCustom(key string) (any, bool) {
	// Implement fetching of custom parameters if needed.
	return nil, false
}

// Implement any other methods required by the chain.Rules interface.
// If certain methods are not applicable to your project, you can leave them unimplemented or return default values.

// Example of implementing methods from chain.Rules interface that are not applicable:
func (r *Rules) GetWarpBaseFee() uint64 {
	// Not applicable for Order Book Matching Engine
	return 0
}

func (r *Rules) GetWarpConfig(sourceChainID string) (bool, uint64, uint64) {
	// Not applicable for Order Book Matching Engine
	return false, 0, 0
}

func (r *Rules) GetWarpFeePerSigner() uint64 {
	// Not applicable for Order Book Matching Engine
	return 0
}

// Add the Rules getter to the Genesis struct.
func (g *Genesis) Rules(int64) *Rules {
	return &Rules{g}
}
