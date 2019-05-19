package automation_hat

// 3.3V ADC, with no LED
type adc struct {
	channel    uint8
	maxVoltage float64
}

// 24V ADC, with a LED
type adcLed struct {
	adc
	ledContainer
}

func GetADC33(channel uint8, maxVoltage float64) ADC33 {
	return adc{channel: channel, maxVoltage: maxVoltage}
}

func GetADC24(control chan<- interface{}, channel uint8, maxVoltage float64, ledPin uint8) ADC24 {
	return adcLed{adc: adc{channel: channel, maxVoltage: maxVoltage}, ledContainer: ledContainer{GetLED(control, ledPin)}}
}

func (adc adc) Get() float64 {
	panic("Not Implemented")
}
