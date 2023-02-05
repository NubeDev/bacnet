package network

import (
	"fmt"
	pprint "github.com/NubeDev/bacnet/helpers/print"

	//"github.com/NubeDev/bacnet"

	"testing"
)

func TestDiscover(t *testing.T) {

	localDevice, err := New(&Network{Interface: iface, Port: 47808})
	if err != nil {
		fmt.Println("ERR-client", err)
		return
	}
	defer localDevice.NetworkClose(false)
	go localDevice.NetworkRun()

	device, err := NewDevice(localDevice, &Device{Ip: deviceIP, DeviceID: deviceID})
	if err != nil {
		return
	}

	objects, err := device.DeviceObjects(12, true)
	if err != nil {
		return
	}
	pprint.PrintJOSN(objects)

}

func TestGetPointsList(t *testing.T) {

	localDevice, err := New(&Network{Interface: iface, Port: 47808})
	if err != nil {
		fmt.Println("ERR-client", err)
		return
	}
	defer localDevice.NetworkClose(false)
	go localDevice.NetworkRun()

	device, err := NewDevice(localDevice, &Device{Ip: deviceIP, DeviceID: deviceID})
	if err != nil {
		return
	}

	objects, err := device.GetDevicePoints(12)
	if err != nil {
		return
	}
	pprint.PrintJOSN(objects)

}
