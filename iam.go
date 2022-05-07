package bacnet

import (
	"github.com/NubeDev/bacnet/btypes"
	"github.com/NubeDev/bacnet/encoding"
)

func (c *client) iAm(dest btypes.Address) error {
	npdu := &btypes.NPDU{
		Version:               btypes.ProtocolVersion,
		Destination:           &dest,
		IsNetworkLayerMessage: false,
		ExpectingReply:        false,
		Priority:              btypes.Normal,
		HopCount:              btypes.DefaultHopCount,
	}
	enc := encoding.NewEncoder()
	enc.NPDU(npdu)

	//	iams := []btypes.ObjectID{btypes.ObjectID{Instance: 1, Type: 5}}
	//	enc.IAm(iams)
	_, err := c.Send(dest, npdu, enc.Bytes())
	return err
}
