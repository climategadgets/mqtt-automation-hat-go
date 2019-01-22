package automation_hat

import (
	"testing"
	"time"
)

// Flip all the lights except relay for 2 seconds, then shut off
func TestLights(t *testing.T) {

	hat := GetAutomationHAT()
	lights := make([]Light, 0)

	for _, adc := range hat.ADC24() {

		lights = append(lights, adc.Light())
	}

	for _, input := range hat.Input() {

		lights = append(lights, input.Light())
	}

	for _, output := range hat.Output() {

		lights = append(lights, output.Light())
	}

	lights = append(lights, hat.StatusLights().Power())
	lights = append(lights, hat.StatusLights().Comms())
	lights = append(lights, hat.StatusLights().Warn())

	for _, light := range lights {
		light.Set(true)
	}

	time.Sleep(2 * time.Second)

	for _, light := range lights {
		light.Set(false)
	}
}
