package coeus

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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
func Ollama(ip, port, model string, temperature float64) error {
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
		Temperature:  temperature,
	}

	return nil
}

/*
sendOllama is a function that setups and sends a request to Ollama.

@param request: the request to send

@return A response and an error if the request fails
*/
func sendOllama(con *Conversation) (ResponseStruct, error) {

	config := Provider.(ollamaStruct)

	url := config.HTTPProtocol + "://" + config.ServerIP + ":" + config.Port + ollamaSuffix

	ollamaReq := ollamaRequest{
		Model:  config.Model,
		Stream: config.Stream,
		Tools:  ollamaToolsWrapper(),
		Options: struct {
			Temperature float64 "json:\"temperature,omitempty\""
			Seed        int     "json:\"seed,omitempty\""
		}{Temperature: config.Temperature},
	}

	history, err := memory(con)
	if err != nil {
		history = []HistoryStruct{{Role: "system", Content: sp}}
	}

	ollamaReq.Messages = append(ollamaReq.Messages, convertHistoryToOllama(history, con)...)

	jData, err := ollamaNetworkSender(ollamaReq, url)
	if err != nil {
		return ResponseStruct{}, err
	}

	if len(jData.Message.ToolCalls) > 0 {

		for {

			con.History = append(con.History, HistoryStruct{
				Role:      "assistant",
				ToolCalls: convertOllamaToolCallsToHistory(jData.Message.ToolCalls),
			})

			for _, t := range jData.Message.ToolCalls {
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
					fmt.Println(err.Error())
					continue
				}

				con.History = append(con.History, HistoryStruct{
					Role:       "tool",
					ToolCallID: t.ID,
					Content:    toolResponse,
				})

			}

			history, err := memory(con)
			if err != nil {
				history = []HistoryStruct{{Role: "system", Content: sp}}
			}

			ollamaReq.Messages = convertHistoryToOllama(history, con)

			jData, err = ollamaNetworkSender(ollamaReq, url)
			if err != nil {
				return ResponseStruct{}, err
			}

			if len(jData.Message.ToolCalls) == 0 {
				break
			}
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
func ollamaNetworkSender(reqData ollamaRequest, url string) (ollamaResponse, error) {
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

func convertOllamaToolCallsToHistory(t []ollamaToolCall) []ToolCall {
	var array []ToolCall
	for _, call := range t {
		data, err := json.Marshal(call.Function.Arguments)
		if err != nil {
			return nil
		}

		array = append(array, ToolCall{
			Index: call.Index,
			ID:    call.ID,
			Type:  call.Type,
			Function: struct {
				Name      string "json:\"name,omitempty\""
				Arguments string "json:\"arguments,omitempty\""
			}{Name: call.Function.Name,
				Arguments: string(data)},
		})
	}
	return array
}

func convertHistoryToolCallsToOllama(t []ToolCall) []ollamaToolCall {
	var array []ollamaToolCall

	for _, call := range t {

		var m = make(map[string]interface{})

		err := json.Unmarshal([]byte(call.Function.Arguments), &m)
		if err != nil {
			return nil
		}

		array = append(array, ollamaToolCall{
			Index: call.Index,
			ID:    call.ID,
			Type:  call.Type,
			Function: struct {
				Arguments map[string]interface{} "json:\"arguments\""
				Name      string                 "json:\"name\""
			}{Name: call.Function.Name,
				Arguments: m},
		})
	}
	return array
}

func convertHistoryToOllama(h []HistoryStruct, con *Conversation) []ollamaMessage {
	var array []ollamaMessage
	for _, his := range h {
		array = append(array, ollamaMessage{
			Content:    his.Content,
			Role:       his.Role,
			ToolCalls:  convertHistoryToolCallsToOllama(his.ToolCalls),
			ToolCallID: his.ToolCallID,
		})
	}

	if rag {
		array = append(array, ollamaMessage{
			Role:    "system",
			Content: "This is innformation from RAG: " + GetRAG(con.History[len(con.History)-1].Content),
		})
	}
	return array
}
