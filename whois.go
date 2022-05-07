package bacnet

import "C"
import (
	"fmt"
	"github.com/NubeDev/bacnet/btypes"
	"github.com/NubeDev/bacnet/encoding"
)

type WhoIsOpts struct {
	High            int
	Low             int
	GlobalBroadcast bool
	NetworkNumber   uint16
}

// WhoIs finds all devices with ids between the provided low and high values.
// Use constant ArrayAll for both fields to scan the entire network at once.
// Using ArrayAll is highly discouraged for most networks since it can lead
// to a high congested network.
func (c *client) WhoIs(wh *WhoIsOpts) ([]btypes.Device, error) {
	dest := *c.dataLink.GetBroadcastAddress()
	enc := encoding.NewEncoder()
	low := wh.Low
	high := wh.High

	if wh.GlobalBroadcast {
		wh.NetworkNumber = btypes.GlobalBroadcast //65535

	}
	dest.Net = uint16(wh.NetworkNumber)
	npdu := &btypes.NPDU{
		Version:               btypes.ProtocolVersion,
		Destination:           &dest,
		Source:                c.dataLink.GetMyAddress(),
		IsNetworkLayerMessage: false,

		// We are not expecting a direct reply from a single destination
		ExpectingReply: false,
		Priority:       btypes.Normal,
		HopCount:       btypes.DefaultHopCount,
	}
	enc.NPDU(npdu)

	err := enc.WhoIs(int32(low), int32(high))
	if err != nil {
		return nil, err
	}

	// Subscribe to any changes in the the range. If it is a broadcast,
	var start, end int
	if low == -1 || high == -1 {
		start = 0
		end = maxInt
	} else {
		start = low
		end = high
	}

	// Run in parallel
	errChan := make(chan error)
	go func() {
		_, err = c.Send(dest, npdu, enc.Bytes())
		errChan <- err
	}()
	values, err := c.utsm.Subscribe(start, end)
	if err != nil {
		return nil, err
	}
	err = <-errChan
	if err != nil {
		return nil, err
	}
	fmt.Println("values", values)
	// Weed out values that are not important such as non object type
	// and that are not
	uniqueMap := make(map[btypes.ObjectInstance]btypes.Device)
	uniqueList := make([]btypes.Device, len(uniqueMap))

	//func cAddrToGoAddr(addr *C.BACNET_ADDRESS) *btypes.Address {
	//	var result btypes.Address
	//	for i := 0; i < len(result.Mac) && i < len(addr.mac); i++ {
	//	result.Mac[i] = uint8(addr.mac[i])
	//}
	//	result.MacLen = uint8(addr.mac_len)
	//	for i := 0; i < len(result.Adr) && i < len(addr.adr); i++ {
	//	result.Adr[i] = uint8(addr.adr[i])
	//}
	//	result.Len = uint8(addr.len)
	//	result.Net = uint16(addr.net)
	//	return &result
	//}
	//
	for _, v := range values {
		r, ok := v.(btypes.IAm)
		fmt.Println("WHOIS RES")
		fmt.Println("ADDR", r.Addr)
		fmt.Println("ALL", v)
		fmt.Println("ALL-ID", r.ID)
		fmt.Println("ALL-Vendor", r.Vendor)
		fmt.Println("ALL-MaxApdu", r.MaxApdu)
		fmt.Println("ALL-Segmentation", r.Segmentation)
		fmt.Println("WHOIS RES")

		// Skip non I AM responses
		if !ok {
			continue
		}

		// Check to see if we are in the map before inserting
		if _, ok := uniqueMap[r.ID.Instance]; !ok {
			dev := btypes.Device{
				Addr:         r.Addr,
				ID:           r.ID,
				MaxApdu:      r.MaxApdu,
				Segmentation: r.Segmentation,
				Vendor:       r.Vendor,
			}
			uniqueMap[r.ID.Instance] = btypes.Device(dev)
			uniqueList = append(uniqueList, dev)
		}
	}
	return uniqueList, err
}
