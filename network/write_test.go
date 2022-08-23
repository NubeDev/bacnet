package network

import (
	"fmt"
	"github.com/NubeDev/bacnet/btypes"
	"testing"
)

func TestDevice_Write(t *testing.T) {
	localDevice, err := New(&Network{Interface: iface, Port: 47809})
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

	err = device.Write(&Write{
		ObjectID:      0,
		ObjectType:    btypes.MultiStateOutput,
		Prop:          85,
		WriteValue:    1,
		WriteNull:     false,
		WritePriority: 16,
	})
	fmt.Println(err)
	if err != nil {
		return
	}
}
