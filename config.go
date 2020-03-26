package main

import (
	"encoding/json"
	"log"
)

type Config struct {
	Version    string `json:"version"`
	Id         string `json:"id"`
	LoopSleep  string `json:"loop_sleep"`
	UnixSocket string `json:"unix_socket"`
}

func (c *Config) GetLoopSleep() string {
	return c.LoopSleep
}

func (c *Config) GetUnixSocket() string {
	return c.UnixSocket
}

func (c *Config) SetFromJSON(b []byte) {
	err := json.Unmarshal(b, c)
	if err != nil {
		log.Fatal("Error setting config from JSON:", err.Error())
	}
}
