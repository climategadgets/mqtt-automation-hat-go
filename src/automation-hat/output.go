package automation_hat

type output struct {
	pin uint8
	ledContainer
}

func GetOutput(pin uint8, ledPin uint8) Output {
	return output{pin: pin, ledContainer: ledContainer{GetLED(ledPin)}}
}

func (o output) Light() Light {
	return o.led
}
