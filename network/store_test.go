package network

import (
	"fmt"
	"github.com/NubeDev/bacnet"
	pprint "github.com/NubeDev/bacnet/helpers/print"
	"testing"
)

func TestNetwork(t *testing.T) {

	client, err := New(&Network{Interface: iface, Port: 47808})
	if err != nil {
		fmt.Println("ERR-client", err)
		return
	}
	defer client.NetworkClose()
	go client.NetworkRun()

	wi := &bacnet.WhoIsOpts{
		High:            0,
		Low:             0,
		GlobalBroadcast: true,
		NetworkNumber:   0,
	}

	//client.Whois(wi)

	cli, ok := BacStore.Get("1")

	fmt.Println(cli, ok)

	aa := cli.(*Network)
	aa.Port = 47808
	aa.Interface = "wlp3s0"

	//if err != nil {
	//	return
	//}
	defer aa.NetworkClose()
	go aa.NetworkRun()

	whois, err := aa.Whois(wi)
	if err != nil {

	}

	pprint.PrintJOSN(err)
	pprint.PrintJOSN(whois)
	//
	//aa.Whois(wi)
	////close
	//aa.ClientClose()

}

func TestUpdateDev(t *testing.T) {

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

	objects, err := device.GetDevicePoints(202)
	if err != nil {
		return
	}
	pprint.PrintJOSN(objects)
	//err = device.UpdateDevice("123")
	fmt.Println(err)
	if err != nil {
		fmt.Println(err)
	}
	device.Ip = "192.168.15.15"
	//err = device.UpdateDevice("123")
	fmt.Println(err)
	if err != nil {
		fmt.Println(err)
	}
	objects, err = device.GetDevicePoints(202)
	if err != nil {
		fmt.Println(err)
	}
	pprint.PrintJOSN(objects)

	device.Ip = "192.168.15.191"
	//err = device.UpdateDevice("123")
	fmt.Println(err)
	if err != nil {
		fmt.Println(err)
	}
	objects, err = device.GetDevicePoints(202)
	if err != nil {
		//return
	}
	pprint.PrintJOSN(objects)

}
