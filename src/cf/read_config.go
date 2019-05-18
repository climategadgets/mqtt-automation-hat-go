package cf

import (
	"encoding/json"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"os/user"
)

type SwitchConfig struct {
	Topic    string `json:"topic"`
	Inverted bool   `json:"inverted"`
}

var defaultCfDir = ".mqtt-automation-hat"
var defaultCfFile = "config.json"

// The key is the switch ID, the value is the switch config
type ConfigSwitchMap map[string]SwitchConfig
type ConfigHAT map[string]ConfigSwitchMap

func getDefaultConfigDir() string {

	usr, err := user.Current()

	if err != nil {
		zap.S().Fatal(err)
	}

	return usr.HomeDir + string(os.PathSeparator) + defaultCfDir
}

func GetDefaultConfigLocation() string {

	return getDefaultConfigDir() + string(os.PathSeparator) + defaultCfFile
}

func ReadConfig(source string) ConfigHAT {

	return readConfig(source, false)
}

// Read the configuration from given source file.
// If 'again' is true, that means that we've just tried to create a default configuration, and should fail immediately.
func readConfig(source string, again bool) ConfigHAT {

	// Let's first make sure that the file is there, readable and parseable
	jsonFile, err0 := os.Open(source)

	if err0 != nil {

		// VT: FIXME: Learn to distinguish error types
		defaultCf := GetDefaultConfigLocation()

		if source == defaultCf {

			if again {
				zap.S().Fatal("looks like we failed to create default config, terminating")
			}

			// VT: FIXME: Let's assume the file doesn't exist
			zap.S().Warnf("default configuration location (" + defaultCf + ") specified, but the file can't be read")
			zap.S().Warn("we'll try to create a default configuration there,")
			zap.S().Warn("see https://github.com/climategadgets/mqtt-automation-hat-go/wiki/Configuration-File for details")

			createDefaultConfig(defaultCf)
			return readConfig(source, true)
		}
	}

	buffer, err1 := ioutil.ReadAll(jsonFile)

	if err1 != nil {
		zap.S().Fatalf("couldn't read configuration from "+source+", %v", err1)
	}

	result := make(ConfigHAT)
	err2 := json.Unmarshal(buffer, &result)

	if err2 != nil {
		zap.S().Fatalf("couldn't read configuration from "+source+", %v", err2)
	}

	zap.S().Infof("configuration: %v", result)

	return result
}

var defaultConfig = `{
  "switchMap": {
    "0": {
      "topic": "/pimoroni/automation-hat/switch/0",
      "inverted": false
    },
    "1": {
      "topic": "/pimoroni/automation-hat/switch/1",
      "inverted": false
    },
    "2": {
      "topic": "/pimoroni/automation-hat/switch/2",
      "inverted": false
    }
  }
}`

func createDefaultConfig(target string) {

	cfDir := getDefaultConfigDir()

	_, err := os.Stat(cfDir)

	if os.IsNotExist(err) {
		err := os.Mkdir(cfDir, 0755)

		if err != nil {
			zap.S().Fatalf("can't create directory for default configuration: %v", err)
		}
	}

	{
		payload := []byte(defaultConfig)
		err := ioutil.WriteFile(target, payload, 0644)

		if err != nil {
			zap.S().Fatal(err)
		}
	}
}
