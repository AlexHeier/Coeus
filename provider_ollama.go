package coeus

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"strconv"
)

// Sufix for the Ollama API endpoint
const ollamaSuffix = "/api/chat"

/*
Ollama is a function that sets the provider to Ollama.

@param ip: the IP address of the Ollama server

@param port: the port of the Ollama server

@param model: the model to use

@return An error if the IP address, port or model is invalid
*/
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

	Provider = ollamaStruct{
		HTTPProtocol: "http",
		ServerIP:     ip,
		Port:         port,
		Model:        model,
		Stream:       false,
	}

	return nil
}

/*
sendOllama is a function that setups and sends a request to Ollama.

@param request: the request to send

@return A response and an error if the request fails
*/
func sendOllama(request RequestStruct) (ResponseStruct, error) {

	config := Provider.(ollamaStruct)

	url := config.HTTPProtocol + "://" + config.ServerIP + ":" + config.Port + ollamaSuffix

	reqData := make(map[string]interface{})

	reqData["model"] = config.Model
	reqData["messages"] = []ollamaRole{
		{Role: "system", Content: request.Systemprompt},
		{Role: "user", Content: request.Userprompt},
	}
	reqData["stream"] = config.Stream
	reqData["tools"] = ollamaToolsWrapper()

	jData, err := ollamaNetworkSender(reqData, url)
	if err != nil {
		return ResponseStruct{}, err
	}

	if len(jData.Message.ToolCalls) > 0 {
		calls := jData.Message.ToolCalls
		for _, t := range calls {
			tool, err := FindTool(t.Function.Name)
			if err != nil {
				continue
			}

			var args []interface{}
			for _, a := range t.Function.Arguments {
				args = append(args, a)
			}

			toolResponse, err := tool.RunTool(args...)
			if err != nil {
				continue
			}

			reqData["messages"] = append(reqData["messages"].([]ollamaRole),
				ollamaRole{Role: "tool", Content: toolResponse})

			*request.History = append(*request.History, HistoryStruct{
				Role:      "assistant",
				ToolCalls: []ToolCall{t},
			})

			*request.History = append(*request.History, HistoryStruct{
				Role:       "tool",
				Content:    toolResponse,
				ToolCallID: t.ID,
			})
		}

		jData, err = ollamaNetworkSender(reqData, url)
		if err != nil {
			return ResponseStruct{}, err
		}
	}

	return ResponseStruct{
		Response:           jData.Message.Content,
		TotalLoadDuration:  float64(jData.TotalDuration),
		Eval_count:         float64(jData.EvalCount),
		PromptEvalCount:    float64(jData.PromptEvalCount),
		PromptEvalDuration: float64(jData.PromptEvalDuration),
	}, nil
}

/*
ollaToolsWrapper is a function that wraps the tools in the request to Ollama.

@return A list of tools in the request to Ollama
*/
func ollamaToolsWrapper() []ollamaTool {

	var ollamaTools []ollamaTool

	for _, t := range Tools {
		temp := ollamaTool{}
		temp.Type = "function"
		temp.Function.Name = t.Name
		temp.Function.Description = t.Desc
		temp.Function.Parameters = t.Params
		ollamaTools = append(ollamaTools, temp)
	}

	return ollamaTools
}

/*
ollaNetworkSender is a function that sends a request to Ollama.

@param reqData: the data to send

@param url: the URL to send the data to

@return A response and an error if the request fails
*/
func ollamaNetworkSender(reqData map[string]interface{}, url string) (ollamaResponse, error) {
	data, err := json.Marshal(reqData)
	if err != nil {
		return ollamaResponse{}, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return ollamaResponse{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return ollamaResponse{}, err
	}
	defer res.Body.Close()

	resData, err := io.ReadAll(res.Body)
	if err != nil {
		return ollamaResponse{}, err
	}

	var jData ollamaResponse
	if err := json.Unmarshal(resData, &jData); err != nil {
		return ollamaResponse{}, err
	}

	return jData, nil
}
