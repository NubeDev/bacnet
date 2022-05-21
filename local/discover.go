package local

import (
	"github.com/NubeDev/bacnet"
	"github.com/NubeDev/bacnet/btypes"
)

type Discover struct {
	Name                      string            `json:"name"`
	MaxApdu                   uint32            `json:"max_apdu"`
	VendorName                string            `json:"vendor_name"`
	Segmentation              uint32            `json:"segmentation"`
	ProtocolServicesSupported *btypes.BitString `json:"protocol_services_supported"`
}

//DeviceDiscover get the device name, max adpu and so on
//first read device and see what it supports and get the name and so on
//try and get the object list if it's an error then loop through the arrayIndex to build the object list
//with the object list do a point's discovery, get the name, units and so on
func (device *Device) DeviceDiscover(deviceID btypes.ObjectInstance) (resp *Discover, err error) {
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
