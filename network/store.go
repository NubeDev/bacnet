package network

import (
	"github.com/NubeDev/bacnet"
	"github.com/NubeDev/bacnet/helpers/store"
)

var cache *store.Handler

//UpdateClient updated a cached
func (local *Local) UpdateClient(storeID string) error {
	//first close the client
	local.ClientClose()

	cb := &bacnet.ClientBuilder{
		Interface:  local.Interface,
		Ip:         local.Ip,
		Port:       local.Port,
		SubnetCIDR: local.SubnetCIDR,
	}

	bc, err := bacnet.NewClient(cb)
	if err != nil {
		return err
	}
	local.bacnet = bc
	cache.Set(storeID, local, -1)
	return nil
}
