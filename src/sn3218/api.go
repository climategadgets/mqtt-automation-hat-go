package sn3218

import "io"

type SN3218 interface {
	io.Closer

	Reset() error
	Enable(state bool) error
	EnableLEDs(mask uint32) error
	GetChannelGamma(channel uint8) *[256]byte
	SetChannelGamma(channel uint8, gamma *[256]byte)
	Output(values [18]byte) error
	SetLED(channel uint8, intensity byte) error
	GetLED(channel uint8) byte
}
