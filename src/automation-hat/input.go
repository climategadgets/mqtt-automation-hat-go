package automation_hat

type input struct {
	pin uint8
	ledContainer
}

func GetInput(pin uint8, ledPin uint8) Input {
	return input{pin: pin, ledContainer: ledContainer{GetLED(ledPin)}}
}

func (i input) Light() Light {
	return i.led
}
