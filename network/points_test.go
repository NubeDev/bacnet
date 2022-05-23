package network

import (
	"fmt"
	"github.com/NubeDev/bacnet/btypes"
	"testing"
)

func TestPointDetails(t *testing.T) {

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

	pnt := &Point{
		ObjectID:   2,
		ObjectType: btypes.AnalogInput,
	}

	readFloat64, err := device.PointDetails(pnt)
	if err != nil {
		//return
	}

	fmt.Println(readFloat64, err)

}

func TestRead(t *testing.T) {

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

	pnt := &Point{
		ObjectID:   2,
		ObjectType: btypes.AnalogOutput,
	}

	readFloat64, err := device.PointReadFloat64(pnt)
	if err != nil {
		return
	}

	fmt.Println(readFloat64, err)

}
