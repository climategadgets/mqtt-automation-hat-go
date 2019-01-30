package protocol_handler

import (
	"encoding/json"
	"fmt"
	"github.com/climategadgets/mqtt-automation-hat-go/src/automation-hat"
	"github.com/climategadgets/mqtt-automation-hat-go/src/hcc-shared"
	"github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
	"sync"
)

func Receive(client mqtt.Client, message mqtt.Message, automationHat automation_hat.AutomationHAT) {
	log.Debug(fmt.Sprintf("%s %s", message.Topic(), message.Payload()))

	// Let's allow some redundancy and try parsing the payload for several types we know about in parallel,
	// then discard the ones we don't care about (or the ones that simply didn't parse)
	var done sync.WaitGroup
	done.Add(2)

	// Zero or one goroutine will set this to true
	parsed := false

	go func() {
		defer done.Done()

		var payload hcc_shared.HccMessageSwitch
		err := json.Unmarshal(message.Payload(), &payload)

		if err != nil {
			log.Debug("not a switch: ", err)
			return
		}

		log.Info("switch: ", payload)
		parsed = true
	}()

	go func() {
		defer done.Done()

		var payload hcc_shared.HccMessageSensor
		err := json.Unmarshal(message.Payload(), &payload)

		if err != nil {
			log.Debug("not a sensor: ", err)
			return
		}

		log.Info("sensor: ", payload)
		parsed = true
	}()

	done.Wait()
	log.Debug("all parsers returned")

	if !parsed {
		log.Warn(fmt.Sprintf("no parsers were able to understand: %s", message.Payload()))
	}
}
