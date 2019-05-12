package automation_hat

import (
	"fmt"
	"go.uber.org/zap"
	"testing"
)

// Make sure only a single instance is created
func TestGetAutomationHAT(t *testing.T) {

	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
	defer logger.Sync()

	hat1 := GetAutomationHAT()
	hat2 := GetAutomationHAT()

	// Pointers must be the same

	t.Log(fmt.Sprintf("hat1: %d", hat1))
	t.Log(fmt.Sprintf("hat2: %d", hat2))
}

// Make sure lights and relays behave as switches and keep state
func TestRelaySwitch(t *testing.T) {

	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
	defer logger.Sync()

	hat := GetAutomationHAT()

	testSwitch(t, hat.Relay()[0])
}

// Make sure lights and relays behave as switches and keep state
func TestLightSwitch(t *testing.T) {

	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
	defer logger.Sync()

	hat := GetAutomationHAT()

	testSwitch(t, hat.StatusLights().Power())
}

func testSwitch(t *testing.T, s Switch) {

	state := s.Get()
	changed := s.Set(!state)

	if !changed {
		t.Fatalf(fmt.Sprintf("failed to change state from %t: %s", state, s))
	}

	if s.Get() != !state {
		t.Fatalf(fmt.Sprintf("state didn't change from %t to %t for %s", state, !state, s))
	}
}
