// CLOB/genesis/errors.go

package genesis

import "errors"

var (
	// ErrInvalidGenesisConfig is returned when the genesis configuration is invalid
	ErrInvalidGenesisConfig = errors.New("invalid genesis configuration")

	// ErrInitialOrderSetupFailed is returned when setting up initial orders fails
	ErrInitialOrderSetupFailed = errors.New("failed to set up initial orders")

	// ErrDuplicateInitialOrder is returned when an initial order ID is duplicated
	ErrDuplicateInitialOrder = errors.New("duplicate initial order ID found")
)
