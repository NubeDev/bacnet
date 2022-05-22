package units

import (
	"fmt"
	"strconv"
	"testing"
)

func (i Unit) String2() string {

	fmt.Println()
	aa := _Unit_index_0
	for i2, aaa := range aa {
		fmt.Println("aa", i2, aaa)
	}
	fmt.Println("aa", aa)

	//aa := _Unit_name_0[_Unit_index_0[i]:_Unit_index_0[i+1]]
	//
	//for ii, ww := range _Unit_name_0[_Unit_index_0] {
	//	fmt.Println(ii, ww)
	//}

	switch {
	case i <= 104:
		//fmt.Println(_Unit_name_0[_Unit_index_0[i]:_Unit_index_0[i+1]])
		return _Unit_name_0[_Unit_index_0[i]:_Unit_index_0[i+1]]
	case 115 <= i && i <= 222:
		i -= 115
		return _Unit_name_1[_Unit_index_1[i]:_Unit_index_1[i+1]]
	case 224 <= i && i <= 236:
		i -= 224
		return _Unit_name_2[_Unit_index_2[i]:_Unit_index_2[i+1]]
	default:
		return "Unit(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}

func TestMarshal(t *testing.T) {
	u := Unit.String2(62)
	fmt.Println(u)
	//get()

}
