package bacnet

import (
	"fmt"
	"github.com/NubeDev/bacnet/btypes"
	"github.com/NubeDev/bacnet/datalink"
	ip2bytes "github.com/NubeDev/bacnet/helpers/ipbytes"
	log "github.com/sirupsen/logrus"
	"testing"
)

var iface = "wlp3s0"
var deviceIP = "192.168.15.202"
var deviceID = 202
var networkNumber = 0
var deviceHardwareMac = 0
var objectID = 1

func TestRead(t *testing.T) {

	dataLink, err := datalink.NewUDPDataLink(iface, 47808)
	if err != nil {
		log.Fatal(err)
	}
	c := NewClient(dataLink, 0)
	defer c.Close()
	go c.Run()

	ip, err := ip2bytes.New(deviceIP, uint16(47808))
	if err != nil {
		return
	}

	addr := btypes.Address{
		Net: uint16(networkNumber),
		Mac: ip,
		Adr: []uint8{uint8(deviceHardwareMac)},
	}
	object := btypes.ObjectID{
		Type:     btypes.DeviceType,
		Instance: btypes.ObjectInstance(deviceID),
	}

	dest := btypes.Device{
		ID:   object,
		Addr: addr,
	}

	rp := btypes.PropertyData{
		Object: btypes.Object{
			ID: btypes.ObjectID{
				Type:     btypes.AnalogOutput,
				Instance: btypes.ObjectInstance(objectID),
			},
			Properties: []btypes.Property{
				{
					Type:       btypes.PropPresentValue,
					ArrayIndex: ArrayAll,
				},
			},
		},
	}
	out, err := c.ReadProperty(dest, rp)
	if err != nil {
		if rp.Object.Properties[0].Type == btypes.PropObjectList {
			log.Error("Note: PropObjectList reads may need to be broken up into multiple reads due to length. Read index 0 for array length")
		}
		fmt.Println("ReadProperty err")
		log.Fatal(err)
	}
	if len(out.Object.Properties) == 0 {
		fmt.Println("No value returned")
		return
	}
	fmt.Println(out.Object.Properties[0].Data)

}
