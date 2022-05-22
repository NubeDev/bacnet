package services

import (
	"fmt"
	"testing"
)

func TestSupported(t *testing.T) {

	ss := Supported{}
	for supported := range ss.ListAll() {
		fmt.Println(supported)
	}

}
