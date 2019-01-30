package cf

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/user"
)

type SwitchConfig struct {
	topic    string
	inverted bool
}

type ConfigSwitchMap map[string]SwitchConfig

func GetDefaultConfigLocation() string {

	usr, err := user.Current()

	if err != nil {
		log.Fatal(err)
	}

	return usr.HomeDir + "/.mqtt-automation-hat/config.json"
}

func ReadConfig(source string) ConfigSwitchMap {

	return readConfig(source, false)
}

// Read the configuration from given source file.
// If 'again' is true, that means that we've just tried to create a default configuration, and should fail immediately.
func readConfig(source string, again bool) ConfigSwitchMap {

	// Let's first make sure that the file is there, readable and parseable
	jsonFile, err0 := os.Open(source)

	if err0 != nil {

		// VT: FIXME: Learn to distinguish error types
		defaultCf := GetDefaultConfigLocation()

		if source == defaultCf {

			if again {
				log.Fatal("looks like we failed to create default config, terminating")
			}

			// VT: FIXME: Let's assume the file doesn't exist
			log.Warnf("default configuration location (" + defaultCf + ") specified, but the file can't be read")
			log.Warn("we'll try to create a default configuration there,")
			log.Warn("see https://github.com/climategadgets/mqtt-automation-hat-go/wiki/Configuration-File for details")

			createDefaultConfig(defaultCf)
			return readConfig(source, true)
		}
	}

	buffer, err1 := ioutil.ReadAll(jsonFile)

	if err1 != nil {
		log.Fatalf("couldn't read configuration from "+source+", %v", err1)
	}

	var result ConfigSwitchMap
	err2 := json.Unmarshal(buffer, &result)

	if err2 != nil {
		log.Fatalf("couldn't read configuration from "+source+", %v", err2)
	}

	log.Warnf("configuration: %v", result)

	return result
}

func createDefaultConfig(source string) {

	log.Warn("FIXME: createDefaultConfig()")
}
