package bacnet

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"testing"
)

func TestIP(t *testing.T) {

	buf := new(bytes.Buffer)
	var num uint16 = 47808
	err := binary.Write(buf, binary.BigEndian, num)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}

}
