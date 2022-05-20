package local

import (
	"fmt"
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

	discover, err := device.DeviceDiscover(202)
	fmt.Println(discover, err)

}
