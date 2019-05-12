package automation_hat

import (
	"fmt"
	"go.uber.org/zap"
	"testing"
	"time"
)

// Flip all the lights except relay for 2 seconds, then shut off
func TestLights(t *testing.T) {

	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
	defer logger.Sync()

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

	for _, relay := range hat.Relay() {
		lights = append(lights, relay.Light()[0])
		lights = append(lights, relay.Light()[1])
	}

	lights = append(lights, hat.StatusLights().Power())
	lights = append(lights, hat.StatusLights().Comms())
	lights = append(lights, hat.StatusLights().Warn())

	zap.S().Infof("lights collected: %d", len(lights))

	for offset, light := range lights {
		// VT: NOTE: Index will be different from the pin
		zap.S().Infow("light", "index", offset, "light", fmt.Sprintf("%v", light))
	}

	for _, light := range lights {
		light.Set(true)
	}

	time.Sleep(2 * time.Second)

	for _, light := range lights {
		light.Set(false)
	}
}
