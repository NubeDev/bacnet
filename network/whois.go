package network

import (
	"github.com/NubeDev/bacnet"
	"github.com/NubeDev/bacnet/btypes"
)

func (device *Device) Whois(options *bacnet.WhoIsOpts) ([]btypes.Device, error) {
	go device.network.ClientRun()
	resp, err := device.network.WhoIs(options)
	return resp, err
}

func (net *Network) Whois(options *bacnet.WhoIsOpts) ([]btypes.Device, error) {
	go net.NetworkRun()
	resp, err := net.bacnet.WhoIs(options)
	return resp, err
}
