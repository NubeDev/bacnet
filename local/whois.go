package local

import (
	"github.com/NubeDev/bacnet"
	"github.com/NubeDev/bacnet/btypes"
)

func (local *Local) Whois(options *bacnet.WhoIsOpts) ([]btypes.Device, error) {
	return local.bacnet.WhoIs(options)
}
