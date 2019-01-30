// +build go1.11

package main

import (
	"flag"
	"fmt"
	"github.com/climategadgets/mqtt-automation-hat-go/src/automation-hat"
	"github.com/climategadgets/mqtt-automation-hat-go/src/cf"
	"github.com/climategadgets/mqtt-automation-hat-go/src/protocol-handler"
	"github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"sync"
	"time"
)

var (
	config     = flag.String("config", cf.GetDefaultConfigLocation(), "Configuration file location")
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
	log.Info("configuration: " + *config)

	switchMap := cf.ReadConfig(*config)
	automationHat := automation_hat.GetAutomationHAT()
	mqttClient := initMqttClient(target, topicFilter, automationHat, switchMap)

	installShutDownHandler(mqttClient, automationHat)

	// VT: NOTE: Now we wait until we're interrupted

	select {}
}

// Create MQTT client and subscribe
func initMqttClient(target string, topicFilter string, automationHat automation_hat.AutomationHAT, switchMap cf.ConfigSwitchMap) mqtt.Client {

	opts := mqtt.NewClientOptions().AddBroker(target).SetClientID("mqtt-automation-hat").SetCleanSession(true)

	result := mqtt.NewClient(opts)

	if token := result.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	receive := func(client mqtt.Client, message mqtt.Message) {
		protocol_handler.Receive(client, message, automationHat)
	}

	if token := result.Subscribe(topicFilter, 2, receive); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	return result
}

func installShutDownHandler(mqttClient mqtt.Client, automationHat automation_hat.AutomationHAT) {

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
				defer done.Done()

				now := time.Now()

				mqttClient.Disconnect(250)

				duration := time.Since(now)

				log.Infof("MQTT disconnected in %v", duration)
			}()

			// Shut down the Automation HAT
			go func() {
				defer done.Done()

				now := time.Now()

				// VT: FIXME: Commented out until the protocol handler is finished - blows up
				// on i386
				// automationHat.Close()
				log.Error("FIXME: AutomationHat.Close() disabled")

				duration := time.Since(now)

				log.Infof("AutomationHAT shut down in %v", duration)
			}()

			done.Wait()
			duration := time.Since(now)

			log.Infof("shut down in %v", duration)

			log.Infof("bye")
			os.Exit(0)
		}
	}()
}
