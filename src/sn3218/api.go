package sn3218

import "io"

type SN3218 interface {
	io.Closer

	Reset() error
	Enable(state bool) error
	EnableLEDs(mask uint32) error
	SetChannelGamma(channel uint8, gamma [256]uint8) error
	Output(values [18]byte) error
}
