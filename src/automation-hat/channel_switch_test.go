package automation_hat

import (
	"fmt"
	"testing"
)

// Make sure channel/switch combination works the way it is expected to
func TestChannelSwitch(t *testing.T) {

	var pipe chan interface{} = make(chan interface{}, 3)

	pipe <- relay{}
	pipe <- light{}
	pipe <- output{}

	close(pipe)

	for count := 0; count < 3; count++ {

		payload := <-pipe

		t.Logf("payload: %v", payload)

		switch tp := payload.(type) {
		case relay:
			t.Logf("relay: %v", tp)
		case light:
			t.Logf("light: %v", tp)
		case output:
			t.Logf("output: %v", tp)
		default:
			panic(fmt.Sprintf("unknown type: %v", tp))
		}
	}
}
