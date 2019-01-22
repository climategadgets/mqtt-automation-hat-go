package automation_hat

type light struct {
	pin        int
	state      bool
	brightness float64
}

func (l light) Get() bool {
	return l.state
}

func (l light) Set(state bool) bool {

	changed := l.state == state

	// VT: FIXME: Need to implement the state change, though
	l.state = state

	return changed
}

func GetLED(pin int) Light {

	return light{pin: pin}
}

type led_container struct {
	led Light
}

func (l led_container) Light() Light {
	return l.led
}
