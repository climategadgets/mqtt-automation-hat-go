package automation_hat

import (
	"go.uber.org/zap"
)

type relay struct {
	messageBus
	state bool
	pin   uint8
	ledNO uint8
	ledNC uint8
	led   [2]Light
}

type relayCommand struct {
	relay
	changed bool
}

func GetRelay(control chan<- interface{}, pin uint8, ledNO uint8, ledNC uint8) Relay {

	r := &relay{messageBus: messageBus{control}, pin: pin, ledNO: ledNO, ledNC: ledNC}

	r.led[0] = GetLED(control, ledNO)
	r.led[1] = GetLED(control, ledNC)

	return r
}

func (r relay) Get() bool {
	return r.state
}

func (r *relay) Set(state bool) bool {

	if r.control == nil {
		panic("nil control channel")
	}

	changed := r.state != state

	r.state = state

	zap.S().Infof("relay: pin=%v, ledNO=%v, ledNC=%v, state=%v, changed=%v", r.pin, r.ledNO, r.ledNC, state, changed)

	// VT: NOTE: Counterintuitively, 'changed' is not always true. Remains to be seen how useful it is, though
	r.control <- relayCommand{*r, changed}

	if state {
		r.led[0].Set(true)
		r.led[1].Set(false)
	} else {
		r.led[0].Set(false)
		r.led[1].Set(true)
	}

	return changed
}

func (r relay) Light() [2]Light {
	return r.led
}
