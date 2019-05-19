package automation_hat

type input struct {
	messageBus
	pin uint8
	ledContainer
}

func GetInput(control chan<- interface{}, pin uint8, ledPin uint8) Input {
	return input{messageBus: messageBus{control}, pin: pin, ledContainer: ledContainer{GetLED(control, ledPin)}}
}

func (input input) Get() bool {
	panic("Not Implemented")
}

func (i input) Light() Light {
	return i.led
}
