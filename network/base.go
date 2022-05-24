package network

import (
	"github.com/NubeDev/bacnet"
)

type Network struct {
	Interface  string
	Ip         string
	Port       int
	SubnetCIDR int
	StoreID    string
	bacnet     bacnet.Client
}

// New returns a new instance of bacnet network
func New(net *Network) (*Network, error) {
	cb := &bacnet.ClientBuilder{
		Interface:  net.Interface,
		Ip:         net.Ip,
		Port:       net.Port,
		SubnetCIDR: net.SubnetCIDR,
	}

	bc, err := bacnet.NewClient(cb)
	if err != nil {
		return nil, err
	}

	net.bacnet = bc
	if BacStore != nil {
		BacStore.Set(net.StoreID, net, -1)
	}
	return net, nil
}

func (net *Network) NetworkClose() {
	if net.bacnet != nil {
		net.bacnet.Close()
	}

}

func (net *Network) NetworkRun() {
	if net.bacnet != nil {
		net.bacnet.ClientRun()
	}

}

func (net *Network) store() {

}
