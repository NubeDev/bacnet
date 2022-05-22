package local

import (
	"fmt"
	pprint "github.com/NubeDev/bacnet/helpers/print"

	//"github.com/NubeDev/bacnet"

	"testing"
)

func TestDiscover(t *testing.T) {

	localDevice, err := New(&Local{Interface: iface, Port: localDevicePort})
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

	objects, err := device.DeviceObjects(202, true)
	if err != nil {
		return
	}
	pprint.PrintJOSN(objects)

}

func TestGetPointsList(t *testing.T) {

	localDevice, err := New(&Local{Interface: iface, Port: localDevicePort})
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

	objects, err := device.GetDevicePoints(202)
	if err != nil {
		return
	}
	pprint.PrintJOSN(objects)

}
