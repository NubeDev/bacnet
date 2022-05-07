package btypes

import (
	"fmt"
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
	ID           ObjectID
	MaxApdu      uint32 //maxApduLengthAccepted	62
	Segmentation Enumerated
	Vendor       uint32
	Addr         Address
	Objects      ObjectMap
	SupportsRPM  bool //support read prob multiple
	SupportsWPM  bool //support read prob multiple
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
