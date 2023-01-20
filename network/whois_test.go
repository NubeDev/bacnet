package network

import (
	"fmt"
	"github.com/NubeDev/bacnet"
	pprint "github.com/NubeDev/bacnet/helpers/print"
	"github.com/kr/pretty"
	"testing"
)

func TestNetwork_Whois(t *testing.T) {
	localDevice, err := New(&Network{Interface: iface, Port: 47808})
	if err != nil {
		fmt.Println("ERR-client", err)
		return
	}
	defer localDevice.NetworkClose()
	go localDevice.NetworkRun()

	whois, err := localDevice.Whois(&bacnet.WhoIsOpts{
		Low:             0,
		High:            0,
		GlobalBroadcast: true,
		NetworkNumber:   0,
	})
	fmt.Println(err)
	if err != nil {
		return
	}

	pretty.Print(whois)
}

func TestNetwork_DeviceDiscover(t *testing.T) {
	localDevice, err := New(&Network{Interface: iface, Port: 47808})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer localDevice.NetworkClose()
	go localDevice.NetworkRun()

	devices, err := localDevice.NetworkDiscover(&bacnet.WhoIsOpts{
		Low:             0,
		High:            0,
		GlobalBroadcast: true,
		NetworkNumber:   0,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	pprint.PrintJOSN(devices)
	fmt.Println(devices)
}
