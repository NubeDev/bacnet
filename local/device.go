package local

import (
	"fmt"
	"github.com/NubeDev/bacnet"
	"github.com/NubeDev/bacnet/btypes"
)

type Device struct {
	Ip            string
	Port          int
	DeviceID      int
	NetworkNumber int
	MSTPMac       int
	Dev           btypes.Device
	bacnet        bacnet.Client
}

// NewDevice returns a new instance of ta bacnet device
func NewDevice(bacnetDevice *Local, device *Device) (*Device, error) {
	dev, err := btypes.NewDevice(&btypes.Device{Ip: device.Ip, DeviceID: device.DeviceID})
	if err != nil {
		return nil, err
	}
	if dev == nil {
		fmt.Println("dev is nil")
		return nil, err
	}
	device.bacnet = bacnetDevice.bacnet
	device.Dev = *dev
	return device, nil
}
