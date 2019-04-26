package protocol_handler

import (
	"encoding/json"
	"fmt"
	"github.com/climategadgets/mqtt-automation-hat-go/src/automation-hat"
	"github.com/climategadgets/mqtt-automation-hat-go/src/hcc-shared"
	"github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
	"sync"
	"sync/atomic"
)

func Receive(client mqtt.Client, message mqtt.Message, automationHat automation_hat.AutomationHAT) {
	log.Debug(fmt.Sprintf("%s %s", message.Topic(), message.Payload()))

	// Let's allow some redundancy and try parsing the payload for several types we know about in parallel,
	// then discard the ones we don't care about (or the ones that simply didn't parse)
	var done sync.WaitGroup
	done.Add(2)

	// This will contain the number of successfully parsed messages. Default JSON parser is quite forgiving
	// and may yield false positives, so it is a responsibility of each parser function to ensure they see
	// the message they *really* expect. If more than one message is parsed successfully, it's a problem -
	// we'll complain and ignore it.
	parsed := uint32(0)

	go func() {
		defer done.Done()

		var payload hcc_shared.HccMessageSwitch
		err := json.Unmarshal(message.Payload(), &payload)

		if err != nil {
			log.Debug("not a switch: ", err)
			return
		}

		log.Infof("switch: %v %v", message.Topic(), payload)
		atomic.AddUint32(&parsed, 1)
	}()

	go func() {
		defer done.Done()

		var payload hcc_shared.HccMessageSensor
		err := json.Unmarshal(message.Payload(), &payload)

		if err != nil {
			log.Debug("not a sensor: ", err)
			return
		}

		log.Infof("sensor: %v %v", message.Topic(), payload)
		atomic.AddUint32(&parsed, 1)
	}()

	done.Wait()
	log.Debug("all parsers returned")

	if parsed == 0 {
		log.Warn(fmt.Sprintf("no parsers were able to understand: %s", message.Payload()))
		return
	}

	if parsed > 1 {
		log.Error(fmt.Sprintf("ambiguous message, parsed %d times: %s", parsed, message.Payload()))
	}
}
