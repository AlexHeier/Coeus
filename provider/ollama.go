package provider

import (
	"Coeus/llm/tool"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"strconv"
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

type ollamaRole struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ollamaToolCall struct {
	Function struct {
		Arguments map[string]interface{} `json:"arguments"`
		Name      string                 `json:"name"`
	} `json:"function"`
}

type ollamaMessage struct {
	Content   string           `json:"content"`
	Role      string           `json:"role"`
	ToolCalls []ollamaToolCall `json:"tool_calls"`
}

type ollamaResponse struct {
	CreatedAt          string        `json:"created_at"`
	Done               bool          `json:"done"`
	DoneReason         string        `json:"done_reason"`
	EvalCount          int           `json:"eval_count"`
	EvalDuration       int64         `json:"eval_duration"`
	LoadDuration       int64         `json:"load_duration"`
	Message            ollamaMessage `json:"message"`
	Model              string        `json:"model"`
	PromptEvalCount    int           `json:"prompt_eval_count"`
	PromptEvalDuration int64         `json:"prompt_eval_duration"`
	TotalDuration      int64         `json:"total_duration"`
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
			tool, err := tool.Find(t.Function.Name)
			if err != nil {
				continue
			}

			var args []interface{}
			for _, a := range t.Function.Arguments {
				args = append(args, a)
			}

			toolResponse, err := tool.Run(args...)
			if err != nil {
				continue
			}

			reqData["messages"] = append(reqData["messages"].([]ollamaRole),
				ollamaRole{Role: "tool", Content: string(toolResponse)})
		}

		// Printer ut reqData["messages"] til terminalen
		/**
		messagesJSON, err := json.MarshalIndent(reqData["messages"], "", "  ")
		if err != nil {
			fmt.Println("Error marshalling messages:", err)
		} else {
			fmt.Println("Updated messages:", string(messagesJSON))
		}
		*/

		jData, err = ollamaNetworkSender(reqData, url)
		if err != nil {
			return ResponseStruct{}, err
		}

	}

	return ResponseStruct{
		Response: jData.Message.Content,
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
