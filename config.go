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

// GetLoopSleep returns main loop iteration delay in miliseconds
func (c *Config) GetLoopSleep() string {
	return c.LoopSleep
}

// GetUnixSocket returns unix_socket value which is a file used to communiate
// with running ntree
func (c *Config) GetUnixSocket() string {
	return c.UnixSocket
}

// SetFromJSON sets config instance values from JSON
func (c *Config) SetFromJSON(b []byte) {
	err := json.Unmarshal(b, c)
	if err != nil {
		log.Fatal("Error setting config from JSON:", err.Error())
	}
}
