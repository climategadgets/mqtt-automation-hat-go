package protocol_handler

import (
	"encoding/json"
	"fmt"
	"github.com/climategadgets/mqtt-automation-hat-go/src/automation-hat"
	"github.com/climategadgets/mqtt-automation-hat-go/src/cf"
	"github.com/climategadgets/mqtt-automation-hat-go/src/hcc-shared"
	"github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
	"strconv"
	"sync"
	"sync/atomic"
)

// Parse a JSON message.
// First argument is the payload to parse, second is the topic name (for debugging purposes only)
func parse(source []byte, topic string) interface{} {

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
		err := json.Unmarshal(source, &payload)

		if err != nil {
			zap.S().Debugf("not a sensor: %v", err)
			return
		}

		if payload.Signal == nil {
			zap.S().Debug("not a sensor: no signal")
			return
		}

		zap.S().Debugf("receive: sensor: %v = %v", topic, *payload.Signal)
		c <- payload
		atomic.AddUint32(&parsed, 1)
	}(cSensor)

	go func(c chan<- hcc_shared.HccMessageSwitch) {
		defer done.Done()

		var payload hcc_shared.HccMessageSwitch
		err := json.Unmarshal(source, &payload)

		if err != nil {
			zap.S().Debugf("not a switch: %v", err)
			return
		}

		if payload.State == nil {
			zap.S().Debug("not a switch: no state")
			return
		}

		zap.S().Debugf("receive: switch: %v = %v", topic, *payload.State)
		c <- payload
		atomic.AddUint32(&parsed, 1)
	}(cSwitch)

	go func(c chan<- hcc_shared.HccMessageZone) {
		defer done.Done()

		var payload hcc_shared.HccMessageZone
		err := json.Unmarshal(source, &payload)

		if err != nil {
			zap.S().Debugf("not a zone: ", err)
			return
		}

		if payload.ThermostatSignal == nil {
			zap.S().Debug("not a zone: no thermostat signal")
			return
		}

		zap.S().Debugf("receive: zone: %v = %v", topic, payload)
		c <- payload
		atomic.AddUint32(&parsed, 1)
	}(cZone)

	done.Wait()
	zap.S().Debug("all parsers returned")

	if parsed == 0 {
		zap.S().Warnf("no parsers were able to understand: %s = %s", topic, source)
		return nil
	}

	if parsed > 1 {
		zap.S().Errorf("ambiguous message, parsed %d times: %s = %s", parsed, topic, source)
		return nil
	}

	// We are only supposed to end up here if there was exactly one message parsed

	select {
	case m := <-cSensor:
		return m

	case m := <-cSwitch:
		//process(message.Topic(), m, automationHat, config)
		return m

	case m := <-cZone:
		return m
	}
}

// Receive an MQTT message.
func Receive(client mqtt.Client, message mqtt.Message, automationHat automation_hat.AutomationHAT, config cf.ConfigHAT) {
	zap.S().Debugf("%s %s", message.Topic(), message.Payload())

	command := parse(message.Payload(), message.Topic())

	if command == nil {
		// Nothing parsed
		return
	}

	switch c := command.(type) {

	case hcc_shared.HccMessageSensor:
		zap.S().Infof("sensor: %v", c)
		// VT: FIXME: Implement this

	case hcc_shared.HccMessageSwitch:
		processSwitch(message.Topic(), c, automationHat, config)

	case hcc_shared.HccMessageZone:
		zap.S().Infof("zone: %v", c)
		// VT: FIXME: Implement this

	default:
		panic(fmt.Sprintf("unknown type: %v", c))
	}
}

func processSwitch(topic string, m hcc_shared.HccMessageSwitch, automationHat automation_hat.AutomationHAT, config cf.ConfigHAT) {

	zap.S().Infof(fmt.Sprintf("process: switch: %v = %v", topic, *m.State))

	switchMap := config["switchMap"]
	for switchId, switchConfig := range switchMap {

		if switchConfig.Topic == topic {

			var state bool

			if switchConfig.Inverted {
				state = !*m.State
			} else {
				state = *m.State
			}

			zap.S().Infof("topic=%v, id=%v, inverted=%v, state=%v", topic, switchId, switchConfig.Inverted, state)

			switchOffset, _ := strconv.Atoi(switchId)
			automationHat.Relay()[switchOffset].Set(state)

			return
		}
	}

	zap.S().Warnf("no matching topic in the config: %v", topic)
}
