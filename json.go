package main

import (
	"encoding/json"
	"fmt"
)

type Data struct {
	IntF    int
	StringF string
	FloatF  float64
	SliceF  []string
	MapF    map[string]int
}

func (d Data) print() {
	fmt.Printf("IntF: %d\n", d.IntF)
	fmt.Printf("StringF: %s\n", d.StringF)
	fmt.Printf("FloatF: %f\n", d.FloatF)
	fmt.Printf("SliceF: %s\n", d.SliceF)
	fmt.Printf("MapF: %s\n", d.MapF)
}

func main() {

	data := Data{
		5,
		"55",
		55.5,
		[]string{"5", "55", "555"},
		map[string]int{"5": 5, "55": 55, "555": 555},
	}
	data.print()

	b, _ := json.Marshal(data)
	fmt.Printf("\nMarshalled data: %s\n", b)

}
