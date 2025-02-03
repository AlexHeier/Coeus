package provider

type AzureStruct struct {
	Provider     string
	HttpProtocol string
	ServerIP     string
	Port         string
	Model        string
	Stream       bool
}

/*
func NewAzure(ip, port, model string) (Azure, error) {
    // Validate IP address
    if net.ParseIP(ip) == nil {
        return Azure{}, errors.New("invalid IP address")
    }

    // Validate port
    if _, err := strconv.Atoi(port); err != nil {
        return Azure{}, errors.New("invalid port")
    }

    // Validate model (example: non-empty string)
    if model == "" {
        return Azure{}, errors.New("model cannot be empty")
    }

    return Azure{
        Provider:     "Ollama",
        HttpProtocol: "http",
        ServerIP:     ip,
        Port:         port,
        Model:        model,
        Stream:       false,
    }, nil
}
*/

func SendAzure(prompt string) (ResponseStruct, error) {
	return ResponseStruct{}, nil
}
