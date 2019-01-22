package automation_hat

type adc struct {
	channel    int
	maxVoltage float64
	led        Light
}

func (adc adc) Light() Light {
	return adc.led
}

func GetADC24(channel int, maxVoltage float64, ledPin int) ADC {

	return adc{channel: channel, maxVoltage: maxVoltage, led: GetLED(ledPin)}
}
