package automation_hat

// 3.3V ADC, with no LED
type adc struct {
	channel    int
	maxVoltage float64
}

// 24V ADC, with a LED
type adcLed struct {
	adc
	ledContainer
}

func GetADC33(channel int, maxVoltage float64) ADC33 {
	return adc{channel: channel, maxVoltage: maxVoltage}
}

func GetADC24(channel int, maxVoltage float64, ledPin int) ADC24 {
	return adcLed{adc: adc{channel: channel, maxVoltage: maxVoltage}, ledContainer: ledContainer{GetLED(ledPin)}}
}
