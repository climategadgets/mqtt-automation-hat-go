package cf

import (
	"encoding/json"
	"testing"
)

func TestMarshalReadSwitchConfig(t *testing.T) {

	sc := SwitchConfig{"Topic0", true}
	buffer, err := json.Marshal(sc)

	if err != nil {
		t.Fatal(err)
	}

	t.Logf("switch config/marshal: %s", buffer)
}

func TestUnmarshalReadSwitchConfigStraight(t *testing.T) {

	buffer := []byte("{\"topic\": \"stringA\"}")
	var sc SwitchConfig

	err := json.Unmarshal(buffer, &sc)

	if err != nil {
		t.Fatal(err)
	}

	t.Logf("switch config/unmarshal: %v", sc)
}

func TestUnmarshalReadSwitchConfigInverted(t *testing.T) {

	buffer := []byte("{\"topic\": \"stringA\",\"inverted\":true}")
	var sc SwitchConfig

	err := json.Unmarshal(buffer, &sc)

	if err != nil {
		t.Fatal(err)
	}

	t.Logf("switch config/unmarshal: %v", sc)
}

func TestMarshalSwitchMap(t *testing.T) {
	switchMap := make(ConfigSwitchMap)
	switchMap["0"] = SwitchConfig{"Topic0", false}
	switchMap["1"] = SwitchConfig{"Topic1", true}

	buffer, err := json.Marshal(switchMap)

	if err != nil {
		t.Fatal(err)
	}

	t.Logf("switch map/marshal: %s", buffer)
}

func TestUnmarshalSwitchMapStraight(t *testing.T) {

	buffer := []byte("{\"0\":{\"topic\":\"Topic0\"},\"1\":{\"topic\":\"Topic1\",\"inverted\":true}}")
	switchMap := make(ConfigSwitchMap)

	err := json.Unmarshal(buffer, &switchMap)

	if err != nil {
		t.Fatal(err)
	}

	t.Logf("switch map/unmarshal/A: %v", switchMap)
}

func TestUnmarshalSwitchMapSwapped(t *testing.T) {

	buffer := []byte("{\"0\":{\"topic\":\"Topic0\"},\"1\":{\"inverted\":true,\"topic\":\"Topic1\"}}")
	switchMap := make(ConfigSwitchMap)

	err := json.Unmarshal(buffer, &switchMap)

	if err != nil {
		t.Fatal(err)
	}

	t.Logf("switch map/unmarshal/B: %v", switchMap)
}

func TestUnmarshalConfig(t *testing.T) {

	buffer := []byte("{\"entry\":{\"0\":{\"topic\":\"Topic0\"},\"1\":{\"topic\":\"Topic1\",\"inverted\":true}}}")
	config := make(map[string]ConfigSwitchMap)

	err := json.Unmarshal(buffer, &config)

	if err != nil {
		t.Fatal(err)
	}

	t.Logf("config/unmarshal: %v", config)
}
