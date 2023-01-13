package network

import (
	"github.com/NubeDev/bacnet"
	"github.com/NubeDev/bacnet/btypes"
	"github.com/NubeDev/bacnet/helpers/data"
	log "github.com/sirupsen/logrus"
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

func (net *Network) NetworkDiscover(options *bacnet.WhoIsOpts) ([]btypes.Device, error) {
	go net.NetworkRun()
	resp, err := net.bacnet.WhoIs(options)
	var devices []btypes.Device
	for _, device := range resp {
		buildDevice := &btypes.Device{
			DeviceID:      device.DeviceID,
			Ip:            device.Ip,
			Port:          device.Port,
			NetworkNumber: device.NetworkNumber,
			MacMSTP:       device.MacMSTP,
			MaxApdu:       device.MaxApdu,
			Segmentation:  device.Segmentation,
		}
		dev, err := btypes.NewDevice(buildDevice)
		details, err := net.GetDeviceDetails(*dev)
		if err != nil {
			return nil, err
		}
		device.DeviceName = details.Name
		device.VendorName = details.VendorName
		device.MaxApdu = details.MaxApdu
		device.Segmentation = btypes.Enumerated(details.Segmentation)
		devices = append(devices, device)
	}
	return devices, err
}

func (net *Network) GetDeviceDetails(device btypes.Device) (resp *DeviceDetails, err error) {
	resp = &DeviceDetails{}
	props := []btypes.PropertyType{btypes.PropObjectName, btypes.PropMaxAPDU, btypes.PropVendorName, btypes.PropSegmentationSupported}
	for _, prop := range props {
		property, err := net.bacnet.ReadProperty(device, buildObj(device.DeviceID, prop))
		if err != nil {
			log.Errorln("bacnet-master-GetDeviceDetails()", err.Error())
		}
		switch prop {
		case btypes.PropObjectName:
			_, read := data.ToStr(property)
			resp.Name = read
		case btypes.PropMaxAPDU:
			_, read := data.ToUint32(property)
			resp.MaxApdu = read
		case btypes.PropVendorName:
			_, read := data.ToStr(property)
			resp.VendorName = read
		case btypes.PropSegmentationSupported:
			_, read := data.ToUint32(property)
			resp.Segmentation = read
		case btypes.ProtocolServicesSupported:
			_, read := data.ToBitString(property)
			resp.ProtocolServicesSupported = read
		}
	}
	log.Infoln("bacnet-device name:", resp.Name)
	log.Infoln("bacnet-device vendor-name:", resp.VendorName)
	return resp, nil
}

func buildObj(id int, propertyType btypes.PropertyType) btypes.PropertyData {
	rp := btypes.PropertyData{
		Object: btypes.Object{
			ID: btypes.ObjectID{
				Type:     btypes.TypeDeviceType,
				Instance: btypes.ObjectInstance(id),
			},
			Properties: []btypes.Property{
				{
					Type:       propertyType,
					ArrayIndex: bacnet.ArrayAll, //bacnet.ArrayAll
				},
			},
		},
	}

	return rp
}
