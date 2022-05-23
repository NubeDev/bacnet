package network

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
	MacMSTP       int
	MaxApdu       uint32
	Segmentation  uint32
	Dev           btypes.Device
	bacnet        bacnet.Client
}

// NewDevice returns a new instance of ta bacnet device
func NewDevice(bacnetDevice *Local, device *Device) (*Device, error) {
	dev := &btypes.Device{
		Ip:            device.Ip,
		DeviceID:      device.DeviceID,
		NetworkNumber: device.NetworkNumber,
		MacMSTP:       device.MacMSTP,
		MaxApdu:       device.MaxApdu,
		Segmentation:  btypes.Enumerated(device.Segmentation),
	}

	dev, err := btypes.NewDevice(dev)
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
