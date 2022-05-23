package network

import (
	"github.com/NubeDev/bacnet"
	"github.com/NubeDev/bacnet/helpers/store"
)

type Local struct {
	Interface  string
	Ip         string
	Port       int
	SubnetCIDR int
	bacnet     bacnet.Client
	StoreID    string
}

// New returns a new instance of bacnet network
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

	cache = store.Init()
	local.bacnet = bc
	cache.Set("1", local, -1)
	return local, nil
}

func (local *Local) ClientClose() {
	local.bacnet.Close()
}

func (local *Local) ClientRun() {
	local.bacnet.ClientRun()
}

func (local *Local) store() {

}
