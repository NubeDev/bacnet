package local

import (
	"fmt"
	"github.com/NubeDev/bacnet"
	"github.com/NubeDev/bacnet/btypes"
	log "github.com/sirupsen/logrus"
)

type Write struct {
	ObjectID      btypes.ObjectInstance
	ObjectType    btypes.ObjectType
	Prop          btypes.PropertyType
	WriteValue    interface{}
	WriteNull     bool
	WritePriority uint8
}

func (device *Device) Write(write *Write) error {
	var err error
	writeValue := write.WriteValue

	rp := btypes.PropertyData{
		Object: btypes.Object{
			ID: btypes.ObjectID{
				Type:     write.ObjectType,
				Instance: btypes.ObjectInstance(write.ObjectID),
			},
			Properties: []btypes.Property{
				{
					Type:       write.Prop,
					ArrayIndex: bacnet.ArrayAll,
				},
			},
		},
	}

	if write.WriteNull {
		writeValue = btypes.Null{}
	} else {
		switch writeValue.(type) {
		case uint32:
			out, _ := writeValue.(uint32)
			writeValue = out
		case float32:
			out, _ := writeValue.(float32)
			writeValue = out
		case float64:
			out, _ := writeValue.(float64)
			writeValue = out
		case string:
			writeValue = fmt.Sprintf("%s", writeValue)
		default:
			err = fmt.Errorf("unable to handle a type %T", writeValue)
			return err
		}
		if err != nil {
			log.Printf("Expects a %T", rp.Object.Properties[0].Data)
			return err
		}
	}

	rp.Object.Properties[0].Data = writeValue
	log.Printf("Writting: %v", writeValue)
	err = device.bacnet.WriteProperty(device.Dev, rp)
	if err != nil {
		return err
	}
	return nil
}
