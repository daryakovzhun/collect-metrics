package app

import "flag"

type AgentConfig struct {
	ServerAddress  string
	ReportInterval int
	PollInterval   int
}

type ServerConfig struct {
	Address string
}

func getAgentConfig() *AgentConfig {
	cfg := &AgentConfig{}

	flag.StringVar(&cfg.ServerAddress, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&cfg.ReportInterval, "r", 10, "frequency of sending metrics to the server")
	flag.IntVar(&cfg.PollInterval, "p", 2, "the frequency of polling metrics from the package")
	flag.Parse()

	return cfg
}

func getServerConfig() *ServerConfig {
	cfg := &ServerConfig{}

	flag.StringVar(&cfg.Address, "a", "localhost:8080", "address and port to run server")
	flag.Parse()

	return cfg
}
