package automation_hat

import (
	"fmt"
	"testing"
)

// Make sure only a single instance is created
func TestGetAutomationHAT(t *testing.T) {

	hat1 := GetAutomationHAT()
	hat2 := GetAutomationHAT()

	// Pointers must be the same

	t.Log(fmt.Sprintf("hat1: %d", hat1))
	t.Log(fmt.Sprintf("hat2: %d", hat2))
}

// Make sure lights and relays behave as switches and keep state
func TestRelaySwitch(t *testing.T) {

	hat := GetAutomationHAT()

	testSwitch(t, hat.Relay()[0])
}

// Make sure lights and relays behave as switches and keep state
func TestLightSwitch(t *testing.T) {

	hat := GetAutomationHAT()

	testSwitch(t, hat.StatusLights().Power())
}

func testSwitch(t *testing.T, s Switch) {

	state := s.Get()
	changed := s.Set(!state)

	if !changed {
		t.Fatalf(fmt.Sprintf("failed to change state from %t: %s", state, s))
	}
}
