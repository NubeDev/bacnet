package network

import (
	"fmt"
	"github.com/NubeDev/bacnet/btypes"
	"testing"
)

func TestPointDetails(t *testing.T) {

	localDevice, err := New(&Network{Interface: iface, Port: 47808})
	if err != nil {
		fmt.Println("ERR-client", err)
		return
	}
	defer localDevice.NetworkClose()
	go localDevice.NetworkRun()

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

	localDevice, err := New(&Network{Interface: iface, Port: 47808})
	if err != nil {
		fmt.Println("ERR-client", err)
		return
	}
	defer localDevice.NetworkClose()
	go localDevice.NetworkRun()

	device, err := NewDevice(localDevice, &Device{Ip: deviceIP, DeviceID: deviceID})
	if err != nil {
		return
	}

	pnt := &Point{
		ObjectID:   2,
		ObjectType: btypes.AnalogOutput,
	}

	readFloat64, err := device.PointReadFloat32(pnt)
	if err != nil {
		return
	}

	fmt.Println(readFloat64, err)

}

func TestReadWrite(t *testing.T) {

	localDevice, err := New(&Network{Interface: iface, Port: 47808})
	if err != nil {
		fmt.Println("ERR-client", err)
		return
	}
	defer localDevice.NetworkClose()
	go localDevice.NetworkRun()

	device, err := NewDevice(localDevice, &Device{Ip: deviceIP, DeviceID: deviceID})
	if err != nil {
		return
	}

	pnt := &Point{
		ObjectID:   1,
		ObjectType: btypes.AnalogValue,
	}

	err = device.PointWriteAnalogue(pnt, 1)
	if err != nil {
		//return
	}

	fmt.Println(err)

}
