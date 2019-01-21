// +build go1.11

package main

import (
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/eclipse/paho.mqtt.golang"
)

var (
	debug      = flag.Bool("debug", false, "enable verbose logging")
	brokerHost = flag.String("broker-host", "localhost", "MQTT broker host to use")
	brokerPort = flag.Int("broker-port", 1883, "MQTT broker port to use")
	rootTopic  = flag.String("root-topic", "", "Root topic to listen to")
)
func main() {
	flag.Parse()

	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	target := fmt.Sprintf("tcp://%s:%d", *brokerHost, *brokerPort)
	topicFilter := fmt.Sprintf("%s/#", *rootTopic)

	log.Info("connecting to broker: " + target)
	log.Info("topic filter: " + topicFilter)

	opts := mqtt.NewClientOptions().AddBroker(target).SetClientID("mqtt-automation-hat").SetCleanSession(true)

	c := mqtt.NewClient(opts)

	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	if token := c.Subscribe(topicFilter, 2, receive); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	defer c.Disconnect(250);

	// VT: NOTE: Now we wait until we're interrupted

	select {}
}

func receive(client mqtt.Client, message mqtt.Message) {
	log.Info(fmt.Sprintf("%s %s", message.Topic(), message.Payload()))


}