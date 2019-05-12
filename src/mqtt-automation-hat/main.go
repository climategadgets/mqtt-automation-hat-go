// +build go1.11

package main

import (
	"flag"
	"fmt"
	"github.com/climategadgets/mqtt-automation-hat-go/src/automation-hat"
	"github.com/climategadgets/mqtt-automation-hat-go/src/cf"
	"github.com/climategadgets/mqtt-automation-hat-go/src/protocol-handler"
	"github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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

	var logger *zap.Logger

	if *debug {
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		logger, _ = config.Build()
	} else {
		logger, _ = zap.NewProduction()
	}
	zap.ReplaceGlobals(logger)
	defer logger.Sync()

	target := fmt.Sprintf("tcp://%s:%d", *brokerHost, *brokerPort)
	topicFilter := fmt.Sprintf("%s/#", *rootTopic)

	zap.S().Infof("connecting to broker: %v", target)
	zap.S().Infof("topic filter: %v", topicFilter)
	zap.S().Infof("configuration: %v", *config)

	config := cf.ReadConfig(*config)
	automationHat := automation_hat.GetAutomationHAT()
	mqttClient := initMqttClient(target, topicFilter, automationHat, config)

	installShutDownHandler(mqttClient, automationHat)

	// VT: NOTE: Now we wait until we're interrupted

	select {}
}

// Create MQTT client and subscribe
func initMqttClient(target string, topicFilter string, automationHat automation_hat.AutomationHAT, config cf.ConfigHAT) mqtt.Client {

	opts := mqtt.NewClientOptions().AddBroker(target).SetClientID("mqtt-automation-hat").SetCleanSession(true)

	result := mqtt.NewClient(opts)

	if token := result.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	receive := func(client mqtt.Client, message mqtt.Message) {
		protocol_handler.Receive(client, message, automationHat, config)
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
			zap.S().Infow("shutting down", "signal", sig)

			var done sync.WaitGroup
			done.Add(2)

			now := time.Now()

			// Shut down MQTT client
			go func() {
				defer done.Done()

				now := time.Now()

				mqttClient.Disconnect(250)

				duration := time.Since(now)

				zap.S().Infow("MQTT disconnected", "duration", duration)
			}()

			// Shut down the Automation HAT
			go func() {
				defer done.Done()

				now := time.Now()

				automationHat.Close()

				duration := time.Since(now)

				zap.S().Infow("AutomationHAT shut down", "duration", duration)
			}()

			done.Wait()
			duration := time.Since(now)

			zap.S().Infof("shut down in %v", duration)

			zap.S().Infof("bye")
			os.Exit(0)
		}
	}()
}
