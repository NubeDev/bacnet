package bacnet

import (
	"encoding/json"
	"fmt"
	"github.com/NubeDev/bacnet/datalink"
	"github.com/NubeDev/bacnet/encoding"
	"log"
	"testing"

	"github.com/NubeDev/bacnet/btypes"
)

const interfaceName = "eth0"
const testServer = 1234

// TestMain are general test
func TestUdpDataLink(t *testing.T) {
	dataLink, err := datalink.NewUDPDataLink(interfaceName, 0)
	if err != nil {
		t.Fatal(err)
	}
	c := NewClient(dataLink, 0)
	c.Close()

	_, err = datalink.NewUDPDataLink("pizzainterfacenotreal", 0)
	if err == nil {
		t.Fatal("Successfully passed a false interface.")
	}
}

func TestMac(t *testing.T) {
	var mac []byte
	json.Unmarshal([]byte("\"ChQAzLrA\""), &mac)
	l := len(mac)
	p := uint16(mac[l-1])<<8 | uint16(mac[l-1])
	log.Printf("%d", p)
}

func TestServices(t *testing.T) {
	dataLink, err := datalink.NewUDPDataLink(interfaceName, 0)
	if err != nil {
		t.Fatal(err)
	}
	c := NewClient(dataLink, 0)
	defer c.Close()

	t.Run("Read Property", func(t *testing.T) {
		testReadPropertyService(c, t)
	})

	t.Run("Who Is", func(t *testing.T) {
		testWhoIs(c, t)
	})

	t.Run("WriteProperty", func(t *testing.T) {
		testWritePropertyService(c, t)
	})

}

func testReadPropertyService(c Client, t *testing.T) {
	dev, err := c.WhoIs(testServer, testServer)
	read := btypes.PropertyData{
		Object: btypes.Object{
			ID: btypes.ObjectID{
				Type:     btypes.AnalogValue,
				Instance: 1,
			},
			Properties: []btypes.Property{
				btypes.Property{
					Type:       btypes.PropObjectName, // Present value
					ArrayIndex: ArrayAll,
				},
			},
		},
	}
	if len(dev) == 0 {
		t.Fatalf("Unable to find device id %d", testServer)
	}

	resp, err := c.ReadProperty(dev[0], read)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Response: %v", resp.Object.Properties[0].Data)
}

func testWhoIs(c Client, t *testing.T) {
	dev, err := c.WhoIs(testServer-1, testServer+1)
	if err != nil {
		t.Fatal(err)
	}
	if len(dev) == 0 {
		t.Fatalf("Unable to find device id %d", testServer)
	}
}

// This test will first cconver the name of an analogue sensor to a different
// value, read the property to make sure the name was changed, revert back, and
// ensure that the revert was successful
func testWritePropertyService(c Client, t *testing.T) {
	const targetName = "Hotdog"
	dev, err := c.WhoIs(testServer, testServer)
	wp := btypes.PropertyData{
		Object: btypes.Object{
			ID: btypes.ObjectID{
				Type:     btypes.AnalogValue,
				Instance: 1,
			},
			Properties: []btypes.Property{
				btypes.Property{
					Type:       btypes.PropObjectName, // Present value
					ArrayIndex: ArrayAll,
					Priority:   btypes.Normal,
				},
			},
		},
	}

	if len(dev) == 0 {
		t.Fatalf("Unable to find device id %d", testServer)
	}
	resp, err := c.ReadProperty(dev[0], wp)
	if err != nil {
		t.Fatal(err)
	}
	// Store the original response since we plan to put it back in after
	org := resp.Object.Properties[0].Data
	t.Logf("original name is: %d", org)

	wp.Object.Properties[0].Data = targetName
	err = c.WriteProperty(dev[0], wp)
	if err != nil {
		t.Fatal(err)
	}

	resp, err = c.ReadProperty(dev[0], wp)
	if err != nil {
		t.Fatal(err)
	}

	d := resp.Object.Properties[0].Data
	s, ok := d.(string)
	if !ok {
		log.Fatalf("unexpected return type %T", d)
	}

	if s != targetName {
		log.Fatalf("write to name %s did not successed, name was %s", targetName, s)
	}

	// Revert Changes
	wp.Object.Properties[0].Data = org
	err = c.WriteProperty(dev[0], wp)
	if err != nil {
		t.Fatal(err)
	}

	resp, err = c.ReadProperty(dev[0], wp)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Object.Properties[0].Data != org {
		t.Fatalf("unable to revert name back to original value %v: name is %v", org, resp.Object.Properties[0].Data)
	}
}

func TestDeviceClient(t *testing.T) {
	dataLink, err := datalink.NewUDPDataLink("本地连接", 47809)
	if err != nil {
		fmt.Println(err)
		return
	}
	c := NewClient(dataLink, 0)
	go c.Run()

	devs, err := c.WhoIs(-1, -1)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%+v\n", devs)
	//	c.Objects(devs[0])

	prop, err := c.ReadProperty(
		devs[0],
		btypes.PropertyData{
			Object: btypes.Object{
				ID: btypes.ObjectID{
					Type:     btypes.AnalogInput,
					Instance: 0,
				},
				Properties: []btypes.Property{{
					Type:       85,
					ArrayIndex: encoding.ArrayAll,
				}},
			},
			ErrorClass: 0,
			ErrorCode:  0,
		})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(prop.Object.Properties)

	props, err := c.ReadMultiProperty(devs[0], btypes.MultiplePropertyData{Objects: []btypes.Object{
		{
			ID: btypes.ObjectID{
				Type:     btypes.AnalogInput,
				Instance: 0,
			},
			Properties: []btypes.Property{
				{
					Type:       8,
					ArrayIndex: encoding.ArrayAll,
				},
				/*	{
					Type:       85,
					ArrayIndex: encoding.ArrayAll,
				},*/
			},
		},
	}})

	fmt.Println(props)
	if err != nil {
		fmt.Println(err)
		return
	}
}