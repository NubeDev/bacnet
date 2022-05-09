package local

import (
	"fmt"
	"github.com/NubeDev/bacnet/btypes"
	log "github.com/sirupsen/logrus"
)

func (device *Device) DeviceObjects() ([]btypes.ObjectID, error) {
	//get object list
	rp := btypes.PropertyData{
		Object: btypes.Object{
			ID: btypes.ObjectID{
				Type:     btypes.DeviceType,
				Instance: btypes.ObjectInstance(device.DeviceID),
			},
			Properties: []btypes.Property{
				btypes.Property{
					Type:       btypes.PropObjectList,
					ArrayIndex: 97,
				},
			},
		},
	}
	out, err := device.bacnet.ReadProperty(device.Dev, rp)
	if err != nil {
		if rp.Object.Properties[0].Type == btypes.PropObjectList {
			log.Error("Note: PropObjectList reads may need to be broken up into multiple reads due to length. Read index 0 for array length")
		}
		return nil, nil
	}
	if len(out.Object.Properties) == 0 {
		fmt.Println("No value returned")
		return nil, nil
	}
	ids, ok := out.Object.Properties[0].Data.([]interface{})
	if !ok {

	}
	var objectIDS []btypes.ObjectID
	for _, id := range ids {
		objectID := id.(btypes.ObjectID) // do something with this
		objectIDS = append(objectIDS, objectID)
	}

	return objectIDS, nil
}
