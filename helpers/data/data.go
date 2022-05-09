package data

import (
	"fmt"
	"github.com/NubeDev/bacnet/btypes"
)

func ToArr(d btypes.PropertyData) (ok bool, out []interface{}) {
	out, ok = d.Object.Properties[0].Data.([]interface{})
	if !ok {
		fmt.Println("unable to get object list")
		return ok, out
	}
	return

}

func ToFloat64(d btypes.PropertyData) (ok bool, out float64) {
	if len(d.Object.Properties) == 0 {
		fmt.Println("No value returned")
		return ok, out
	}
	out, ok = d.Object.Properties[0].Data.(float64)
	return ok, out
}

func ToBool(d btypes.PropertyData) (ok bool, out bool) {
	if len(d.Object.Properties) == 0 {
		fmt.Println("No value returned")
		return ok, out
	}
	out, ok = d.Object.Properties[0].Data.(bool)
	return ok, out
}

func ToUint32(d btypes.PropertyData) (ok bool, out uint32) {
	if len(d.Object.Properties) == 0 {
		fmt.Println("No value returned")
		return ok, out
	}
	out, ok = d.Object.Properties[0].Data.(uint32)
	return ok, out
}
