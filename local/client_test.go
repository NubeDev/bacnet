package local

import (
	"fmt"
	"github.com/NubeDev/bacnet"
	"github.com/NubeDev/bacnet/btypes"

	//"github.com/NubeDev/bacnet"

	"testing"
)

var iface = "wlp3s0"
var localDevicePort = 47809
var deviceIP = "192.168.15.194"
var deviceID = 1234
var networkNumber = 0
var macMSTP = 0
var segmentation = 3
var MaxApdu = 480

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

	device, err := NewDevice(localDevice, &Device{Ip: deviceIP, DeviceID: deviceID, NetworkNumber: networkNumber, MacMSTP: macMSTP, MaxApdu: 206, Segmentation: uint32(segmentation)})
	if err != nil {
		return
	}

	out, err := device.Read(1234, btypes.DeviceType, btypes.PropObjectName)
	fmt.Println(err)
	fmt.Println(out)
	fmt.Println("DATA", out.Object.Properties[0].Data)

}
