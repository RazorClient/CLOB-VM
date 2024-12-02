
/*
Package controller provides the implementation of a Controller that manages
the interactions between the virtual machine (VM), storage, and genesis
configuration for a blockchain application. It implements the vm.Controller
interface, ensuring that it adheres to the expected behavior for managing
blockchain operations.
*/

package controller

import (
	"context"
	"fmt"

	ametrics "github.com/ava-labs/avalanchego/api/metrics"
	"github.com/ava-labs/avalanchego/database"
	"github.com/ava-labs/avalanchego/snow"
	"github.com/ava-labs/avalanchego/snow/engine/common"
	"github.com/ava-labs/hypersdk/builder"
	"github.com/ava-labs/hypersdk/chain"
	"github.com/ava-labs/hypersdk/gossiper"
	"github.com/ava-labs/hypersdk/pebble"
	"github.com/ava-labs/hypersdk/utils"
	"github.com/ava-labs/hypersdk/vm"
	"go.uber.org/zap"

	"CLOB/actions"
	"CLOB/config"
	"CLOB/consts"
	"CLOB/genesis"
	"CLOB/storage"
	"CLOB/version"
)

// Ensure Controller implements the vm.Controller interface
var _ vm.Controller = (*Controller)(nil)

// Controller manages the VM and its interactions with the storage and genesis
// configuration. It is responsible for initializing the VM, loading
// configurations, handling accepted and rejected blocks, and managing
// state transitions.
type Controller struct {
	inner        *vm.VM          // The underlying VM instance
	snowCtx      *snow.Context    // Context for snow operations
	genesis      *genesis.Genesis // Genesis configuration
	config       *config.Config    // Configuration settings
	stateManager *StateManager     // Manages the state of the chain
	metrics      *Metrics          // Metrics for tracking performance
	metaDB       database.Database  // Database for metadata storage
}

// New creates a new instance of the VM with the Controller. It initializes
// the Controller and returns a pointer to the VM instance.
func New() *vm.VM {
	return vm.New(&Controller{}, version.Version)
}

// Initialize sets up the Controller by loading configurations, genesis,
// and initializing databases. It returns various components needed for
// the VM's operation, including configuration, genesis information,
// builders, gossipers, and databases. It also handles errors during
// initialization.
func (c *Controller) Initialize(
	inner *vm.VM,
	snowCtx *snow.Context,
	gatherer ametrics.MultiGatherer,
	genesisBytes []byte,
	upgradeBytes []byte,
	configBytes []byte,
) (
	vm.Config,          // Configuration for the VM
	vm.Genesis,        // Genesis information
	builder.Builder,    // Builder for block creation
	gossiper.Gossiper,  // Gossiper for block propagation
	database.Database,  // Database for block storage
	database.Database,  // Database for state storage
	vm.Handlers,       // Handlers for API endpoints
	chain.ActionRegistry, // Registry for actions
	chain.AuthRegistry,    // Registry for authentication
	error,             // Error if initialization fails
) {
	// Set the inner VM and snow context
	c.inner = inner
	c.snowCtx = snowCtx
	c.stateManager = &StateManager{}

	// Initialize metrics
	var err error
	c.metrics, err = newMetrics(gatherer)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	// Load configuration
	c.config, err = config.New(c.snowCtx.NodeID, configBytes)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	c.snowCtx.Log.SetLevel(c.config.GetLogLevel())
	snowCtx.Log.Info("loaded config", zap.Any("contents", c.config))

	// Load genesis
	c.genesis, err = genesis.New(genesisBytes)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf(
			"unable to read genesis: %w",
			err,
		)
	}
	snowCtx.Log.Info("loaded genesis", zap.Any("genesis", c.genesis))

	// Initialize databases
	blockPath, err := utils.InitSubDirectory(snowCtx.ChainDataDir, "block")
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	cfg := pebble.NewDefaultConfig()
	blockDB, err := pebble.New(blockPath, cfg)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	statePath, err := utils.InitSubDirectory(snowCtx.ChainDataDir, "state")
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	stateDB, err := pebble.New(statePath, cfg)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	metaPath, err := utils.InitSubDirectory(snowCtx.ChainDataDir, "metadata")
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	c.metaDB, err = pebble.New(metaPath, cfg)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	// Initialize handlers
	apis := map[string]*common.HTTPHandler{}
	endpoint, err := utils.NewHandler(consts.Name, &Handler{inner.Handler(), c})
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	apis[vm.Endpoint] = endpoint

	// Create builder and gossiper
	var (
		build  builder.Builder
		gossip gossiper.Gossiper
	)
	if c.config.TestMode {
		c.inner.Logger().Info("running build and gossip in test mode")
		build = builder.NewManual(inner)
		gossip = gossiper.NewManual(inner)
	} else {
		build = builder.NewTime(inner, builder.DefaultTimeConfig())
		gcfg := gossiper.DefaultProposerConfig()
		gcfg.BuildProposerDiff = 1 // Don't gossip if producing the next block
		gossip = gossiper.NewProposer(inner, gcfg)
	}

	return c.config, c.genesis, build, gossip, blockDB, stateDB, apis, consts.ActionRegistry, consts.AuthRegistry, nil
}

// Rules returns the chain.Rules interface based on the genesis configuration.
// It provides the rules that govern the blockchain's behavior.
func (c *Controller) Rules(t int64) chain.Rules {
	return c.genesis.Rules(t)
}

// StateManager returns the chain.StateManager interface. It provides access
// to the state management functionalities of the blockchain.
func (c *Controller) StateManager() chain.StateManager {
	return c.stateManager
}

// Accepted processes accepted blocks and stores transaction results. It
// iterates through the transactions in the block, storing their results
// in the metadata database and updating metrics based on the transaction
// actions.
func (c *Controller) Accepted(ctx context.Context, blk *chain.StatelessBlock) error {
	batch := c.metaDB.NewBatch()
	defer batch.Reset()

	results := blk.Results()
	for i, tx := range blk.Txs {
		result := results[i]
		err := storage.StoreTransaction(
			ctx,
			batch,
			tx.ID(),
			blk.GetTimestamp(),
			result.Success,
			result.Units,
		)
		if err != nil {
			return err
		}
		if result.Success {
			switch tx.Action.(type) {
			case *actions.AddOrderAction:
				c.metrics.addOrder.Inc()
			case *actions.CancelOrderAction:
				c.metrics.cancelOrder.Inc()
			case *actions.MatchOrderAction:
				c.metrics.matchOrder.Inc()
			}
		}
	}
	return batch.Write()
}

// Rejected handles rejected blocks. This implementation is a no-op,
// meaning it does not perform any actions when a block is rejected.
func (*Controller) Rejected(context.Context, *chain.StatelessBlock) error {
	return nil
}

// Shutdown gracefully shuts down the controller. This implementation
// does not close any databases provided during initialization, as the
// VM is responsible for closing them.
func (*Controller) Shutdown(context.Context) error {
	// Do not close any databases provided during initialization. The VM will
	// close any databases you're provided.
	return nil
}
