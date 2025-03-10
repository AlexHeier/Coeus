package coeus

/*	  ____  _ _
	 / __ \| | |
	| |  | | | | __ _ _ __ ___   __ _
	| |  | | | |/ _` | '_ ` _ \ / _` |
	| |__| | | | (_| | | | | | | (_| |
 	 \____/|_|_|\__,_|_| |_| |_|\__,_|
*/

type ollamaStruct struct {
	HTTPProtocol string
	ServerIP     string
	Port         string
	Model        string
	Stream       bool
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

type ollamaTool struct {
	Type     string `json:"type"`
	Function struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Parameters  any    `json:"parameters"`
	}
}

type ollamaMessage struct {
	Content   string           `json:"content"`
	Role      string           `json:"role"`
	ToolCalls []ollamaToolCall `json:"tool_calls"`
}

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
type openAIStruct struct {
	Model  string
	ApiKey string
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
type azureResponse struct {
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
