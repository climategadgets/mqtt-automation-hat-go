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
	"strings"
)

// Parse a JSON message.
// First argument is the payload to parse, second is the topic name
func parse(source []byte, topic string) interface{} {

	guessed := parseGuess(source, topic)

	if guessed != nil {
		return guessed
	}

	var entity hcc_shared.HccMessageEntity
	err := json.Unmarshal(source, &entity)

	if err != nil {
		zap.S().Warnf("can't parse even as a entity: %v", entity)
	}

	switch entity.Type {

	case hcc_shared.TypeSensor:
		return parseSensor(source, topic)
	case hcc_shared.TypeSwitch:
		return parseSwitch(source, topic)
	case hcc_shared.TypeZone:
		return parseZone(source, topic)
	default:
		zap.S().Warnf("unknown entity type: '%v', entityType missing from JSON? %s", entity.Type, string(source))
		return nil
	}
}

// Parse a JSON message.
// First argument is the payload to parse, second is the topic name
func parseGuess(source []byte, topic string) interface{} {

	if strings.Contains(topic, string(hcc_shared.TypeSensor)) {
		zap.S().Debugf("guessing this is a sensor: %s", topic)
		return parseSensor(source, topic)
	}

	if strings.Contains(topic, string(hcc_shared.TypeSwitch)) {
		zap.S().Debugf("guessing this is a switch: %s", topic)
		return parseSwitch(source, topic)
	}

	if strings.Contains(topic, string(hcc_shared.TypeZone)) {
		zap.S().Debugf("guessing this is a zone: %s", topic)
		return parseZone(source, topic)
	}

	return nil
}

// Parse a sensor JSON message.
// First argument is the payload to parse, second is the topic name (for debugging purposes only)
func parseSensor(source []byte, topic string) interface{} {

	var payload hcc_shared.HccMessageSensor
	err := json.Unmarshal(source, &payload)

	if err != nil {
		zap.S().Errorf("malformed sensor (%v): %v", err, source)
		return nil
	}

	if payload.Signal == nil {
		zap.S().Error("not a sensor (no signal): %v", source)
		return nil
	}

	zap.S().Debugf("receive: sensor: %v = %v", topic, *payload.Signal)
	return payload
}

// Parse a switch JSON message.
// First argument is the payload to parse, second is the topic name (for debugging purposes only)
func parseSwitch(source []byte, topic string) interface{} {

	var payload hcc_shared.HccMessageSwitch
	err := json.Unmarshal(source, &payload)

	if err != nil {
		zap.S().Errorf("malformed switch (%v): %v", err, source)
		return nil
	}

	if payload.State == nil {
		zap.S().Error("not a switch (no state): %v", source)
		return nil
	}

	zap.S().Debugf("receive: switch: %v = %v", topic, *payload.State)
	return payload
}

// Parse a zone JSON message.
// First argument is the payload to parse, second is the topic name (for debugging purposes only)
func parseZone(source []byte, topic string) interface{} {

	var payload hcc_shared.HccMessageZone
	err := json.Unmarshal(source, &payload)

	if err != nil {
		zap.S().Errorf("malformed zone (%v): %v", err, source)
		return nil
	}

	if payload.ThermostatSignal == nil {
		zap.S().Debug("not a zone: no thermostat signal")
		return nil
	}

	zap.S().Debugf("receive: zone: %v = %v", topic, payload)
	return payload
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
