package automation_hat

import (
	"go.uber.org/zap"
	"sync"
)

// 3.3V ADC, with no LED
type adc struct {
	messageBus
	channel    uint8
	maxVoltage float64
}

// 24V ADC, with a LED
type adcLed struct {
	adc
	ledContainer
}

type adcCommand struct {
	adc
	signal *float64
	done   sync.WaitGroup
}

func GetADC33(control chan<- interface{}, channel uint8, maxVoltage float64) ADC33 {
	return adc{messageBus: messageBus{control}, channel: channel, maxVoltage: maxVoltage}
}

func GetADC24(control chan<- interface{}, channel uint8, maxVoltage float64, ledPin uint8) ADC24 {
	return adcLed{
		adc:          adc{messageBus: messageBus{control}, channel: channel, maxVoltage: maxVoltage},
		ledContainer: ledContainer{GetLED(control, ledPin)}}
}

func (adc adc) Get() float64 {

	// VT: NOTE: This read operation is done in such an obscure way to ensure that all I/O operations
	// are performed in one place, sequentially, with no locks and no race conditions

	if adc.control == nil {
		panic("nil control channel")
	}

	var done sync.WaitGroup

	done.Add(1)

	command := adcCommand{adc: adc, done: done}

	adc.control <- command

	done.Wait()

	zap.S().Warnw("analog input", "signal", *command.signal)

	return *command.signal
}
