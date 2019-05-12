// +build !arm

package automation_hat

import "go.uber.org/zap"

func newAutomationFake() AutomationHAT {

	zap.S().Warn("using AutomationHAT fake")
	hat := automationHatBase{}
	initialize(&hat)

	go func(control <-chan interface{}) {

		for {
			select {
			case m, ok := <-control:

				if !ok {
					zap.S().Errorw("control/fake channel closed?")
					break
				}
				// VT: NOTE: This is all we do here in the fake, log.
				// VT: FIXME: Errorw so it is visible in the log
				zap.S().Errorw("control/fake", "message", m)
			}
		}

	}(hat.control)

	return automationHatFake{hat}
}

func (hat automationHatFake) Close() error {
	return nil
}
