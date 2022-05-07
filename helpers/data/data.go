package data

import (
	"fmt"
	"github.com/NubeDev/bacnet/btypes"
)

func ToArr(d btypes.PropertyData) (ok bool, data []interface{}) {
	data, ok = d.Object.Properties[0].Data.([]interface{})
	if !ok {
		fmt.Println("unable to get object list")
		return ok, data
	}
	return

}

func ToUint32(d btypes.PropertyData) (ok bool, data uint32) {
	if len(d.Object.Properties) == 0 {
		fmt.Println("No value returned")
		return ok, data
	}
	data, ok = d.Object.Properties[0].Data.(uint32)
	return ok, data
}
