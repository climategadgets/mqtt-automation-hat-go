// Package hcc_shared contains data structures that are common for
// all incoming and outgoing communications between HCC and the outside world
package hcc_shared

type EntityType string

const (
	TypeSensor   EntityType = "sensor"
	TypeSwitch     EntityType = "switch"
	TypeZone EntityType = "thermostat"
)

type HccMessageBase struct {
	Type      EntityType `json:"entityType""`
	Timestamp uint64     `json:"timestamp"`
	Name      string     `json:"name"`
	Signature string     `json:"signature"`
}

type HccMessageUnique struct {
	Id string `json:"id"`
}

type HccMessageResponse struct {
	ReplyTo HccMessageUnique `json:"reply-to"`
}

type HccMessageSensor struct {
	HccMessageBase
	// Need this as a pointer to catch ambiguous JSON parser output
	Signal *float64 `json:"signal"`
}

type HccMessageSwitch struct {
	HccMessageBase
	HccMessageUnique
	// Need this as a pointer to catch ambiguous JSON parser output
	State *bool `json:"state"`
}

type HvacMode int8

const (
	ModeCooling HvacMode = -1
	ModeOff     HvacMode = 0
	ModeHeating HvacMode = 1
)

type ZoneState string

const (
	ZoneError   ZoneState = "ERROR"
	ZoneOff     ZoneState = "OFF"
	ZoneCalling ZoneState = "CALLING"
	ZoneHappy   ZoneState = "HAPPY"
)

type HccMessageZone struct {
	HccMessageBase
	CurrentTemperature  *float64 `json:currentTemperature`
	DeviationEnabled    bool     `json:deviation.enabled`
	DeviationSetpoint   float64  `json:deviation.setpoint`
	DeviationVoting     bool     `json:deviation.voting`
	Enabled             bool     `json:enabled`
	HvacMode            `json:mode`
	OnHold              bool      `json:onHold`
	SetpointTemperature float64   `json:setpointTemperature`
	State               ZoneState `json:state`
	// Need this as a pointer to catch ambiguous JSON parser output
	ThermostatSignal *float64 `json:thermostatSignal`
	Voting           bool     `json:voting`
}
