package automation_hat

type automation_hat struct {
	adc24  [3]ADC
	input  [3]Input
	output [3]Output
	relay  [3]Relay
	adc12  ADC
	status StatusLights
}

type status_lights struct {
	power Light
	comms Light
	warn  Light
}

// Obtain access to the AutomationHAT singleton instance.
// This method performs lazy initialization, the instance doesn't exist before first invocation.
func GetAutomationHAT() AutomationHAT {

	hat := automation_hat{}

	hat.adc24[0] = GetADC24(0, 25.85, 0)
	hat.adc24[1] = GetADC24(1, 25.85, 1)
	hat.adc24[2] = GetADC24(2, 25.85, 2)

	hat.input[0] = GetInput(26, 14)
	hat.input[1] = GetInput(20, 13)
	hat.input[2] = GetInput(21, 12)

	hat.output[0] = GetOutput(5, 3)
	hat.output[1] = GetOutput(12, 4)
	hat.output[2] = GetOutput(6, 5)

	hat.relay[0] = GetRelay(13, 6, 7)
	hat.relay[1] = GetRelay(19, 8, 9)
	hat.relay[2] = GetRelay(16, 10, 11)

	return hat
}

func (hat automation_hat) Relay() [3]Relay {
	return hat.relay
}

func (hat automation_hat) ADC24() [3]ADC {
	return hat.adc24
}

func (hat automation_hat) Input() [3]Input {
	return hat.input
}

func (hat automation_hat) Output() [3]Output {
	return hat.output
}

func (hat automation_hat) ADC12() ADC {
	return hat.adc12
}

func (hat automation_hat) StatusLights() StatusLights {
	return hat.status
}

func (sl status_lights) Power() Light {
	return sl.power
}

func (sl status_lights) Comms() Light {
	return sl.comms
}

func (sl status_lights) Warn() Light {
	return sl.warn
}
