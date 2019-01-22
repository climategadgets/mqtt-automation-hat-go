// Package hcc_shared contains data structures that are common for
// all incoming and outgoing communications between HCC and the outside world
package hcc_shared

type hcc_message_base struct {
	Timestamp string `json:"timestamp"`
	Name      string `json:"name"`
	Signature string `json:"signature"`
}

type hcc_message_unique struct {
	Id string `json:"id"`
}

type hcc_message_response struct {
	ReplyTo hcc_message_unique `json:"reply-to"`
}

type hcc_message_sensor struct {
	hcc_message_base
	Signal float64 `json:"signal"`
}

type hcc_message_switch struct {
	hcc_message_base
	hcc_message_unique
	State bool `json:"state"`
}
