package local

import (
	"fmt"
	"github.com/NubeDev/bacnet"
	"github.com/NubeDev/bacnet/btypes"

	//"github.com/NubeDev/bacnet"

	"testing"
)

var iface = "wlp3s0"
var localDevicePort = 47808
var deviceIP = "192.168.15.202"
var deviceID = 202
var networkNumber = 0
var macMSTP = 0
var segmentation = btypes.SegmentedBoth
var MaxApdu uint32 = btypes.MaxAPDU1476

/*
MaxApdu
0 = 50
1 = 128
2 = 206 jci PCG
3 = 480 honeywell spyder
4 = 1024
5 = 1476  easyIO-30p when over IP

BACnetSegmentation:
segmented-both:0
segmented-transmit:1
segmented-receive:2
no-segmentation: 3
*/

func TestWhoIs(t *testing.T) {

	client, err := New(&Local{Interface: iface, Port: localDevicePort})
	if err != nil {
		fmt.Println("ERR-client", err)
		return
	}
	defer client.ClientClose()
	go client.ClientRun()

	whoIs, err := client.bacnet.WhoIs(&bacnet.WhoIsOpts{NetworkNumber: 4})
	if err != nil {
		fmt.Println("ERR-whoIs", err)
		return
	}

	for _, dev := range whoIs {
		fmt.Println(dev.ID)
		fmt.Println(dev.Vendor)
	}

}

func TestReadObjects(t *testing.T) {

	localDevice, err := New(&Local{Interface: iface, Port: localDevicePort})
	if err != nil {
		fmt.Println("ERR-client", err)
		return
	}
	defer localDevice.ClientClose()
	go localDevice.ClientRun()

	device, err := NewDevice(localDevice, &Device{Ip: deviceIP, DeviceID: deviceID, NetworkNumber: networkNumber, MacMSTP: macMSTP, MaxApdu: uint32(MaxApdu), Segmentation: uint32(segmentation)})
	if err != nil {
		return
	}

	objects, err := device.DeviceObjects()
	fmt.Println(objects, err)
	if err != nil {
		return
	}
	for i, a := range objects {
		fmt.Println(i, a.Type)
	}

}

func TestRead(t *testing.T) {

	localDevice, err := New(&Local{Interface: iface, Port: localDevicePort})
	if err != nil {
		fmt.Println("ERR-client", err)
		return
	}
	defer localDevice.ClientClose()
	go localDevice.ClientRun()

	device, err := NewDevice(localDevice, &Device{Ip: deviceIP, DeviceID: deviceID, NetworkNumber: networkNumber, MacMSTP: macMSTP, MaxApdu: MaxApdu, Segmentation: uint32(segmentation)})
	if err != nil {
		return
	}

	obj := &Object{
		ObjectID:   202,
		ObjectType: btypes.DeviceType,
		Prop:       btypes.ProtocolServicesSupported,
		ArrayIndex: btypes.ArrayAll, //btypes.ArrayAll

	}

	out, err := device.Read(obj)
	fmt.Println(err)
	fmt.Println(out)
	//fmt.Println("DATA", out.Object.Properties[0].Data)

}
