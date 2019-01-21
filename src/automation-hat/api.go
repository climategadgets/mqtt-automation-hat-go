// Package automation_hat contains abstractions necessary to control the Pimoroni Automation HAT
// (https://shop.pimoroni.com/products/automation-hat)
// Baseline Python code is available at https://github.com/pimoroni/automation-hat
package automation_hat

type AutomationHAT interface {

	// 3 x 24V @ 2A relays (NC and NO terminals)
	Relay()  [3]Relay

	// 3 x 12-bit ADC @ 0-24V (±2% accuracy)
	ADC24()  [3]ADC

	// 3 x 24V tolerant buffered inputs
	Input()  [3]Input

	// 3 x 24V tolerant sinking outputs
	Output() [3] Output

	// 1 x 12-bit ADC @ 0-3.3V
	ADC12()  ADC

	StatusLights() StatusLights
}

type Switch interface {
	Get() bool
	Set(bool)
}

type Relay interface {

	Switch()
	Light() [2]Light
}

type ADC interface {

	Light() [3]Light
}

type Input interface {

	Light() [3]Light
}

type Output interface {

	Light() [3]Light
}


type StatusLights interface {

	Power() Light
	Comms() Light
	Warn() Light
}

type Light interface {

	Switch()
}
