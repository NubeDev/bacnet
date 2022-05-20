package btypes

import (
	"fmt"
	"testing"
)

func TestSupported(t *testing.T) {

	ss := ServicesSupported{}
	for supported := range ss.ListAll() {
		fmt.Println(supported)
	}

}
