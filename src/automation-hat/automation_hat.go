package automation_hat

import (
	log "github.com/sirupsen/logrus"
	"github.com/stianeikeland/go-rpio"
	"sync"
)

type automationHatBase struct {
	adc24  [3]ADC24
	input  [3]Input
	output [3]Output
	relay  [3]Relay
	adc33  ADC33
	status StatusLights
}

type statusLights struct {
	power Light
	comms Light
	warn  Light
}

type automationHatFake struct {
	automationHatBase
}

type automationHatPi struct {
	automationHatBase
}

type hatLocker struct {
	mu  sync.Mutex
	hat AutomationHAT
}

var theHat hatLocker

// Obtain access to the AutomationHAT singleton instance.
// This method performs lazy initialization, the instance doesn't exist before first invocation.
func GetAutomationHAT() AutomationHAT {

	theHat.mu.Lock()
	defer theHat.mu.Unlock()

	if theHat.hat == nil {

		if GetRaspberryPiRevision() != nil {
			theHat.hat = newAutomationHAT()
		} else {
			theHat.hat = newAutomationFake()
		}
	}

	return theHat.hat
}

func newAutomationFake() AutomationHAT {
	log.Warn("using AutomationHAT fake")
	hat := automationHatBase{}

	initialize(&hat)

	return automationHatFake{hat}
}

func newAutomationHAT() AutomationHAT {

	log.Info("creating new instance of AutomationHAT")

	hat := automationHatBase{}

	initialize(&hat)

	// VT: NOTE: We can safely assume that since someone's created an instance,
	// they're going to use it

	rpio.Open()

	return automationHatPi{hat}
}

func initialize(hat *automationHatBase) {

	// Pinout: https://pinout.xyz/pinout/automation_hat#

	// Input and output numbers are BCM pin numbers
	// LED numbers are from a different namespace (SN3218 PWM driver)

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

	hat.adc33 = GetADC33(3, 3.3)

	hat.status = statusLights{power: GetLED(17), comms: GetLED(16), warn: GetLED(15)}
}

func (hat automationHatFake) Close() error {

	return nil
}

func (hat automationHatPi) Close() error {

	return rpio.Close()
}

func (hat automationHatBase) Relay() [3]Relay {
	return hat.relay
}

func (hat automationHatBase) ADC24() [3]ADC24 {
	return hat.adc24
}

func (hat automationHatBase) Input() [3]Input {
	return hat.input
}

func (hat automationHatBase) Output() [3]Output {
	return hat.output
}

func (hat automationHatBase) ADC33() ADC33 {
	return hat.adc33
}

func (hat automationHatBase) StatusLights() StatusLights {
	return hat.status
}

func (sl statusLights) Power() Light {
	return sl.power
}

func (sl statusLights) Comms() Light {
	return sl.comms
}

func (sl statusLights) Warn() Light {
	return sl.warn
}
