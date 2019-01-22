package automation_hat

type output struct {
	pin int
	ledContainer
}

func GetOutput(pin int, ledPin int) Output {

	result := output{pin: pin}

	result.led = GetLED(ledPin)

	return result
}

func (o output) Light() Light {
	return o.led
}
