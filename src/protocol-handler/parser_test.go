package protocol_handler

import (
	"encoding/json"
	"github.com/climategadgets/mqtt-automation-hat-go/src/hcc-shared"
	"testing"
)

// VT: NOTE: Benchmarks show that the time spent on the operation doesn't depend as much on whether it passes or fails,
// but rather roughly on source JSON size. Still, there's a small overhead for a failed operation.

var switchJson string = "{\"entityType\":\"switch\",\"timestamp\":1547614624328,\"name\":\"mode switch\",\"signature\":\"mode switch\",\"state\":false,\"id\":\"af0c0e5802d625aaaddddb14ea7a8731\"}"
var sensorJson string = "{\"entityType\":\"thermostat\",\"timestamp\":1548739962980,\"name\":\"Back Wall\",\"signature\":\"4e4b070508820f5913f\",\"mode\":\"Cooling\",\"state\":\"HAPPY\",\"thermostatSignal\":-7.5,\"currentTemperature\":21.5,\"setpointTemperature\":30.0,\"enabled\":true,\"onHold\":false,\"voting\":true,\"deviation.setpoint\":0.0,\"deviation.enabled\":false,\"deviation.voting\":false}"

func BenchmarkParserSwitchPass(b *testing.B) {

	for i := 0; i < b.N; i++ {
		var payload hcc_shared.HccMessageSwitch
		json.Unmarshal([]byte(switchJson), &payload)
	}
}

func BenchmarkParserSwitchFail(b *testing.B) {

	for i := 0; i < b.N; i++ {
		var payload hcc_shared.HccMessageSwitch
		json.Unmarshal([]byte(sensorJson), &payload)
	}
}

func BenchmarkParserSensorPass(b *testing.B) {

	for i := 0; i < b.N; i++ {
		var payload hcc_shared.HccMessageSensor
		json.Unmarshal([]byte(sensorJson), &payload)
	}
}

func BenchmarkParserSensorFail(b *testing.B) {

	for i := 0; i < b.N; i++ {
		var payload hcc_shared.HccMessageSensor
		json.Unmarshal([]byte(switchJson), &payload)
	}
}

func BenchmarkParserSpeedParallel(b *testing.B) {

	for i := 0; i < b.N; i++ {
		parseParallel([]byte(sensorJson), "test-topic")
		parseParallel([]byte(switchJson), "test-topic")
	}
}

func BenchmarkParserSpeedFromSeed(b *testing.B) {

	for i := 0; i < b.N; i++ {
		parseFromSeed([]byte(sensorJson), "test-topic")
		parseFromSeed([]byte(switchJson), "test-topic")
	}
}

func BenchmarkParserSpeedEntity(b *testing.B) {

	for i := 0; i < b.N; i++ {
		var payload hcc_shared.HccMessageEntity
		json.Unmarshal([]byte(sensorJson), &payload)
	}
}

func BenchmarkParserSpeedZone(b *testing.B) {

	for i := 0; i < b.N; i++ {
		var payload hcc_shared.HccMessageSensor
		json.Unmarshal([]byte(sensorJson), &payload)
	}
}
