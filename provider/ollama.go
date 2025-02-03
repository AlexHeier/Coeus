package provider

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"strconv"
)

const OLLAMA_GENERATE_SUFFIX = "/api/generate"

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

func SendOllama(prompt string) (ResponseStruct, error) {

	config := Provider.(OllamaStruct)

	url := "http://" + config.ServerIP + ":" + config.Port + OLLAMA_GENERATE_SUFFIX

	reqData := make(map[string]interface{})

	reqData["model"] = config.Model
	reqData["stream"] = config.Stream
	reqData["prompt"] = prompt

	data := new(bytes.Buffer)

	json.NewEncoder(data).Encode(reqData)

	req, err := http.NewRequest(http.MethodPost, url, data)
	if err != nil {
		return ResponseStruct{}, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return ResponseStruct{}, err
	}
	defer res.Body.Close()

	resData, err := io.ReadAll(res.Body)
	if err != nil {
		return ResponseStruct{}, err
	}

	jData := make(map[string]interface{})

	err = json.Unmarshal(resData, &jData)
	if err != nil {
		return ResponseStruct{}, err
	}

	return ResponseStruct{
		Response:             jData["response"].(string),
		LoadDuration:         jData["load_duration"].(float64),
		eval_count:           jData["eval_count"].(float64),
		prompt_eval_count:    jData["prompt_eval_count"].(float64),
		prompt_eval_duration: jData["prompt_eval_duration"].(float64),
	}, nil
}
