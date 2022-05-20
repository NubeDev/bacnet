package local

import (
	"github.com/NubeDev/bacnet"
	"github.com/NubeDev/bacnet/btypes"
)

type Discover struct {
	Name                      string
	MaxApdu                   uint32
	VendorName                string
	Segmentation              uint32
	ProtocolServicesSupported *btypes.BitString
}

//_, dd := data.ToBitString(out)

func (device *Device) DeviceDiscoverObjects(deviceID btypes.ObjectInstance) (resp *Discover, err error) {
	resp = &Discover{}
	obj := &Object{
		ObjectID:   deviceID,
		ObjectType: btypes.TypeDeviceType,
		Prop:       btypes.PropObjectName,
		ArrayIndex: bacnet.ArrayAll,
	}
	props := []btypes.PropertyType{btypes.PropObjectName, btypes.PropMaxAPDU, btypes.PropVendorName, btypes.PropSegmentationSupported, btypes.ProtocolServicesSupported}
	for _, prop := range props {
		obj.Prop = prop
		read, _ := device.Read(obj)
		switch prop {
		case btypes.PropObjectName:
			resp.Name = device.toStr(read)
		case btypes.PropMaxAPDU:
			resp.MaxApdu = device.toUint32(read)
		case btypes.PropVendorName:
			resp.VendorName = device.toStr(read)
		case btypes.PropSegmentationSupported:
			resp.Segmentation = device.toUint32(read)
		case btypes.ProtocolServicesSupported:
			resp.ProtocolServicesSupported = device.ToBitString(read)
		}
	}

	return resp, nil
}
