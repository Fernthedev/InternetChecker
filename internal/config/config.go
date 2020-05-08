package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"runtime"
)

type PingAddressSetting struct {
	PingAddress string `json:"ping_address"`
}

type NetInterfaceSettings struct {
	PingAddressSetting[] PingAddressSetting `json:"ping_addresses"`
}

type Config struct {
	CheckNetworkInterfaces map[string]NetInterfaceSettings  `json:"check_network_interfaces"`

	GlobalPingAddresses[] PingAddressSetting `json:"global_ping_addresses"`
}

func SaveDefaultConfig(filename string) (Config, error) {


	var defaultEth string

	if runtime.GOOS == "windows" {
		defaultEth = "Ethernet"
	} else {
		defaultEth = "eth0"
	}

	var config = Config{
		map[string]NetInterfaceSettings{
			defaultEth: {
				[]PingAddressSetting{{"1.1.1.1"}},
			},
		},
		[]PingAddressSetting{{"8.8.8.8"}},
	}


	jsonMarsh, err := json.Marshal(config)

	if err != nil {
		return config, err
	}

	err = ioutil.WriteFile(filename, jsonMarsh, 0644)

	if err != nil {
		return config, err
	}

	return config, err
}

func LoadConfiguration(filename string) (Config, error) {
	var config Config
	configFile, err := os.Open(filename)

	if err != nil {
		return config, err
	}

	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)

	return config, err
}