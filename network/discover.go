package network

import (
	"github.com/NubeDev/bacnet"
	"github.com/NubeDev/bacnet/btypes"
	log "github.com/sirupsen/logrus"
)

type DevicePoints struct {
	Name                      string            `json:"name"`
	MaxApdu                   uint32            `json:"max_apdu"`
	VendorName                string            `json:"vendor_name"`
	Segmentation              uint32            `json:"segmentation"`
	ProtocolServicesSupported *btypes.BitString `json:"protocol_services_supported"`
}

//GetDevicePoints build a device points list
//first read device and see what it supports and get the name and so on
//try and get the object list if it's an error then loop through the arrayIndex to build the object list
//with the object list do a point's discovery, get the name, units and so on
func (device *Device) GetDevicePoints(deviceID btypes.ObjectInstance) (resp []*PointDetails, err error) {
	resp = []*PointDetails{}
	list, err := device.DeviceObjects(deviceID, true)
	if err != nil {
		return nil, err
	}
	pntDetails := &Point{}
	for _, obj := range list {

		if obj.Type != 8 {
			pntDetails = &Point{
				ObjectID:   obj.Instance,
				ObjectType: obj.Type,
			}
			details, _ := device.PointDetails(pntDetails)
			resp = append(resp, details)
		}
	}
	return resp, nil

}

type DeviceDetails struct {
	Name                      string            `json:"name"`
	MaxApdu                   uint32            `json:"max_apdu"`
	VendorName                string            `json:"vendor_name"`
	Segmentation              uint32            `json:"segmentation"`
	ProtocolServicesSupported *btypes.BitString `json:"protocol_services_supported"`
}

//GetDeviceDetails get the device name, max adpu and so on
//first read device and see what it supports and get the name and so on
//try and get the object list if it's an error then loop through the arrayIndex to build the object list
func (device *Device) GetDeviceDetails(deviceID btypes.ObjectInstance) (resp *DeviceDetails, err error) {
	resp = &DeviceDetails{}
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
	log.Println("bacnet-device name:", resp.Name)
	log.Println("bacnet-device vendor-name:", resp.VendorName)
	return resp, nil
}
