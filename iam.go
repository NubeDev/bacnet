package bacnet

import (
	"github.com/NubeDev/bacnet/btypes"
	"github.com/NubeDev/bacnet/encoding"
)

func (c *client) IAm(dest btypes.Address, iam btypes.IAm) error {
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
	enc.IAm(iam)
	_, err := c.Send(dest, npdu, enc.Bytes())
	return err
}
