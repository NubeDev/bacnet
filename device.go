package bacnet

import (
	"fmt"
	"github.com/NubeDev/bacnet/btypes"

	"github.com/NubeDev/bacnet/datalink"
	"github.com/NubeDev/bacnet/encoding"
	"io"
	"sync"
	"time"

	"github.com/NubeDev/bacnet/tsm"
	"github.com/NubeDev/bacnet/utsm"
	"github.com/sirupsen/logrus"
)

const mtuHeaderLength = 4
const defaultStateSize = 20
const forwardHeaderLength = 10

type Client interface {
	io.Closer
	Run()
	WhoIs(low, high int) ([]btypes.Device, error)
	Objects(dev btypes.Device) (btypes.Device, error)
	ReadProperty(dest btypes.Device, rp btypes.PropertyData) (btypes.PropertyData, error)
	ReadMultiProperty(dev btypes.Device, rp btypes.MultiplePropertyData) (btypes.MultiplePropertyData, error)
	WriteProperty(dest btypes.Device, wp btypes.PropertyData) error
	WriteMultiProperty(dev btypes.Device, wp btypes.MultiplePropertyData) error
}

type client struct {
	dataLink       datalink.DataLink
	tsm            *tsm.TSM
	utsm           *utsm.Manager
	readBufferPool sync.Pool
	log            *logrus.Logger
}

// NewClient creates a new client with the given interface and
// port.
func NewClient(dataLink datalink.DataLink, maxPDU uint16) Client {
	if maxPDU == 0 {
		maxPDU = btypes.MaxAPDU
	}
	log := logrus.New()
	log.Formatter = &logrus.TextFormatter{}
	log.SetLevel(logrus.InfoLevel)
	return &client{
		dataLink: dataLink,
		tsm:      tsm.New(defaultStateSize),
		utsm: utsm.NewManager(
			utsm.DefaultSubscriberTimeout(time.Second*time.Duration(10)),
			utsm.DefaultSubscriberLastReceivedTimeout(time.Second*time.Duration(2)),
		),
		readBufferPool: sync.Pool{New: func() interface{} {
			return make([]byte, maxPDU)
		}},
		log: log,
	}
}

func (c *client) Run() {
	var err error = nil
	for err == nil {
		b := c.readBufferPool.Get().([]byte)
		var addr *btypes.Address
		var n int
		addr, n, err = c.dataLink.Receive(b)
		if err != nil {
			continue
		}
		go c.handleMsg(addr, b[:n])
	}
}

func (c *client) handleMsg(src *btypes.Address, b []byte) {
	var header btypes.BVLC
	var npdu btypes.NPDU
	var apdu btypes.APDU

	dec := encoding.NewDecoder(b)
	err := dec.BVLC(&header)
	if err != nil {
		c.log.Error(err)
		return
	}

	if header.Function == btypes.BacFuncBroadcast || header.Function == btypes.BacFuncUnicast || header.Function == btypes.BacFuncForwardedNPDU {
		// Remove the header information
		b = b[mtuHeaderLength:]
		err = dec.NPDU(&npdu)
		if err != nil {
			return
		}

		if npdu.IsNetworkLayerMessage {
			c.log.Debug("Ignored Network Layer Message")
			return
		}

		// We want to keep the APDU intact so we will get a snapshot before decoding
		// further
		send := dec.Bytes()
		err = dec.APDU(&apdu)
		if err != nil {
			c.log.Errorf("Issue decoding APDU: %v", err)
			return
		}

		switch apdu.DataType {
		case btypes.UnconfirmedServiceRequest:
			if apdu.UnconfirmedService == btypes.ServiceUnconfirmedIAm {
				c.log.Debug("Received IAm Message")
				dec = encoding.NewDecoder(apdu.RawData)
				var iam btypes.IAm

				err = dec.IAm(&iam)

				iam.Addr = *src
				if err != nil {
					c.log.Error(err)
					return
				}
				c.utsm.Publish(int(iam.ID.Instance), iam)
			} else if apdu.UnconfirmedService == btypes.ServiceUnconfirmedWhoIs {
				dec := encoding.NewDecoder(apdu.RawData)
				var low, high int32
				dec.WhoIs(&low, &high)
				// For now we are going to ignore who is request.
				//log.WithFields(log.Fields{"low": low, "high": high}).Debug("WHO IS Request")
			} else {
				c.log.Errorf("Unconfirmed: %d %v", apdu.UnconfirmedService, apdu.RawData)
			}
		case btypes.SimpleAck:
			c.log.Debug("Received Simple Ack")
			err := c.tsm.Send(int(apdu.InvokeId), send)
			if err != nil {
				return
			}
		case btypes.ComplexAck:
			c.log.Debug("Received Complex Ack")
			err := c.tsm.Send(int(apdu.InvokeId), send)
			if err != nil {
				return
			}
		case btypes.ConfirmedServiceRequest:
			c.log.Debug("Received  Confirmed Service Request")
			err := c.tsm.Send(int(apdu.InvokeId), send)
			if err != nil {
				return
			}
		case btypes.Error:
			err := fmt.Errorf("error class %d code %d", apdu.Error.Class, apdu.Error.Code)
			err = c.tsm.Send(int(apdu.InvokeId), err)
			if err != nil {
				c.log.Debugf("unable to Send error to %d: %v", apdu.InvokeId, err)
			}
		default:
			// Ignore it
			//log.WithFields(log.Fields{"raw": b}).Debug("An ignored packet went through")
		}
	}

	if header.Function == btypes.BacFuncForwardedNPDU {
		// Right now we are ignoring the NPDU data that is stored in the packet. Eventually
		// we will need to check it for any additional information we can gleam.
		// NDPU has source
		b = b[forwardHeaderLength:]
		c.log.Debug("Ignored NDPU Forwarded")
	}
}

// Send transfers the raw apdu byte slice to the destination address.
func (c *client) Send(dest btypes.Address, npdu *btypes.NPDU, data []byte) (int, error) {
	var header btypes.BVLC

	// Set packet type
	header.Type = btypes.BVLCTypeBacnetIP

	if dest.IsBroadcast() || dest.IsSubBroadcast() {
		// SET BROADCAST FLAG
		header.Function = btypes.BacFuncBroadcast
	} else {
		// SET UNICAST FLAG
		header.Function = btypes.BacFuncUnicast
	}
	header.Function = btypes.BacFuncUnicast //TODO AIDAN had to set this to make it work over a network router
	header.Length = uint16(mtuHeaderLength + len(data))
	header.Data = data
	e := encoding.NewEncoder()
	err := e.BVLC(header)
	if err != nil {
		return 0, err
	}

	// use default udp type, src = local address (nil)
	return c.dataLink.Send(e.Bytes(), npdu, &dest)
}

// Close free resources for the client. Always call this function when using NewClient
func (c *client) Close() error {
	if c.dataLink != nil {
		c.dataLink.Close()
	}
	if f, ok := c.log.Out.(io.Closer); ok {
		return f.Close()
	}
	return nil
}