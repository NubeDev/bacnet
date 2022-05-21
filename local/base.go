package local

import (
	"github.com/NubeDev/bacnet"
	"github.com/NubeDev/bacnet/btypes"
)

type Local struct {
	Interface  string
	Ip         string
	Port       int
	SubnetCIDR int
	bacnet     bacnet.Client
}

// New returns a new instance of the nube common apis
func New(local *Local) (*Local, error) {
	cb := &bacnet.ClientBuilder{
		Interface:  local.Interface,
		Ip:         local.Ip,
		Port:       local.Port,
		SubnetCIDR: local.SubnetCIDR,
	}
	bc, err := bacnet.NewClient(cb)
	if err != nil {
		return nil, err
	}
	local.bacnet = bc
	return local, nil
}

func (local *Local) Whois(options *bacnet.WhoIsOpts) ([]btypes.Device, error) {

	return local.bacnet.WhoIs(options)
}

func (local *Local) ClientClose() {
	local.bacnet.Close()
}

func (local *Local) ClientRun() {
	local.bacnet.ClientRun()
}
