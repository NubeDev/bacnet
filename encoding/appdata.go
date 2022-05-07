package encoding

import (
	"fmt"

	"github.com/NubeDev/bacnet/btypes"
)

const (
	tagNull            uint8 = 0
	tagBool            uint8 = 1
	tagUint            uint8 = 2
	tagInt             uint8 = 3
	tagReal            uint8 = 4
	tagDouble          uint8 = 5
	tagOctetString     uint8 = 6
	tagCharacterString uint8 = 7
	tagBitString       uint8 = 8
	tagEnumerated      uint8 = 9
	tagDate            uint8 = 10
	tagTime            uint8 = 11
	tagObjectID        uint8 = 12
	tagReserve1        uint8 = 13
	tagReserve2        uint8 = 14
	tagReserve3        uint8 = 15
	maxTag             uint8 = 16
)

// Other values omitted here can have variable length
const (
	realLen     uint32 = 4
	doubleLen   uint32 = 8
	dateLen     uint32 = 4
	timeLen     uint32 = 4
	objectIDLen uint32 = 4
)

// epochYear is an increment to all non-stored values. This year is chosen in
// the standard. Why? No idea. God help us all if bacnet hits the 255 + 1990
// limit
const epochYear = 1990

// If the values == 0XFF, that means it is not specified. We will take that to
const notDefined = 0xff

// All app layer is non-context specific
const appLayerContext = false

func IsOddMonth(month int) bool {
	return month == 13
}

func IsEvenMonth(month int) bool {
	return month == 14
}

func IsLastDayOfMonth(day int) bool {
	return day == 32
}

func IsEvenDayOfMonth(day int) bool {
	return day == 33
}

func IsOddDayOfMonth(day int) bool {
	return day == 32
}
func (e *Encoder) string(s string) {
	e.write(stringUTF8)
	e.write([]byte(s))
}
func (d *Decoder) string(s *string, len int) error {
	var t stringType
	d.decode(&t)
	if t != stringUTF8 {
		return fmt.Errorf("unsupported string format %d", t)
	}

	b := make([]byte, len)
	d.decode(b)
	*s = string(b)
	return d.Error()
}
func (e *Encoder) octetstring(b []byte) {
	e.write([]byte(b))
}
func (d *Decoder) octetstring(b *[]byte, len int) {
	*b = make([]byte, len)
	d.decode(b)
}

func (e *Encoder) date(dt btypes.Date) {
	// We don't want to override an unspecified time date
	if dt.Year != btypes.UnspecifiedTime {
		e.write(uint8(dt.Year - epochYear))
	} else {
		e.write(uint8(dt.Year))
	}
	e.write(uint8(dt.Month))
	e.write(uint8(dt.Day))
	e.write(uint8(dt.DayOfWeek))
}

func (d *Decoder) date(dt *btypes.Date, length int) {
	if length <= 0 {
		return
	}
	data := make([]byte, length)
	_, d.err = d.Read(data)
	if d.err != nil {
		return
	}
	if len(data) < 4 {
		return
	}

	if dt.Year != btypes.UnspecifiedTime {
		dt.Year = int(data[0]) + epochYear
	} else {
		dt.Year = int(data[0])
	}

	dt.Month = int(data[1])
	dt.Day = int(data[2])
	dt.DayOfWeek = btypes.DayOfWeek(data[3])
}

func (e *Encoder) time(t btypes.Time) {
	e.write(uint8(t.Hour))
	e.write(uint8(t.Minute))
	e.write(uint8(t.Second))

	// Stored as 1/100 of a second
	e.write(uint8(t.Millisecond / 10))
}
func (d *Decoder) time(t *btypes.Time, length int) {
	if length <= 0 {
		return
	}
	data := make([]byte, length)
	if _, d.err = d.Read(data); d.err != nil {
		return
	}

	t.Hour = int(data[0])
	t.Minute = int(data[1])
	t.Second = int(data[2])
	t.Millisecond = int(data[3]) * 10

}

func (e *Encoder) boolean(x bool) {
	// Boolean information is stored into the length field
	var length uint32
	if x {
		length = 1
	} else {
		length = 0
	}
	e.tag(tagInfo{ID: tagBool, Context: appLayerContext, Value: length})
}

func (e *Encoder) real(x float32) {
	e.write(x)
}

func (d *Decoder) real(x *float32) {
	d.decode(x)
}

func (e *Encoder) double(x float64) {
	e.write(x)
}

func (d *Decoder) double(x *float64) {
	d.decode(x)
}

func (e *Encoder) AppData(i interface{}) error {
	switch val := i.(type) {
	case float32:
		e.tag(tagInfo{ID: tagReal, Context: appLayerContext, Value: realLen})
		e.real(val)
	case float64:
		e.tag(tagInfo{ID: tagDouble, Context: appLayerContext, Value: realLen})
		e.double(val)
	case bool:
		e.boolean(val)
	case string:
		// Add 1 to length to account for the encoding byte
		e.tag(tagInfo{ID: tagCharacterString, Context: appLayerContext, Value: uint32(len(val) + 1)})
		e.string(val)
	case uint32:
		//AIDAN changed TAG from  tagUint to tagEnumerated to get BO, BVs working
		length := valueLength(val)
		e.tag(tagInfo{ID: tagEnumerated, Context: appLayerContext, Value: uint32(length)})
		e.unsigned(val)
	case int32:
		v := uint32(val)
		length := valueLength(v)
		e.tag(tagInfo{ID: tagInt, Context: appLayerContext, Value: uint32(length)})
		e.unsigned(v)
	// Enumerated is pretty much a wrapper for a uint32 with an enumerated associated with it.
	case btypes.Enumerated:
		v := uint32(val)
		length := valueLength(v)
		e.tag(tagInfo{ID: tagEnumerated, Context: appLayerContext, Value: uint32(length)})
		e.enumerated(v)
	case btypes.ObjectID:
		e.tag(tagInfo{ID: tagObjectID, Context: appLayerContext, Value: objectIDLen})
		e.objectId(val.Type, val.Instance)

	case btypes.Null:
		e.tag(tagInfo{ID: tagNull, Context: appLayerContext})

	default:
		err := fmt.Errorf("Unknown type %T", i)
		// Set global error
		e.err = err
		return err
	}
	return nil
}

func (d *Decoder) AppDataOfTag(tag uint8, len int) (interface{}, error) {
	switch tag {
	case tagNull:
		return btypes.Null{}, nil
	case tagBool:
		// Originally this was in C so non 0 values are considered
		// true
		return len > 0, d.Error()
	case tagUint:
		return d.unsigned(len), d.Error()
	case tagInt:
		return d.signed(len), d.Error()
	case tagReal:
		var x float32
		d.real(&x)
		return x, d.Error()
	case tagDouble:
		var x float64
		d.double(&x)
		return x, d.Error()
	case tagOctetString:
		var b []byte
		d.octetstring(&b, len)
		return b, d.Error()

	case tagCharacterString:
		var s string
		// Subtract 1 to length to account for the encoding byte
		err := d.string(&s, len-1)
		return s, err
	case tagBitString:
		return d.bitString(len), d.Error()
	case tagEnumerated:
		return d.enumerated(len), d.Error()
	case tagDate:
		var date btypes.Date
		d.date(&date, len)
		return date, d.Error()
	case tagTime:
		var t btypes.Time
		d.time(&t, len)
		return t, d.Error()
	case tagObjectID:
		objType, objInstance := d.objectId()
		return btypes.ObjectID{
			Type:     btypes.ObjectType(objType),
			Instance: objInstance,
		}, d.Error()
	default:
		return nil, fmt.Errorf("Unsupported tag: %d", tag)
	}
}
func (d *Decoder) AppData() (interface{}, error) {
	tag, _, lenvalue := d.tagNumberAndValue()
	return d.AppDataOfTag(tag, int(lenvalue))
}
