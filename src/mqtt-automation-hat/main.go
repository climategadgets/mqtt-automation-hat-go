// +build go1.11

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/climategadgets/mqtt-automation-hat-go/src/automation-hat"
	"github.com/climategadgets/mqtt-automation-hat-go/src/hcc-shared"
	"github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"sync"
	"time"
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

	automationHat := automation_hat.GetAutomationHAT()
	mqttClient := initMqttClient(target, topicFilter, automationHat)

	installShutDownHandler(mqttClient)

	// VT: NOTE: Now we wait until we're interrupted

	select {}
}

// Create MQTT client and subscribe
func initMqttClient(target string, topicFilter string, automationHat automation_hat.AutomationHAT) mqtt.Client {

	opts := mqtt.NewClientOptions().AddBroker(target).SetClientID("mqtt-automation-hat").SetCleanSession(true)

	result := mqtt.NewClient(opts)

	if token := result.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	receive := func(client mqtt.Client, message mqtt.Message) {
		receive(client, message, automationHat)
	}

	if token := result.Subscribe(topicFilter, 2, receive); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	return result
}

func installShutDownHandler(mqttClient mqtt.Client) {

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	go func() {
		for sig := range sig {
			log.Infof("captured %v signal, shutting down", sig)

			var done sync.WaitGroup
			done.Add(2)

			now := time.Now()

			// Shut down MQTT client
			go func() {

				now := time.Now()

				mqttClient.Disconnect(250)

				duration := time.Since(now)

				log.Infof("MQTT disconnected in %v", duration)

				done.Done()
			}()

			// Shut down the Automation HAT
			go func() {

				// VT: FIXME: Need to actually do it

				done.Done()
			}()

			done.Wait()
			duration := time.Since(now)

			log.Infof("shut down in %v", duration)

			log.Infof("bye")
			os.Exit(0)
		}
	}()
}

func receive(client mqtt.Client, message mqtt.Message, automationHat automation_hat.AutomationHAT) {
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
