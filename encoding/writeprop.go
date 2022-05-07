package encoding

import "github.com/NubeDev/bacnet/btypes"

// WriteProperty encodes a write property request
func (e *Encoder) WriteProperty(invokeID uint8, data btypes.PropertyData) error {
	a := btypes.APDU{
		DataType: btypes.ConfirmedServiceRequest,
		Service:  btypes.ServiceConfirmedWriteProperty,
		MaxSegs:  0,
		MaxApdu:  MaxAPDU,
		InvokeId: invokeID,
	}
	e.APDU(a)

	tagID, err := e.readPropertyHeader(0, &data)
	if err != nil {
		return err
	}

	prop := data.Object.Properties[0]

	// Tag 3 - the value (unlike other values, this is just a raw byte array)
	e.openingTag(tagID)
	e.AppData(prop.Data)
	e.closingTag(tagID)

	tagID++

	// Tag 4 - Optional priorty tag
	// Priority set
	if prop.Priority != btypes.Normal {
		e.contextUnsigned(tagID, uint32(prop.Priority))
	}
	return e.Error()
}
