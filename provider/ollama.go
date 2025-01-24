package provider

import (
	"errors"
	"net"
	"strconv"
)

type Ollama struct {
	Provider     string
	HttpProtocol string
	ServerIP     string
	Port         string
	Model        string
	Stream       bool
}

func NewOllama(ip, port, model string) (Ollama, error) {
	// Validate IP address
	if net.ParseIP(ip) == nil {
		return Ollama{}, errors.New("invalid IP address")
	}

	// Validate port
	if _, err := strconv.Atoi(port); err != nil {
		return Ollama{}, errors.New("invalid port")
	}

	// Validate model (example: non-empty string)
	if model == "" {
		return Ollama{}, errors.New("model cannot be empty")
	}

	return Ollama{
		Provider:     "Ollama",
		HttpProtocol: "http",
		ServerIP:     ip,
		Port:         port,
		Model:        model,
		Stream:       false,
	}, nil
}
