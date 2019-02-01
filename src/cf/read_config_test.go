package cf

import (
	"encoding/json"
	"testing"
)

func TestMarshalReadSwitchConfig(t *testing.T) {

	sc := SwitchConfig{ "Topic0", true}
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
