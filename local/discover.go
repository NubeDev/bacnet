package local

import (
	"fmt"
	"github.com/NubeDev/bacnet"
	"github.com/NubeDev/bacnet/btypes"
)

type Discover struct {
	Name         string
	MaxApdu      uint32
	VendorName   string
	Segmentation uint32
}

func (device *Device) DeviceDiscover(objectId btypes.ObjectInstance) (resp *Discover, err error) {
	resp = &Discover{}
	obj := &Object{
		ObjectID:   objectId,
		ObjectType: btypes.TypeDeviceType,
		Prop:       btypes.PropObjectName,
		ArrayIndex: bacnet.ArrayAll,
	}
	props := []btypes.PropertyType{btypes.PropObjectName, btypes.PropMaxAPDU, btypes.PropVendorName, btypes.PropSegmentationSupported}
	for _, prop := range props {
		obj.Prop = prop
		read, _ := device.Read(obj)
		switch prop {
		case btypes.PropObjectName:
			resp.Name = device.toStr(read)
		case btypes.PropMaxAPDU:
			out, ok := read.Object.Properties[0].Data.(uint32)
			fmt.Println(out, ok, "PropMaxAPDU")
			resp.MaxApdu = device.toUint32(read)
		case btypes.PropVendorName:
			resp.VendorName = device.toStr(read)
		case btypes.PropSegmentationSupported:
			fmt.Println(read, "PropSegmentationSupported")
			resp.Segmentation = device.toUint32(read)
		}
	}
	return resp, nil
}
