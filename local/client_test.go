package local

import (
	"fmt"
	"github.com/NubeDev/bacnet"
	"github.com/NubeDev/bacnet/btypes"

	//"github.com/NubeDev/bacnet"

	"testing"
)

var iface = "wlp3s0"

func TestWhoIs(t *testing.T) {

	client, err := New(&Local{Interface: iface})
	if err != nil {
		fmt.Println("ERR-client", err)
		return
	}
	defer client.ClientClose()
	go client.ClientRun()

	whoIs, err := client.bacnet.WhoIs(&bacnet.WhoIsOpts{})
	if err != nil {
		fmt.Println("ERR-whoIs", err)
		return
	}

	for _, dev := range whoIs {
		fmt.Println(dev.ID)
	}

}

var deviceIP = "192.168.15.202"
var deviceID = 202

func TestReadObjects(t *testing.T) {

	localDevice, err := New(&Local{Interface: iface})
	if err != nil {
		fmt.Println("ERR-client", err)
		return
	}
	defer localDevice.ClientClose()
	go localDevice.ClientRun()

	device, err := NewDevice(localDevice, &Device{Ip: deviceIP, DeviceID: deviceID})
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

	localDevice, err := New(&Local{Interface: iface})
	if err != nil {
		fmt.Println("ERR-client", err)
		return
	}
	defer localDevice.ClientClose()
	go localDevice.ClientRun()

	device, err := NewDevice(localDevice, &Device{Ip: deviceIP, DeviceID: deviceID})
	if err != nil {
		return
	}

	out, err := device.Read(1, btypes.AnalogOutput, btypes.PropPresentValue)
	fmt.Println(err)
	fmt.Println(out)
	fmt.Println("DATA", out.Object.Properties[0].Data)

}
