package rawv2_test

import (
	"fmt"

	"github.com/susji/ruuvi/data/rawv2"
)

func ExampleParse() {
	rawadv := []byte{
		0x05, 0x12, 0xfc, 0x53, 0x94, 0xc3, 0x7c, 0x00, 0x04, 0xff, 0xfc, 0x04,
		0x0c, 0xac, 0x36, 0x42, 0x00, 0xcd, 0xcb, 0xb8, 0x33, 0x4c, 0x88, 0x4f}
	p, _ := rawv2.Parse(rawadv)
	if p.Temperature.Valid {
		fmt.Printf("Temperature is valid and it is %.1f degrees Celsius\n", p.Temperature.Value)
	}
	fmt.Println("Your sensor's MAC address is", p.MAC)
	// Output:
	// Temperature is valid and it is 24.3 degrees Celsius
	// Your sensor's MAC address is cb:b8:33:4c:88:4f
}
