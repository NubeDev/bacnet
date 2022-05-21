package units

import (
	"fmt"
	"testing"
)

func TestMarshal(t *testing.T) {
	u := Unit.String(62)
	fmt.Println(u)
}
