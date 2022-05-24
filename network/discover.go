package network

import (
	"fmt"
	"github.com/NubeDev/bacnet"
	"github.com/NubeDev/bacnet/btypes"
	pprint "github.com/NubeDev/bacnet/helpers/print"
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
	fmt.Println("GetDeviceDetails()")
	pprint.PrintJOSN(device)
	fmt.Println("GetDeviceDetails()")
	props := []btypes.PropertyType{btypes.PropObjectName, btypes.PropMaxAPDU, btypes.PropVendorName, btypes.PropSegmentationSupported}
	for i, prop := range props {
		fmt.Println(i, "Loop Props", prop, " deviceID:", deviceID)
		obj.Prop = prop
		read, err := device.Read(obj)
		if err != nil {
			log.Errorln("bacnet-master-GetDeviceDetails()", err.Error())
		}
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
	log.Infoln("bacnet-device name:", resp.Name)
	log.Infoln("bacnet-device vendor-name:", resp.VendorName)
	return resp, nil
}

func (device *Device) DeviceDiscover() error {
	options := &bacnet.WhoIsOpts{
		Low:             0,
		High:            0,
		GlobalBroadcast: true,
		NetworkNumber:   0,
	}
	whois, err := device.Whois(options)
	if err != nil {
		return err
	}
	fmt.Println("--------devices------------found device count:", len(whois))
	pprint.PrintJOSN(whois)
	fmt.Println("--------devices------------")
	for _, dev := range whois {
		if len(dev.Addr.Adr) > 0 {
			device.MacMSTP = int(dev.Addr.Adr[0])
		}
		host, _ := dev.Addr.UDPAddr()
		device.DeviceID = int(dev.ID.Instance)
		device.Ip = host.IP.String()
		//device.Port = host.Port
		device.NetworkNumber = int(dev.Addr.Net)
		device.MaxApdu = dev.MaxApdu
		device.Segmentation = uint32(dev.Segmentation)
		fmt.Println("--------device------------", dev.ID.Instance)
		pprint.PrintJOSN(device)
		fmt.Println("--------device------------", dev.ID.Instance)

		details, err := device.GetDeviceDetails(dev.ID.Instance)
		if err != nil {
			fmt.Println("discover err", err)
		}
		fmt.Println("--------device---details---------", dev.ID.Instance)
		pprint.PrintJOSN(details)
		fmt.Println("--------device---details---------", dev.ID.Instance)

	}
	return err
}
