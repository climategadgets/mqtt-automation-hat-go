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
	done.Add(3)

	// This will contain the number of successfully parsed messages. Default JSON parser is quite forgiving
	// and may yield false positives, so it is a responsibility of each parser function to ensure they see
	// the message they *really* expect. If more than one message is parsed successfully, it's a problem -
	// we'll complain and ignore it.
	parsed := uint32(0)

	// Only one of these channels will eventually receive a parsed object
	var cSensor chan hcc_shared.HccMessageSensor = make(chan hcc_shared.HccMessageSensor, 1)
	var cSwitch chan hcc_shared.HccMessageSwitch = make(chan hcc_shared.HccMessageSwitch, 1)
	var cZone chan hcc_shared.HccMessageZone = make(chan hcc_shared.HccMessageZone, 1)

	defer close(cSensor)
	defer close(cSwitch)
	defer close(cZone)

	go func(c chan<- hcc_shared.HccMessageSensor) {
		defer done.Done()

		var payload hcc_shared.HccMessageSensor
		err := json.Unmarshal(message.Payload(), &payload)

		if err != nil {
			log.Debug("not a sensor: ", err)
			return
		}

		if payload.Signal == nil {
			log.Debug("not a sensor: no signal")
			return
		}

		log.Debugf("receive: sensor: %v = %v", message.Topic(), *payload.Signal)
		c <- payload
		atomic.AddUint32(&parsed, 1)
	}(cSensor)

	go func(c chan<- hcc_shared.HccMessageSwitch) {
		defer done.Done()

		var payload hcc_shared.HccMessageSwitch
		err := json.Unmarshal(message.Payload(), &payload)

		if err != nil {
			log.Debug("not a switch: ", err)
			return
		}

		if payload.State == nil {
			log.Debug("not a switch: no state")
			return
		}

		log.Debugf("receive: switch: %v = %v", message.Topic(), *payload.State)
		c <- payload
		atomic.AddUint32(&parsed, 1)
	}(cSwitch)

	go func(c chan<- hcc_shared.HccMessageZone) {
		defer done.Done()

		var payload hcc_shared.HccMessageZone
		err := json.Unmarshal(message.Payload(), &payload)

		if err != nil {
			log.Debug("not a zone: ", err)
			return
		}

		if payload.ThermostatSignal == nil {
			log.Debug("not a zone: no thermostat signal")
			return
		}

		log.Debugf("receive: zone: %v = %v", message.Topic(), payload)
		c <- payload
		atomic.AddUint32(&parsed, 1)
	}(cZone)

	done.Wait()
	log.Debug("all parsers returned")

	if parsed == 0 {
		log.Warn(fmt.Sprintf("no parsers were able to understand: %s = %s", message.Topic(), message.Payload()))
		return
	}

	if parsed > 1 {
		log.Error(fmt.Sprintf("ambiguous message, parsed %d times: %s = %s", parsed, message.Topic(), message.Payload()))
		return
	}

	// We are only supposed to end up here if there was exactly one message parsed

	select {
	case <-cSensor:
		// NOP for now
		return

	case m := <-cSwitch:
		process(message.Topic(), m)
		return

	case <-cZone:
		// NOP for now
		return
	}
}

func process(topic string, m hcc_shared.HccMessageSwitch) {

	log.Infof(fmt.Sprintf("process: switch: %v = %v", topic, *m.State))
}
