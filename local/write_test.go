package local

import (
	"fmt"
	"github.com/NubeDev/bacnet/btypes"
	"testing"
)

func TestWrite(t *testing.T) {

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

	//write an AO
	var writeValueAo float32 = -11
	device.Write(&Write{ObjectID: 1, ObjectType: btypes.AnalogOutput, Prop: btypes.PropPresentValue, WriteValue: writeValueAo})

	//write an BO
	var writeValueBO uint32 = 1
	device.PointWriteBool(&Point{ObjectID: 1, ObjectType: btypes.BinaryOutput}, writeValueBO)

	//device.PointWriteAnalogue(&Point{ObjectID: 1, ObjectType: btypes.BinaryOutput}, 1)

}
