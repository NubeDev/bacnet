package btypes

import (
	"fmt"
	ip2bytes "github.com/NubeDev/bacnet/helpers/ipbytes"
	"github.com/NubeDev/bacnet/helpers/validation"
)

type Enumerated uint32

type IAm struct {
	ID           ObjectID
	MaxApdu      uint32
	Segmentation Enumerated
	Vendor       uint32
	Addr         Address
}

type Device struct {
	ID            ObjectID
	DeviceID      int
	Ip            string
	Port          int
	NetworkNumber int
	MacMSTP       int
	MaxApdu       uint32 //maxApduLengthAccepted	62
	Segmentation  Enumerated
	Vendor        uint32
	Addr          Address
	Objects       ObjectMap
	SupportsRPM   bool //support read prob multiple
	SupportsWPM   bool //support read prob multiple
}

// NewDevice returns a new instance of ta bacnet device
func NewDevice(device *Device) (*Device, error) {

	port := device.Port
	//check ip
	ok := validation.ValidIP(device.Ip)
	if !ok {
		fmt.Println("fail ip")
	}
	//check port
	if port == 0 {
		port = 0xBAC0
	}
	ok = validation.ValidPort(port)
	if !ok {
		fmt.Println("fail port")
	}

	ip, err := ip2bytes.New(device.Ip, uint16(port))
	if err != nil {
		fmt.Println("fail ip2bytes")
		return nil, err
	}
	addr := Address{
		Net: uint16(device.NetworkNumber),
		Mac: ip,
		Adr: []uint8{uint8(device.MacMSTP)},
	}
	object := ObjectID{
		Type:     DeviceType,
		Instance: device.ID.Instance,
	}
	device.ID = object
	device.Addr = addr
	return device, nil
}

// ObjectSlice returns all the objects in the device as a slice (not thread-safe)
func (dev *Device) ObjectSlice() []Object {
	var objs []Object
	for _, objMap := range dev.Objects {
		for _, o := range objMap {
			objs = append(objs, o)
		}
	}
	return objs
}

//CheckADPU device max ADPU len (mstp can be > 480, and IP > 1476)
func (dev *Device) CheckADPU() error {
	errMsg := "device.CheckADPU() incorrect ADPU size:"
	size := dev.MaxApdu
	if size == 0 {
		return fmt.Errorf("%s %d", errMsg, size)
	}
	return nil
}
