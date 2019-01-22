package automation_hat

type relay struct {
	state bool
	pin   int
	ledNO int
	ledNC int
	led   [2]Light
}

func GetRelay(pin int, ledNO int, ledNC int) Relay {

	r := relay{pin: pin, ledNO: ledNO, ledNC: ledNC}

	r.led[0] = GetLED(ledNO)
	r.led[1] = GetLED(ledNC)

	return r
}

func (r relay) Get() bool {
	return r.state
}

func (r relay) Set(state bool) bool {

	changed := r.state == state

	// VT: FIXME: Need to implement the state change, though
	r.state = state

	return changed
}

func (r relay) Light() [2]Light {
	return r.led
}
