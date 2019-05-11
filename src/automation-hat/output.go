package automation_hat

type output struct {
	messageBus
	pin uint8
	ledContainer
}

func GetOutput(control chan<- interface{}, pin uint8, ledPin uint8) Output {
	return output{messageBus: messageBus{control}, pin: pin, ledContainer: ledContainer{GetLED(control, ledPin)}}
}

func (o output) Light() Light {
	return o.led
}
