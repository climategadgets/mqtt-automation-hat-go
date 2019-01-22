package automation_hat

import (
	"fmt"
	"testing"
)

func TestGetAutomationHAT(t *testing.T) {

	hat1 := GetAutomationHAT()
	hat2 := GetAutomationHAT()

	// Pointers must be the same

	t.Log(fmt.Sprintf("hat1: %d", hat1))
	t.Log(fmt.Sprintf("hat2: %d", hat2))
}
