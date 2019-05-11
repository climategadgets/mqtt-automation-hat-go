// +build arm

package sn3218

import (
	"fmt"
	"testing"
	"time"
)

func TestLightRing(t *testing.T) {

	values := [18]byte{}
	for offset := 0; offset < 18; offset++ {
		// VT: NOTE: 0xFF is way too bright
		values[offset] = 0x55
	}

	driver := GetSN3218()

	driver.Reset()
	driver.Enable(true)
	defer driver.Enable(false)

	var shift uint32

	driver.output(values)

	for shift = 0; shift < 18; shift++ {

		var mask uint32
		mask = 0x01 << shift
		fmt.Printf("EnableLEDs(%018b)\n", mask)
		driver.EnableLEDs(mask)
		time.Sleep(250 * time.Millisecond)
	}
}
