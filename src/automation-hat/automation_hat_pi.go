// +build arm

package automation_hat

import (
	"fmt"
	"github.com/stianeikeland/go-rpio"
	"go.uber.org/zap"
	"time"
)

func newAutomationHAT() AutomationHAT {

	zap.S().Info("creating new instance of AutomationHAT")
	hatBase := automationHatBase{}
	initialize(&hatBase)

	// VT: NOTE: We can safely assume that since someone's created an instance,
	// they're going to use it

	if err := rpio.Open(); err != nil {

		// VT: NOTE: It makes no sense to continue, just bail out
		panic(fmt.Sprintf("can't open rpio, reason: %v", err))
	}

	hat := automationHatPi{hatBase}

	go func(control <-chan interface{}) {

		for {
			select {
			case m, ok := <-control:

				if !ok {
					zap.S().Errorw("control/pi channel closed?")
					break
				}
				// VT: FIXME: Errorw so it is visible in the log
				zap.S().Errorw("control/rpio", "message", m)
				execute(m)
			}
		}

	}(hat.control)

	// Now that we have the hardware listener in place, we can reset the board state
	reset(&hat)

	zap.S().Info("init: giving the board a chance to settle...")
	time.Sleep(500 * time.Millisecond)
	zap.S().Info("init: done")

	return hat
}

func (hat automationHatPi) Close() error {

	reset(&hat)

	zap.S().Info("close(): giving the board a chance to settle...")
	time.Sleep(250 * time.Millisecond)
	zap.S().Info("close(): done")

	return rpio.Close()
}

// Executes the request
func execute(message interface{}) {

	switch command := message.(type) {
	case relayCommand:
		executeRelay(command)
	case lightCommand:
		executeLight(command)
	default:
		zap.S().Errorw("don't know how to execute", "message", message)
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

func executeLight(command lightCommand) {

	zap.S().Errorf("FIXME: implement %v", command)
}

func reset(pi *automationHatPi) {

	for _, relay := range pi.relay {
		relay.Set(false)
	}

	// VT: FIXME: Reset the lights and other things
}
