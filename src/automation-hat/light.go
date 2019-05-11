package automation_hat

import "go.uber.org/zap"

type light struct {
	messageBus
	pin        uint8
	state      bool
	brightness byte
}

type lightCommand struct {
	light
	changed bool
}

func (l light) Get() bool {
	return l.state
}

func (l *light) Set(state bool) bool {

	if l.control == nil {
		panic("nil control channel")
	}

	changed := l.state != state

	l.state = state

	zap.S().Infow("set", "entityType", "light", "pin", l.pin, "state", l.state, "changed", changed)

	// VT: NOTE: Counterintuitively, 'changed' is not always true. Remains to be seen how useful it is, though
	l.control <- lightCommand{*l, changed}

	return changed
}

func GetLED(control chan<- interface{}, pin uint8) Light {
	return &light{messageBus: messageBus{control}, pin: pin}
}

type ledContainer struct {
	led Light
}

func (l ledContainer) Light() Light {
	return l.led
}
