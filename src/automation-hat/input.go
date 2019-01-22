package automation_hat

type input struct {
	pin int
	led_container
}

func GetInput(pin int, ledPin int) Input {

	result := input{pin: pin}

	result.led = GetLED(ledPin)

	return result
}

func (i input) Light() Light {
	return i.led
}
