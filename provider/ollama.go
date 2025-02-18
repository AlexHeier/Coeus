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
	"time"
)

const OLLAMA_GENERATE_SUFFIX = "/api/chat"

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
		Name        string `json:"name"`
		Description string `json:"description"`
		Parameters  any    `json:"parameters"`
	}
}

type ollamaMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ollamaFunctionCall struct {
	Function struct {
		Name      string                   `json:"name"`
		Arguments []map[string]interface{} `json:"arguments"`
	} `json:"function"`
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
	reqData["messages"] = []ollamaMessage{
		{Role: "system", Content: request.Systemprompt},
		{Role: "user", Content: request.Userprompt},
	}
	reqData["stream"] = config.Stream
	reqData["tools"] = ollamaToolsWrapper()

	jData, err := ollamaNetworkSender(reqData, url)
	if err != nil {
		return ResponseStruct{}, err
	}

	print(jData)

	_, ok := jData["message"].(map[string]interface{})["tool_calls"]
	if ok {
		calls := 
		for _, t := range calls {
			fmt.Printf("\n\n%v\n\n", t)
			time.Sleep(2 * time.Second)
			tool, err := tool.Find(t.Function.Name)
			if err != nil {
				continue
			}

			var args []interface{}
			args = append(args, t.Function.Arguments)

			toolResponse, err := tool.Run(args...)
			if err != nil {
				continue
			}

			reqData["messages"] = append(reqData["messages"].([]ollamaMessage),
				ollamaMessage{Role: "Tools", Content: string(toolResponse)})
		}
		jData, err = ollamaNetworkSender(reqData, url)
		if err != nil {
			return ResponseStruct{}, err
		}

	}

	return ResponseStruct{
		Response: jData["message"].(map[string]interface{})["content"].(string),
	}, nil
}

func ollamaToolsWrapper() []ollamaTool {
	var ollamaTools []ollamaTool

	for _, t := range tool.Tools {
		temp := ollamaTool{}
		temp.Type = "function"
		temp.Function.Name = t.Name
		temp.Function.Description = t.Desc
		temp.Function.Parameters = t.Params
		ollamaTools = append(ollamaTools, temp)
	}

	return ollamaTools
}

// Does Ollamas network handling
func ollamaNetworkSender(reqData map[string]interface{}, url string) (map[string]interface{}, error) {
	data, err := json.Marshal(reqData)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	resData, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var jData map[string]interface{}
	if err := json.Unmarshal(resData, &jData); err != nil {
		return nil, err
	}

	fmt.Println(jData)

	return jData, nil
}
