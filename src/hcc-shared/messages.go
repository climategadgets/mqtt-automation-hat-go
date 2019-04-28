// Package hcc_shared contains data structures that are common for
// all incoming and outgoing communications between HCC and the outside world
package hcc_shared

type HccMessageBase struct {
	Timestamp uint64 `json:"timestamp"`
	Name      string `json:"name"`
	Signature string `json:"signature"`
}

type HccMessageUnique struct {
	Id string `json:"id"`
}

type HccMessageResponse struct {
	ReplyTo HccMessageUnique `json:"reply-to"`
}

type HccMessageSensor struct {
	HccMessageBase
	Signal *float64 `json:"signal"`
}

type HccMessageSwitch struct {
	HccMessageBase
	HccMessageUnique
	State *bool `json:"state"`
}
