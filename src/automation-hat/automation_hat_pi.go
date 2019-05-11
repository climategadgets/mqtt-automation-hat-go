package automation_hat

import (
	"fmt"
	"github.com/stianeikeland/go-rpio"
	"go.uber.org/zap"
)

func newAutomationHAT() AutomationHAT {

	zap.S().Info("creating new instance of AutomationHAT")
	hat := automationHatBase{}
	initialize(&hat)

	// VT: NOTE: We can safely assume that since someone's created an instance,
	// they're going to use it

	if err := rpio.Open(); err != nil {

		// VT: NOTE: It makes no sense to continue, just bail out
		panic(fmt.Sprintf("can't open rpio, reason: %v", err))
	}

	go func(control <-chan interface{}) {

		for {
			select {
			case m, ok := <-control:

				if !ok {
					zap.S().Errorf("control/pi channel closed?")
					break
				}
				// VT: FIXME: Errorf so it is visible in the log
				zap.S().Errorf("control/rpio: %v", m)
				execute(m)
			}
		}

	}(hat.control)

	return automationHatPi{hat}
}

func (hat automationHatPi) Close() error {
	return rpio.Close()
}

// Executes the request
func execute(message interface{}) {

	switch command := message.(type) {
	case relayCommand:
		executeRelay(command)
	default:
		zap.S().Errorf("don't know how to execute %v", message)
	}
}

func executeRelay(command relayCommand) {

	pin := rpio.Pin(command.pin)
	pin.Output()

	if command.state {
		pin.High()
	} else {
		pin.Low()
	}
}
