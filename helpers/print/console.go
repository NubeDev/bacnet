package pprint

import (
	"fmt"
)

func Print(i interface{}) {
	fmt.Printf("%+v\n", i)
	return
}

func Log(i interface{}) string {

	return fmt.Sprintf("%+v\n", i)
}
