package bacnet

import (
	"fmt"
	"go/build"
	"os"
	"testing"
)

var iface = "enp0s31f6"

func TestIam(t *testing.T) {

	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}
	fmt.Println(gopath)

	cb := &ClientBuilder{
		Interface: iface,
	}
	c, _ := NewClient(cb)
	defer c.Close()
	go c.ClientRun()

	c.WhatIsNetworkNumber()

	//ip, err := ip2bytes.New(deviceIP, uint16(47808))
	//if err != nil {
	//	return
	//}
	//
	//addr := btypes.Address{
	//	Mac: ip,
	//}
	//
	//iam := btypes.IAm{
	//	ID:           btypes.ObjectID{Instance: 123, Type: 8},
	//	MaxApdu:      btypes.MaxAPDU480,
	//	Segmentation: btypes.Enumerated(segmentation.SegmentedBoth),
	//	Vendor:       234,
	//	Addr:         addr,
	//}
	//
	//pprint.Print(addr)
	//
	//err = c.IAm(addr, iam)
	//fmt.Println(err)
	//if err != nil {
	//	fmt.Println(err)
	//}

}
