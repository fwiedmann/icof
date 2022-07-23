package gocrazy_permanent_data

import (
	"encoding/json"
	"io/ioutil"
)

// ConfigLocation is used to load the config required for icof to start.
// gokrazy stores the permanent data under the /perm directory
var ConfigLocation = "/perm/icof/start-config.json"

// New loads the startup config from disk
func New() (StartUpConfig, error) {
	content, err := ioutil.ReadFile(ConfigLocation)
	if err != nil {
		return StartUpConfig{}, err
	}

	var c StartUpConfig
	if err := json.Unmarshal(content, &c); err != nil {
		return StartUpConfig{}, err
	}

	return c, nil
}

// StartUpConfig holds all required properties for icof to start
type StartUpConfig struct {
	EmailClientConfig EmailClientConfig `json:"email_config"`
}

type EmailClientConfig struct {
	Host             string `json:"host"`
	Port             int    `json:"port"`
	Username         string `json:"username"`
	Password         string `json:"password"`
	FromEmailAddress string `json:"from_email_address"`
}
