package automation_hat

type output struct {
	pin int
	led Light
}

func GetOutput(pin int, ledPin int) Output {

	return output{ pin: pin, led: GetLED(ledPin)}
}

func (o output) Light() Light {
	return o.led
}
