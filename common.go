package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type config struct {
	TunnelService  string `json:'tunnelService'`
	FrontService   string `json:'frontService'`
	BackendService string `json:'backendService'`
}

func loadConfig() (*config, error) {
	f, err := os.Open("config.json")
	checkError(err)
	defer f.Close()

	var cfg config
	err = json.NewDecoder(f).Decode(&cfg)
	return &cfg, err
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
