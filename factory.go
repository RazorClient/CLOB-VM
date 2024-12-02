
package vm

import (
	"github.com/ava-labs/avalanchego/utils/logging"
	"github.com/ava-labs/avalanchego/vms"

	// controller
	CLOB/controller
)

// Factory is a struct that implements the vms.Factory interface.
// It is responsible for creating new instances of the controller.
var _ vms.Factory = &Factory{}

// Factory struct definition.
type Factory struct{}

// New creates a new instance of the controller and returns it along with any error encountered.
// It takes a logging.Logger as an argument for logging purposes.
func (*Factory) New(logger logging.Logger) (interface{}, error) {
	return controller.New(), nil // Ensure controller.New() handles errors appropriately
}