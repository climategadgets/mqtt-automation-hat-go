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

type SN3218 struct {
	bus    *i2c.I2CBus
	values [18]byte
}

type driverLocker struct {
	mu     sync.Mutex
	driver *SN3218
}

var theDriver driverLocker

// Obtain access to SN3218 singleton instance.
// This method performs lazy initialization, the instance doesn't exist before first invocation.
func GetSN3218() *SN3218 {

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

// Reset resets all hardware registers
func (driver SN3218) Reset() error {
	return driver.bus.WriteByteBlock(i2cAddress, cmdReset, []byte{0xFF})
}

func (driver SN3218) Enable(enable bool) error {
	if enable {
		return driver.bus.WriteByteBlock(i2cAddress, cmdEnableOutput, []byte{0x01})
	} else {
		return driver.bus.WriteByteBlock(i2cAddress, cmdEnableOutput, []byte{0x00})
	}
}

// EnableLEDs enables or disables an individual LED channel.
// The argument is a binary channel mask.
func (driver SN3218) EnableLEDs(mask uint32) error {

	if err := driver.bus.WriteByteBlock(i2cAddress, cmdEnableLeds, []byte{byte(mask & 0x3F), byte((mask >> 6) & 0x3F), byte((mask >> 12) & 0x3F)}); err != nil {
		return err
	}

	return driver.bus.WriteByteBlock(i2cAddress, cmdUpdate, []byte{0xFF})
}

// SetChannelGamma provides Gamma Correction (see the PDF at the top).
func (driver SN3218) SetChannelGamma(channel uint8, gamma [256]uint8) error {
	return nil
}

func (driver SN3218) output(values [18]byte) error {

	if err := driver.bus.WriteByteBlock(i2cAddress, cmdSetPwmValues, values[0:18]); err != nil {
		return err
	}

	return driver.bus.WriteByteBlock(i2cAddress, cmdUpdate, []byte{0xFF})
}

func newSN3218() (*SN3218, error) {

	// VT: FIXME: Bus ID may not be 1 on some revisions
	// See https://github.com/pimoroni/sn3218/blob/master/library/sn3218.py
	bus, err := i2c.Bus(1)

	if err != nil {
		return nil, err
	}

	values := [18]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	return &SN3218{bus, values}, nil
}
