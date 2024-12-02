package vm

import "github.com/ava-labs/hypersdk/vm"

const Namespace = "controller"

type Config struct {
	Enabled bool `json:"enabled"`
}

func NewDefaultConfig() Config {
	return Config{
		Enabled: true,
	}
}

func With() vm.Option {
	return vm.NewOption(Namespace, NewDefaultConfig(), func(v *vm.VM, config Config) error {
		if !config.Enabled {
			return nil
		}
		vm.WithVMAPIs(jsonRPCServerFactory{})(v)
		return nil
	})
}