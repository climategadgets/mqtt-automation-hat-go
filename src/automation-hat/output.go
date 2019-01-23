package automation_hat

type output struct {
	pin int
	ledContainer
}

func GetOutput(pin int, ledPin int) Output {
	return output{pin: pin, ledContainer: ledContainer{GetLED(ledPin)}}
}

func (o output) Light() Light {
	return o.led
}
