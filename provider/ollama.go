package provider

import (
	"Coeus/llm/tool"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
)

const OLLAMA_GENERATE_SUFFIX = "/api/generate"

type OllamaStruct struct {
	HttpProtocol string
	ServerIP     string
	Port         string
	Model        string
	Stream       bool
}

type ollamaTool struct {
	Type     string `json:"type"`
	Function struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Parameters  any      `json:"parameters"`
		Required    []string `json:"required"`
	}
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
		HttpProtocol: "http",
		ServerIP:     ip,
		Port:         port,
		Model:        model,
		Stream:       false,
	}

	return nil
}

func SendOllama(request RequestStruct) (ResponseStruct, error) {

	config := Provider.(OllamaStruct)

	url := "http://" + config.ServerIP + ":" + config.Port + OLLAMA_GENERATE_SUFFIX

	reqData := make(map[string]interface{})

	reqData["model"] = config.Model
	reqData["messages"] = []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}{{Role: "user", Content: request.Userprompt}, {Role: "system", Content: request.Systemprompt}}
	reqData["stream"] = config.Stream
	reqData["tools"] = ollamaToolsWrapper()

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

	fmt.Printf("Response: %v\n", jData)

	return ResponseStruct{
		Response:             jData["response"].(string),
		LoadDuration:         jData["load_duration"].(float64),
		eval_count:           jData["eval_count"].(float64),
		prompt_eval_count:    jData["prompt_eval_count"].(float64),
		prompt_eval_duration: jData["prompt_eval_duration"].(float64),
	}, nil
}

func ollamaToolsWrapper() string {
	var ollamaTools []ollamaTool

	for _, t := range tool.Tools {
		temp := ollamaTool{}
		temp.Type = "function"
		temp.Function.Name = t.Name
		temp.Function.Description = t.Desc
		temp.Function.Parameters = t.Params
		ollamaTools = append(ollamaTools, temp)
	}
	return fmt.Sprintf("%v", ollamaTools)
}
