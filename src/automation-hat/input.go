package automation_hat

type input struct {
	pin int
	ledContainer
}

func GetInput(pin int, ledPin int) Input {
	return input{pin: pin, ledContainer: ledContainer{GetLED(ledPin)}}
}

func (i input) Light() Light {
	return i.led
}
