package automation_hat

import (
	"go.uber.org/zap"
)

type relay struct {
	state bool
	pin   uint8
	ledNO uint8
	ledNC uint8
	led   [2]Light
}

func GetRelay(pin uint8, ledNO uint8, ledNC uint8) Relay {

	r := &relay{pin: pin, ledNO: ledNO, ledNC: ledNC}

	r.led[0] = GetLED(ledNO)
	r.led[1] = GetLED(ledNC)

	return r
}

func (r relay) Get() bool {
	return r.state
}

func (r *relay) Set(state bool) bool {

	changed := r.state != state

	// VT: FIXME: Need to implement the state change, though
	r.state = state

	zap.S().Infof("relay: pin=%v, ledNO=%v, ledNC=%v, state=%v, changed=%v", r.pin, r.ledNO, r.ledNC, state, changed)
	return changed
}

func (r relay) Light() [2]Light {
	return r.led
}
