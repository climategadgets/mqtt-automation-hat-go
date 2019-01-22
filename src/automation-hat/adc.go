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

	result := adc{channel: channel, maxVoltage: maxVoltage}

	return result
}

func GetADC24(channel int, maxVoltage float64, ledPin int) ADC24 {

	result := adcLed{}

	result.channel = channel
	result.maxVoltage = maxVoltage
	result.led = GetLED(ledPin)

	return result
}
