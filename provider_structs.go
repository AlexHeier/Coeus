package coeus

/*	  ____  _ _
	 / __ \| | |
	| |  | | | | __ _ _ __ ___   __ _
	| |  | | | |/ _` | '_ ` _ \ / _` |
	| |__| | | | (_| | | | | | | (_| |
 	 \____/|_|_|\__,_|_| |_| |_|\__,_|
*/

// Struct for Ollama definition
type ollamaStruct struct {
	HTTPProtocol string
	ServerIP     string
	Port         string
	Model        string
	Stream       bool
	Temperature  float64
}

// Struct used in sending requests to Ollama
type ollamaRole struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Struct used in defineing the tools in the request to Ollama
type ollamaTool struct {
	Type     string `json:"type"`
	Function struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Parameters  any    `json:"parameters"`
	}
}

// Struct used for tool calls in the response from Ollama
type ollamaToolCall struct {
	Index    *int   `json:"index,omitempty"`
	ID       string `json:"id,omitempty"`
	Type     string `json:"type,omitempty"`
	Function struct {
		Arguments map[string]interface{} `json:"arguments"`
		Name      string                 `json:"name"`
	} `json:"function"`
}

// Struct used in resonse from Ollama
type ollamaMessage struct {
	Content    string           `json:"content"`
	Role       string           `json:"role"`
	ToolCalls  []ollamaToolCall `json:"tool_calls"`
	ToolCallID string           `json:"tool_call_id"`
}

// Struct for containing the request to Ollama
type ollamaRequest struct {
	Model    string          `json:"model"`
	Messages []ollamaMessage `json:"messages"`
	Tools    []ollamaTool    `json:"tools"`
	Stream   bool            `json:"stream"`
	Options  struct {
		Temperature float64 `json:"temperature,omitempty"`
		Seed        int     `json:"seed,omitempty"`
	} `json:"options,omitempty"`
}

// Struct for containing the response from Ollama
type ollamaResponse struct {
	EvalCount          int           `json:"eval_count"`
	Message            ollamaMessage `json:"message"`
	PromptEvalCount    int           `json:"prompt_eval_count"`
	PromptEvalDuration int64         `json:"prompt_eval_duration"`
	TotalDuration      int64         `json:"total_duration"`
}

/*    ____                            _____
	 / __ \                     /\   |_   _|
	| |  | |_ __   ___ _ __    /  \    | |
	| |  | | '_ \ / _ \ '_ \  / /\ \   | |
	| |__| | |_) |  __/ | | |/ ____ \ _| |_
 	 \____/| .__/ \___|_| |_/_/    \_\_____|
           | |
           |_|
*/

// Struct for OpenAI definition
type openAIStruct struct {
	Model       string
	ApiKey      string
	Temperature float32
}

/*
     	/\
       /  \    _____   _ _ __ ___
   	  / /\ \  |_  / | | | '__/ _ \
  	 / ____ \  / /| |_| | | |  __/
 	/_/    \_\/___|\__,_|_|  \___|
*/

// Azure Configuration struct
type azureProviderStruct struct {
	Endpoint    string  // Azure API endpoint
	APIKey      string  // Azure API Key
	Temperature float64 // How free thinking the LLM should be. Lower equals more free. Can be between 0.1 and 1.0
	MaxTokens   int     // Max amount of tokens a response can use
}

// Struct used in sending requests to an Azure endpoint
type azureRequest struct {
	Messages    []azureMessage `json:"messages"`
	Tools       []azureTool    `json:"tools"`
	MaxTokens   int            `json:"max_tokens"`
	Temperature float64        `json:"temperature"`
}

// Struct for containing the response from Azure
type azureCompletionMessage struct {
	Choices []struct {
		ContentFilterResults map[string]interface{} `json:"content_filter_results"`
		FinishReason         string                 `json:"finish_reason"`
		LogProbs             string                 `json:"logprobs"`
		Message              azureMessage           `json:"message"`
	} `json:"choices"`
	Model               string `json:"model"`
	PromptFilterResults []struct {
		PromptIndex          int `json:"prompt_index"`
		ContentFilterResults map[string]struct {
			Filtered bool   `json:"filtered"`
			Severity string `json:"severity"`
		} `json:"content_filter_results"`
	} `json:"prompt_filter_results"`
	Usage struct {
		CompletionTokens int `json:"completion_tokens"`
		PromptTokens     int `json:"prompt_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

type azureMessage struct {
	Content    string     `json:"content"`
	Refusal    string     `json:"refusal,omitempty"`
	Role       string     `json:"role"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
	ToolCallID string     `json:"tool_call_id,omitempty"`
}

type azureTool struct {
	Type     string `json:"type"`
	Function struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Parameters  any      `json:"parameters"`
		Required    []string `json:"required"`
	} `json:"function"`
}
