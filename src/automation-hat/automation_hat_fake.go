package automation_hat

import "go.uber.org/zap"

func newAutomationFake() AutomationHAT {

	zap.S().Warn("using AutomationHAT fake")
	hat := automationHatBase{}
	initialize(&hat)

	go func(control <-chan interface{}) {

		for {
			select {
			case m := <-control:
				// VT: NOTE: This is all we do here in the fake, log.
				// VT: FIXME: Errorf so it is visible in the log
				zap.S().Errorf("control/fake: %v", m)
			}
		}

	}(hat.control)

	return automationHatFake{hat}
}

func (hat automationHatFake) Close() error {
	return nil
}
