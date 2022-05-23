package network

import (
	"fmt"
	"github.com/NubeDev/bacnet"
	pprint "github.com/NubeDev/bacnet/helpers/print"
	"testing"
)

func TestStore(t *testing.T) {

	client, err := New(&Local{Interface: iface, Port: localDevicePort})
	if err != nil {
		fmt.Println("ERR-client", err)
		return
	}
	defer client.ClientClose()
	go client.ClientRun()

	cli, ok := cache.Get("1")

	fmt.Println(cli, ok)

	aa := cli.(*Local)
	aa.Port = 47808
	aa.Interface = "enp0s31f6"
	aa.UpdateClient()
	//if err != nil {
	//	return
	//}
	defer aa.ClientClose()
	go aa.ClientRun()
	wi := &bacnet.WhoIsOpts{
		High:            0,
		Low:             0,
		GlobalBroadcast: true,
		NetworkNumber:   0,
	}
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
