package main

import (
	metrics "github.com/toxyl/metric-nexus"
)

func main() {
	server := metrics.NewServer(config.Host, config.Port, config.StateFile)
	for _, k := range config.APIKeys {
		server.AddAPIKey(k)
	}
	panic(server.Start(config.KeyFile, config.CertFile))
}
