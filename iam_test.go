package bacnet

import (
	"fmt"
	"github.com/NubeDev/bacnet/btypes"
	ip2bytes "github.com/NubeDev/bacnet/helpers/ipbytes"
	pprint "github.com/NubeDev/bacnet/helpers/print"

	"go/build"
	"os"
	"testing"
)

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

	ip, err := ip2bytes.New(deviceIP, uint16(47808))
	if err != nil {
		return
	}

	addr := btypes.Address{
		Mac: ip,
	}

	iam := btypes.IAm{
		ID:           btypes.ObjectID{Instance: 33, Type: 8},
		MaxApdu:      btypes.MaxAPDU480,
		Segmentation: btypes.SegmentedBoth,
		Vendor:       22,
		Addr:         addr,
	}

	pprint.Print(addr)

	err = c.IAm(addr, iam)
	if err != nil {
		fmt.Println(err)
	}

}
