package btypes

import "fmt"

type ObjectType uint16

const (
	TypeAnalogInput       = 0
	TypeAnalogOutput      = 1
	TypeAnalogValue       = 2
	TypeBinaryInput       = 3
	TypeBinaryOutput      = 4
	TypeBinaryValue       = 5
	TypeDeviceType        = 8
	TypeFile              = 10
	TypeMultiStateInput   = 13
	TypeMultiStateOutput  = 14
	TypeNotificationClass = 15
	TypeMultiStateValue   = 19
	TypeTrendLog          = 20
	TypeCharacterString   = 40
)

const (
	AnalogInput       ObjectType = 0
	AnalogOutput      ObjectType = 1
	AnalogValue       ObjectType = 2
	BinaryInput       ObjectType = 3
	BinaryOutput      ObjectType = 4
	BinaryValue       ObjectType = 5
	DeviceType        ObjectType = 8
	File              ObjectType = 10
	MultiStateInput   ObjectType = 13
	MultiStateOutput  ObjectType = 14
	NotificationClass ObjectType = 15
	MultiStateValue   ObjectType = 19
	TrendLog          ObjectType = 20
	CharacterString   ObjectType = 40
)

const (
	AnalogInputStr       = "Analog Input"
	AnalogOutputStr      = "Analog Output"
	AnalogValueStr       = "Analog Value"
	BinaryInputStr       = "Binary Input"
	BinaryOutputStr      = "Binary Output"
	BinaryValueStr       = "Binary Value"
	Devicebtypestr       = "Device"
	FileStr              = "File"
	NotificationClassStr = "Notification Class"
	MultiStateValueStr   = "Multi-State Value"
	MultiStateInputStr   = "Multi-State Input"
	MultiStateOutputStr  = "Multi-State Output"
	TrendLogStr          = "Trend Log"
	CharacterStringStr   = "Character String"
)

var objTypeMap = map[ObjectType]string{
	AnalogInput:       AnalogInputStr,
	AnalogOutput:      AnalogOutputStr,
	AnalogValue:       AnalogValueStr,
	BinaryInput:       BinaryInputStr,
	BinaryOutput:      BinaryOutputStr,
	BinaryValue:       BinaryValueStr,
	DeviceType:        Devicebtypestr,
	File:              FileStr,
	NotificationClass: NotificationClassStr,
	MultiStateValue:   MultiStateValueStr,
	MultiStateInput:   MultiStateInputStr,
	MultiStateOutput:  MultiStateOutputStr,
	TrendLog:          TrendLogStr,
	CharacterString:   CharacterStringStr,
}

var objStrTypeMap = map[string]ObjectType{
	AnalogInputStr:       AnalogInput,
	AnalogOutputStr:      AnalogOutput,
	AnalogValueStr:       AnalogValue,
	BinaryInputStr:       BinaryInput,
	BinaryOutputStr:      BinaryOutput,
	BinaryValueStr:       BinaryValue,
	Devicebtypestr:       DeviceType,
	FileStr:              File,
	NotificationClassStr: NotificationClass,
	MultiStateValueStr:   MultiStateValue,
	MultiStateInputStr:   MultiStateInput,
	MultiStateOutputStr:  MultiStateOutput,
	TrendLogStr:          TrendLog,
	CharacterStringStr:   CharacterString,
}

func GetType(s string) ObjectType {
	t, ok := objStrTypeMap[s]
	if !ok {
		return 0
	}
	return t
}

func (t ObjectType) String() string {
	s, ok := objTypeMap[t]
	if !ok {
		return fmt.Sprintf("Unknown (%d)", t)
	}
	return fmt.Sprintf("%s", s)
}

type ObjectInstance uint32

type ObjectID struct {
	Type     ObjectType
	Instance ObjectInstance
}

// String returns a pretty print of the ObjectID structure
func (id ObjectID) String() string {
	return fmt.Sprintf("Instance: %d Type: %s", id.Instance, id.Type.String())
}

type Object struct {
	Name        string
	Description string
	ID          ObjectID
	Properties  []Property `json:",omitempty"`
}
