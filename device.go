package bacnet

import (
	"fmt"
	"github.com/NubeDev/bacnet/btypes"
	"github.com/NubeDev/bacnet/btypes/ndpu"
	"github.com/NubeDev/bacnet/datalink"
	"github.com/NubeDev/bacnet/encoding"
	"github.com/NubeDev/bacnet/helpers/validation"
	"github.com/NubeDev/bacnet/tsm"
	"github.com/NubeDev/bacnet/utsm"
	log "github.com/sirupsen/logrus"
	"io"
	"sync"
	"time"
)

const mtuHeaderLength = 4
const defaultStateSize = 20
const forwardHeaderLength = 10

type Client interface {
	io.Closer
	ClientClose(closeLogs bool) error
	ClientRun()
	WhoIs(wh *WhoIsOpts) ([]btypes.Device, error)
	WhatIsNetworkNumber() []*btypes.Address
	IAm(dest btypes.Address, iam btypes.IAm) error
	WhoIsRouterToNetwork() (resp *[]btypes.Address)
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
	log            *log.Logger
}

type ClientBuilder struct {
	DataLink   datalink.DataLink
	Interface  string
	Ip         string
	Port       int
	SubnetCIDR int
	MaxPDU     uint16
}

// NewClient creates a new client with the given interface and
func NewClient(cb *ClientBuilder) (Client, error) {
	var err error
	var dataLink datalink.DataLink
	iface := cb.Interface
	ip := cb.Ip
	port := cb.Port
	maxPDU := cb.MaxPDU
	//check ip
	ok := validation.ValidIP(ip)
	if !ok {

	}
	//check port
	if port == 0 {
		port = datalink.DefaultPort
	}
	ok = validation.ValidPort(port)
	if !ok {

	}
	//check adpu
	if maxPDU == 0 {
		maxPDU = btypes.MaxAPDU
	}
	//build datalink
	if iface != "" {
		dataLink, err = datalink.NewUDPDataLink(iface, port)
		if err != nil {
			//log.Fatal(err)
		}
	} else {
		//check subnet
		sub := cb.SubnetCIDR
		ok = validation.ValidCIDR(ip, sub)
		if !ok {

		}
		dataLink, err = datalink.NewUDPDataLinkFromIP(ip, sub, port)
		if err != nil {
			//log.Fatal(err)
		}
	}

	l := log.New()
	l.Formatter = &log.TextFormatter{}
	l.SetLevel(log.DebugLevel)

	cli := &client{
		dataLink: dataLink,
		tsm:      tsm.New(defaultStateSize),
		utsm: utsm.NewManager(
			utsm.DefaultSubscriberTimeout(time.Second*time.Duration(10)),
			utsm.DefaultSubscriberLastReceivedTimeout(time.Second*time.Duration(2)),
		),
		readBufferPool: sync.Pool{New: func() interface{} {
			return make([]byte, maxPDU)
		}},
		log: l,
	}
	return cli, err
}

func (c *client) ClientRun() {
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
		networkList, err := dec.NPDU(&npdu)
		if err != nil {
			return
		}

		if npdu.IsNetworkLayerMessage {
			c.log.Debug("Ignored Network Layer Message")
			if npdu.NetworkLayerMessageType == ndpu.NetworkIs {
				c.utsm.Publish(int(npdu.Source.Net), npdu)
				//return
			}
			if npdu.NetworkLayerMessageType == ndpu.IamRouterToNetwork {
				c.utsm.Publish(int(npdu.Source.Net), networkList)
				//return
			}

		}

		// We want to keep the APDU intact, so we will get a snapshot before decoding
		send := dec.Bytes()
		err = dec.APDU(&apdu)
		if err != nil {
			c.log.Errorf("Issue decoding APDU: %v", err)
			return
		}
		switch apdu.DataType {
		case btypes.UnconfirmedServiceRequest:
			if apdu.UnconfirmedService == btypes.ServiceUnconfirmedIAm {
				dec = encoding.NewDecoder(apdu.RawData)
				var iam btypes.IAm
				err = dec.IAm(&iam)
				c.log.Debug("Received IAM Message", iam.ID)
				iam.Addr = *src

				if npdu.Source != nil {
					if npdu.Source.Net > 0 { // add in device network number
						log.Println("device-network-address", npdu.Source.Net)
						iam.Addr.Net = npdu.Source.Net
					}
					if len(npdu.Source.Adr) > 0 { // add in hardware mac
						log.Println("device-mstp-mac-address", npdu.Source.Adr)
						iam.Addr.Adr = npdu.Source.Adr
					}
				}
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
			err := fmt.Errorf("error class %s code %s", apdu.Error.Class.String(), apdu.Error.Code.String())
			err = c.tsm.Send(int(apdu.InvokeId), err)
			if err != nil {
				c.log.Debugf("unable to Send error to %d: %v", apdu.InvokeId, err)
			}
		default:
			// Ignore it
			log.WithFields(log.Fields{"raw": b}).Debug("An ignored packet went through")
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

type SetBroadcastType struct { //used to override the header.Function
	Set     bool
	BacFunc btypes.BacFunc
}

// Send transfers the raw apdu byte slice to the destination address.
func (c *client) Send(dest btypes.Address, npdu *btypes.NPDU, data []byte, broadcastType *SetBroadcastType) (int, error) {
	//broadcastType = &SetBroadcastType{}
	var header btypes.BVLC
	// Set packet type
	header.Type = btypes.BVLCTypeBacnetIP
	//if Adr is > 0 it must be an mst-tp device so send a UNICAST
	if len(dest.Adr) > 0 { //(aidan) not sure if this is correct, but it needs to be set to work to send (UNICAST) messages over a bacnet network
		// SET UNICAST FLAG
		// see http://www.bacnet.org/Tutorial/HMN-Overview/sld033.
		// see https://github.com/JoelBender/bacpypes/blob/9fca3f608a97a20807cd188689a2b9ff60b05085/doc/source/gettingstarted/gettingstarted001.rst#udp-communications-issues
		header.Function = btypes.BacFuncUnicast
	} else if dest.IsBroadcast() || dest.IsSubBroadcast() {
		// SET BROADCAST FLAG
		header.Function = btypes.BacFuncBroadcast
	} else {
		// SET UNICAST FLAG
		header.Function = btypes.BacFuncUnicast
	}

	if broadcastType != nil {
		if broadcastType.Set {
			header.Function = broadcastType.BacFunc
		}
	}

	header.Length = uint16(mtuHeaderLength + len(data))
	header.Data = data
	e := encoding.NewEncoder()
	err := e.BVLC(header)
	if err != nil {
		return 0, err
	}
	// use default udp type, src = network address (nil)
	return c.dataLink.Send(e.Bytes(), npdu, &dest)
}

func (c *client) ClientClose(closeLogs bool) error {
	if closeLogs {
		if f, ok := c.log.Out.(io.Closer); ok {
			return f.Close()
		}
	}
	return c.Close()
}

// Close free resources for the client. Always call this function when using NewClient
func (c *client) Close() error {
	if c.dataLink != nil {
		c.dataLink.Close()
	}

	return nil
}
