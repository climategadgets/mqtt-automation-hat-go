package automation_hat

type input struct {
	pin int
	led Light
}

func GetInput(pin int, ledPin int) Input {

	return input{ pin: pin, led: GetLED(ledPin)}
}

func (i input) Light() Light {
	return i.led
}