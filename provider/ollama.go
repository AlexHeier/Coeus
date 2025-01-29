package provider

import (
	"errors"
	"net"
	"strconv"
)

type OllamaStruct struct {
	Provider     string
	HttpProtocol string
	ServerIP     string
	Port         string
	Model        string
	Stream       bool
}

func Ollama(ip, port, model string) error {
	// Validate IP address
	if net.ParseIP(ip) == nil {
		return errors.New("invalid IP address")
	}

	// Validate port
	if _, err := strconv.Atoi(port); err != nil {
		return errors.New("invalid port")
	}

	// Validate model (example: non-empty string)
	if model == "" {
		return errors.New("model cannot be empty")
	}

	Provider = OllamaStruct{
		Provider:     "Ollama",
		HttpProtocol: "http",
		ServerIP:     ip,
		Port:         port,
		Model:        model,
		Stream:       false,
	}

	return nil
}
