package consts

import (
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/version"
	"github.com/ava-labs/avalanchego/vms/platformvm/warp"
	"github.com/ava-labs/hypersdk/chain"
	"github.com/ava-labs/hypersdk/codec"
	"github.com/ava-labs/hypersdk/consts"
)

const (
	HRP="ClbVM"
	Name = "CLOB-vm"
)

var ID ids.ID

func init() {
	b := make([]byte, ids.IDLen) // Create a byte slice of length `ids.IDLen`.
	copy(b, []byte(Name))       // Copy the VM name into the byte slice.
	vmID, err := ids.ToID(b)    // Convert the byte slice to an `ids.ID`.
	if err != nil {
		panic(err)             // If an error occurs, panic to signal a critical issue.
	}
	ID = vmID                   // Assign the computed `ids.ID` to the `ID` variable.
}


var Version = &version.Semantic{
	Major: 0,
	Minor: 0,
	Patch: 1,
}


