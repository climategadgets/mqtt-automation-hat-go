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
	defer driver.Close()

	driver.Reset()
	driver.Enable(true)

	var shift uint32

	driver.Output(values)

	for shift = 0; shift < 18; shift++ {

		var mask uint32
		mask = 0x01 << shift
		fmt.Printf("EnableLEDs(%018b)\n", mask)
		driver.EnableLEDs(mask)
		time.Sleep(250 * time.Millisecond)
	}
}
