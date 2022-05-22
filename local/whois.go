package local

import (
	"github.com/NubeDev/bacnet"
	"github.com/NubeDev/bacnet/btypes"
)

func (local *Local) Whois(options *bacnet.WhoIsOpts) ([]btypes.Device, error) {
	resp, err := local.bacnet.WhoIs(options)
	return resp, err
}
