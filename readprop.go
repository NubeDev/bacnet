package bacnet

import (
	"context"
	"fmt"
	"github.com/NubeDev/bacnet/btypes"
	"github.com/NubeDev/bacnet/encoding"
	"log"
	"time"
)

// ReadProperty reads a single property from a single object in the given device.
func (c *client) ReadProperty(dest btypes.Device, rp btypes.PropertyData) (btypes.PropertyData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	id, err := c.tsm.ID(ctx)
	if err != nil {
		return btypes.PropertyData{}, fmt.Errorf("unable to get an transaction id: %v", err)
	}
	defer c.tsm.Put(id)
	fmt.Println("------")
	fmt.Println("NET", dest.Addr.Net)
	fmt.Println("MAC", dest.Addr.Mac)
	//dest.Addr.Mac = []uint8{0004}

	fmt.Println("MAC", dest.Addr.Mac)
	fmt.Println("MacLen", dest.Addr.MacLen)
	fmt.Println("Source", c.dataLink.GetMyAddress())
	fmt.Println("------")
	//dest.Addr.Adr = []byte{4}
	//dest.Addr.Adr = []uint8{4}
	//dest.Addr.Net = 4
	//dest.Addr.MacLen = 1
	fmt.Println(dest.Addr.Net)
	enc := encoding.NewEncoder()
	npdu := &btypes.NPDU{
		Version:               btypes.ProtocolVersion,
		Destination:           &dest.Addr,
		Source:                c.dataLink.GetMyAddress(),
		IsNetworkLayerMessage: false,
		ExpectingReply:        true,
		Priority:              btypes.Normal,
		HopCount:              btypes.DefaultHopCount,
	}
	enc.NPDU(npdu)

	err = enc.ReadProperty(uint8(id), rp)
	if enc.Error() != nil || err != nil {
		return btypes.PropertyData{}, err
	}

	// the value filled doesn't matter. it just needs to be non nil
	err = fmt.Errorf("go")
	for count := 0; err != nil && count < 2; count++ {
		var b []byte
		var out btypes.PropertyData
		_, err = c.Send(dest.Addr, npdu, enc.Bytes())
		if err != nil {
			log.Print(err)
			continue
		}

		var raw interface{}
		raw, err = c.tsm.Receive(id, time.Duration(5)*time.Second)
		if err != nil {
			continue
		}
		switch v := raw.(type) {
		case error:
			return out, v
		case []byte:
			b = v
		default:
			return out, fmt.Errorf("received unknown datatype %T", raw)
		}

		dec := encoding.NewDecoder(b)

		var apdu btypes.APDU
		if err = dec.APDU(&apdu); err != nil {
			continue
		}
		if apdu.Error.Class != 0 || apdu.Error.Code != 0 {
			err = fmt.Errorf("received error, class: %d, code: %d", apdu.Error.Class, apdu.Error.Code)
			continue
		}

		if err = dec.ReadProperty(&out); err != nil {
			continue
		}

		return out, dec.Error()
	}
	return btypes.PropertyData{}, err
}
