package automation_hat

import "go.uber.org/zap"

type light struct {
	messageBus
	pin        uint8
	state      bool
	brightness float64
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

	zap.S().Infof("light: pin=%v, state=%v, brightness=%v, changed=%v", l.pin, l.state, l.brightness, changed)

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
