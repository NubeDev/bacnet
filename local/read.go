package local

import (
	"fmt"
	"github.com/NubeDev/bacnet"
	"github.com/NubeDev/bacnet/btypes"
	log "github.com/sirupsen/logrus"
)

func (device *Device) Read(objectID int, objectType btypes.ObjectType, prop btypes.PropertyType) (btypes.PropertyData, error) {
	//get object list
	rp := btypes.PropertyData{
		Object: btypes.Object{
			ID: btypes.ObjectID{
				Type:     objectType,
				Instance: btypes.ObjectInstance(objectID),
			},
			Properties: []btypes.Property{
				{
					Type:       prop,
					ArrayIndex: bacnet.ArrayAll,
				},
			},
		},
	}
	out, err := device.bacnet.ReadProperty(device.Dev, rp)
	if err != nil {
		if rp.Object.Properties[0].Type == btypes.PropObjectList {
			log.Error("Note: PropObjectList reads may need to be broken up into multiple reads due to length. Read index 0 for array length")
		}
		return out, nil
	}
	if len(out.Object.Properties) == 0 {
		fmt.Println("No value returned")
		return out, nil
	}
	return out, nil
}
