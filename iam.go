package bacnet

import (
	"github.com/NubeDev/bacnet/btypes"
	"github.com/NubeDev/bacnet/encoding"
)

/*
not working
*/

func (c *client) IAm(dest btypes.Address, iam btypes.IAm) error {
	npdu := &btypes.NPDU{
		Version:     btypes.ProtocolVersion,
		Destination: &dest,
		//IsNetworkLayerMessage:   true,
		//NetworkLayerMessageType: 0x12,
		//Source:         c.dataLink.GetMyAddress(),
		ExpectingReply: false,
		Priority:       btypes.Normal,
		HopCount:       btypes.DefaultHopCount,
	}
	enc := encoding.NewEncoder()
	enc.NPDU(npdu)
	enc.IAm(iam)
	_, err := c.Send(dest, npdu, enc.Bytes())
	return err
}
