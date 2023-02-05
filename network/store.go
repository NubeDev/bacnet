package network

import (
	"fmt"
	"github.com/NubeDev/bacnet"
	"github.com/NubeDev/bacnet/btypes"
	"github.com/NubeDev/bacnet/helpers/store"
	"github.com/pkg/errors"
)

//var memDb *store.Handler

type Store struct {
	BacStore *store.Handler
}

func NewStore() *Store {
	s := &Store{
		BacStore: store.Init(),
	}
	return s
}

//NewNetwork updated a cached
func (store *Store) NewNetwork(storeID, iface, ip string, port, subnet int) error {
	cb := &bacnet.ClientBuilder{
		Interface:  iface,
		Ip:         ip,
		Port:       port,
		SubnetCIDR: subnet,
	}
	bc, err := bacnet.NewClient(cb)
	if err != nil {
		return err
	}
	bacnetNet := &Network{
		Interface: iface,
		Port:      port,
		bacnet:    bc,
	}
	if store.BacStore != nil {
		store.BacStore.Set(storeID, bacnetNet, -1)
		return nil
	} else {
		return errors.New("new network, failed to set bacnet store, bacnet store is empty")
	}
}

func (store *Store) GetNetwork(uuid string) (*Network, error) {
	cli, ok := store.BacStore.Get(uuid)
	if !ok {
		return nil, errors.New(fmt.Sprintf("bacnet: no network found with uuid:%s", uuid))
	}
	parse := cli.(*Network)
	return parse, nil
}

//UpdateDevice updated a cached device
func (store *Store) UpdateDevice(storeID string, net *Network, device *Device) error {
	var err error
	dev := &btypes.Device{
		Ip:            device.Ip,
		DeviceID:      device.DeviceID,
		NetworkNumber: device.NetworkNumber,
		MacMSTP:       device.MacMSTP,
		MaxApdu:       device.MaxApdu,
		Segmentation:  btypes.Enumerated(device.Segmentation),
	}
	dev, err = btypes.NewDevice(dev)
	if err != nil {
		return err
	}
	if dev == nil {
		fmt.Println("dev is nil")
		return err
	}
	device.network = net.bacnet
	device.dev = *dev
	if store.BacStore != nil {
		store.BacStore.Set(storeID, device, -1)
	}
	return nil
}

func (store *Store) GetDevice(uuid string) (*Device, error) {
	cli, ok := store.BacStore.Get(uuid)
	if !ok {
		return nil, errors.New(fmt.Sprintf("bacnet: no device found with uuid:%s", uuid))
	}
	parse := cli.(*Device)
	return parse, nil
}
