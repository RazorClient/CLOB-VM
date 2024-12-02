// controller/registry.go

package controller

import (
	"github.com/ava-labs/avalanchego/utils/wrappers"
	"github.com/ava-labs/hypersdk/chain"
	"github.com/ava-labs/hypersdk/codec"

	"CLOB/actions"
	"CLOB/auth"
	"CLOB/consts"
)

// init initializes the Action and Auth registries in the consts package.
// It registers various action types and authentication methods, collecting
// any errors that occur during the registration process. If any errors
// are encountered, the application will panic to prevent incorrect startup.
func init() {
	// Initialize the Action and Auth registries in the consts package
	consts.ActionRegistry = codec.NewTypeParser[chain.Action, *codec.Message]()
	consts.AuthRegistry = codec.NewTypeParser[chain.Auth, *codec.Message]()

	// Collect errors during registration
	errs := &wrappers.Errs{}
	errs.Add(
		// Register Actions
		consts.ActionRegistry.Register(&actions.AddOrder{}, actions.UnmarshalAddOrder, false),
		consts.ActionRegistry.Register(&actions.CancelOrder{}, actions.UnmarshalCancelOrder, false),
		consts.ActionRegistry.Register(&actions.MatchOrder{}, actions.UnmarshalMatchOrder, false),

		// Register Auth Types
		consts.AuthRegistry.Register(&auth.ED25519{}, auth.UnmarshalED25519, false),
	)
	
	// If any errors occurred during registration, panic to prevent the application from starting incorrectly
	if errs.Errored() {
		panic(errs.Err)
	}
}
