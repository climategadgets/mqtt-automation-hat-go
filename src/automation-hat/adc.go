package automation_hat

type adc struct {
	channel    int
	maxVoltage float64
	led_container
}

func GetADC24(channel int, maxVoltage float64, ledPin int) ADC {

	result := adc{channel: channel, maxVoltage: maxVoltage}

	result.led = GetLED(ledPin)

	return result
}
