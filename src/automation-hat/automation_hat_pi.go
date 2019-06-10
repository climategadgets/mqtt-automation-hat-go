// +build arm

package automation_hat

import (
	"fmt"
	sn32182 "github.com/climategadgets/mqtt-automation-hat-go/src/sn3218"
	"github.com/stianeikeland/go-rpio"
	"go.uber.org/zap"
	"time"
)

const (
	ledIntensity = 0x33 // LEDs are extremely bright even in daylight
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

	ledDriver := sn32182.GetSN3218()

	if ledDriver == nil {
		panic("can't open the LED driver, see the logs for the cause")
	}

	hat := automationHatPi{hatBase, ledDriver}

	go func(control <-chan interface{}) {

		for {
			select {
			case m, ok := <-control:

				if !ok {
					zap.S().Errorw("control/pi channel closed?")
					break
				}
				// VT: FIXME: Errorw so it is visible in the log
				zap.S().Errorw("control/rpio", "message", fmt.Sprintf("%v", m))
				execute(hat, m)
			}
		}

	}(hat.control)

	// Now that we have the hardware listener in place, we can reset the board state
	reset(&hat)

	ledDriver.Enable(true)
	ledDriver.EnableLEDs(0x3FFFF) // 0b111111111111111111, all of them

	zap.S().Info("init: giving the board a chance to settle...")

	// ...and to make sure all LEDs are functional
	for channel := 0; channel < 18; channel++ {
		ledDriver.SetLED(byte(channel), ledIntensity)
		time.Sleep(10 * time.Millisecond)
	}

	time.Sleep(400 * time.Millisecond)

	for channel := 0; channel < 18; channel++ {
		ledDriver.SetLED(byte(channel), 0)
	}

	time.Sleep(100 * time.Millisecond)

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
func execute(hat automationHatPi, message interface{}) {

	switch command := message.(type) {
	case relayCommand:
		executeRelay(hat, command)
	case lightCommand:
		executeLight(hat, command)
	case adcCommand:
		executeAdc(hat, command)
	default:
		zap.S().Errorw("don't know how to execute", "message", message)
	}
}

func executeRelay(hat automationHatPi, command relayCommand) {

	zap.S().Debugw("executeRelay", "pin", command.pin, "state", command.state)
	pin := rpio.Pin(command.pin)
	pin.Output()

	if command.state {
		pin.High()
		hat.ledDriver.SetLED(command.ledNC, 0)
		hat.ledDriver.SetLED(command.ledNO, ledIntensity)
	} else {
		pin.Low()
		hat.ledDriver.SetLED(command.ledNC, ledIntensity)
		hat.ledDriver.SetLED(command.ledNO, 0)
	}
}

func executeLight(hat automationHatPi, command lightCommand) {

	zap.S().Debugw("executeLight", "pin", command.pin, "state", command.state)
	if command.state {
		hat.ledDriver.SetLED(command.pin, ledIntensity)
	} else {
		hat.ledDriver.SetLED(command.pin, 0)
	}
}

func executeAdc(hat automationHatPi, command adcCommand) {

	zap.S().Debugw("executeAdc", "channel", command.channel)

	// VT: FIXME: Need to actually read the value
	*command.signal = 42.0
	command.done.Done()

	zap.S().Errorw("executeAdc: FIXME", "signal", *command.signal)
}

func reset(hat *automationHatPi) {

	// LED driver goes first because relays use LEDs to indicate their status
	hat.ledDriver.Reset()

	for _, relay := range hat.relay {
		relay.Set(false)
	}

	// VT: FIXME: Reset other things
}
