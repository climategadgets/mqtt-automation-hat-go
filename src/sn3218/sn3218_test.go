// +build arm

package sn3218

import (
	"testing"
	"time"
)

// TestLightRing lights one LED at a time using EnableLEDs()
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
		driver.EnableLEDs(mask)
		time.Sleep(100 * time.Millisecond)
	}
}

func TestLightFadeSimple(t *testing.T) {

	driver := GetSN3218()
	defer driver.Close()

	driver.Reset()
	driver.Enable(true)

	// 0b111111111111111111, all of them
	driver.EnableLEDs(0x3FFFF)

	fade(driver)

}

func TestLightGammaAllocation(t *testing.T) {

	driver := GetSN3218()
	defer driver.Close()

	for channel := 0; channel < 18; channel++ {
		if driver.GetChannelGamma(uint8(channel)) != nil {
			t.Fatalf("channel %d gamma is not nil", channel)
		}
	}
}

func TestLightFadeInverted(t *testing.T) {

	driver := GetSN3218()
	defer driver.Close()

	driver.Reset()
	driver.Enable(true)

	// 0b111000000000111, just the first and last 3 LEDs of the linear group
	driver.EnableLEDs(0x7007)

	// Create inverted intensity map
	inversion := [256]byte{}

	for offset := 0; offset < 256; offset++ {
		inversion[offset] = byte(0xFF - offset)
	}
	// Invert intensity for first 3 LEDs
	for channel := 0; channel < 3; channel++ {
		driver.SetChannelGamma(byte(channel), &inversion)
	}

	fade(driver)
}

func fade(driver SN3218) {

	values := [18]byte{}

	for intensity := 0; intensity < 0xFF; intensity++ {

		for offset := 0; offset < 18; offset++ {
			values[offset] = byte(intensity)
		}
		driver.Output(values)
	}

	for intensity := 0xFF; intensity > 0; intensity-- {

		for offset := 0; offset < 18; offset++ {
			values[offset] = byte(intensity)
		}
		driver.Output(values)
	}
}

// TestLightRingSet lights LEDs one by one using SetLED()
func TestLightFill(t *testing.T) {

	values := [18]byte{}
	for offset := 0; offset < 18; offset++ {
		// VT: NOTE: 0xFF is way too bright
		values[offset] = 0x55
	}

	driver := GetSN3218()
	defer driver.Close()

	driver.Reset()
	driver.Enable(true)

	// 0b111111111111111111, all of them
	driver.EnableLEDs(0x3FFFF)

	for channel := 0; channel < 18; channel++ {
		driver.SetLED(byte(channel), 0x55)
		time.Sleep(100 * time.Millisecond)

		if driver.GetLED(byte(channel)) != 0x55 {
			t.Fatalf("value mismatch for channel %d, expected 0x55, received %x", channel, driver.GetLED(byte(channel)))
		}
	}
	time.Sleep(500 * time.Millisecond)
}
