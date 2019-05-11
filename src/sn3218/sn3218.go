package sn3218

import (
	"bitbucket.org/gmcbay/i2c"
	"sync"
)

// Roughly based on:
// https://github.com/schoentoon/piglow
// https://github.com/pimoroni/sn3218/blob/master/library/sn3218.py
// http://www.si-en.com/uploadpdf/s2011517171720.pdf

const (
	i2cAddress      = 0x54
	cmdEnableOutput = 0x00
	cmdSetPwmValues = 0x01
	cmdEnableLeds   = 0x13
	cmdUpdate       = 0x16
	cmdReset        = 0x17
)

type sn3218 struct {
	bus    *i2c.I2CBus
	values [18]byte
	gamma  [18]*[256]byte
}

type driverLocker struct {
	mu     sync.Mutex
	driver *sn3218
}

var theDriver driverLocker

// Obtain access to SN3218 singleton instance.
// This method performs lazy initialization, the instance doesn't exist before first invocation.
func GetSN3218() SN3218 {

	theDriver.mu.Lock()
	defer theDriver.mu.Unlock()

	if theDriver.driver == nil {

		sn3218, err := newSN3218()

		if err != nil {
			// VT: FIXME: This will end badly if this function is called repeatedly, the error is not propagated
			theDriver.driver = nil
		}

		theDriver.driver = sn3218
	}

	return theDriver.driver
}

func (driver *sn3218) Close() error {

	if err := driver.Enable(false); err != nil {
		return err
	}
	return driver.Reset()
}

// Reset resets all hardware registers and gamma correction
func (driver *sn3218) Reset() error {

	for channel := 0; channel < 18; channel++ {
		driver.SetChannelGamma(byte(channel), nil)
	}
	return driver.bus.WriteByteBlock(i2cAddress, cmdReset, []byte{0xFF})
}

func (driver sn3218) Enable(enable bool) error {
	if enable {
		return driver.bus.WriteByteBlock(i2cAddress, cmdEnableOutput, []byte{0x01})
	} else {
		return driver.bus.WriteByteBlock(i2cAddress, cmdEnableOutput, []byte{0x00})
	}
}

// EnableLEDs enables or disables an individual LED channel.
// The argument is a binary channel mask.
func (driver sn3218) EnableLEDs(mask uint32) error {

	if err := driver.bus.WriteByteBlock(i2cAddress, cmdEnableLeds, []byte{byte(mask & 0x3F), byte((mask >> 6) & 0x3F), byte((mask >> 12) & 0x3F)}); err != nil {
		return err
	}

	return driver.bus.WriteByteBlock(i2cAddress, cmdUpdate, []byte{0xFF})
}

// GetChannelGamma returns current Gamma Correction value for the channel (see the PDF at the top).
// nil return value indicates there is no gamma correction in place.
func (driver sn3218) GetChannelGamma(channel uint8) *[256]byte {
	return driver.gamma[channel]
}

// SetChannelGamma provides Gamma Correction for the channel (see the PDF at the top).
// nil value for the gamma argument means gamma correction will not be performed and intensity value will be used raw.
func (driver *sn3218) SetChannelGamma(channel uint8, gamma *[256]byte) {
	driver.gamma[channel] = gamma
}

func (driver sn3218) Output(values [18]byte) error {

	mapped := [18]byte{}

	for channel := 0; channel < 18; channel++ {
		if driver.gamma[channel] == nil {
			mapped[channel] = values[channel]
		} else {
			mapped[channel] = driver.gamma[channel][values[channel]]
		}
	}

	if err := driver.bus.WriteByteBlock(i2cAddress, cmdSetPwmValues, mapped[0:18]); err != nil {
		return err
	}

	return driver.bus.WriteByteBlock(i2cAddress, cmdUpdate, []byte{0xFF})
}

func newSN3218() (*sn3218, error) {

	// VT: FIXME: Bus ID may not be 1 on some revisions
	// See https://github.com/pimoroni/sn3218/blob/master/library/sn3218.py
	bus, err := i2c.Bus(1)

	if err != nil {
		return nil, err
	}

	values := [18]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	return &sn3218{bus, values, [18]*[256]byte{}}, nil
}
