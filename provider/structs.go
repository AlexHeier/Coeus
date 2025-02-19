package provider

/*	  ____  _ _
	 / __ \| | |
	| |  | | | | __ _ _ __ ___   __ _
	| |  | | | |/ _` | '_ ` _ \ / _` |
	| |__| | | | (_| | | | | | | (_| |
 	 \____/|_|_|\__,_|_| |_| |_|\__,_|
*/

type OllamaStruct struct {
	HttpProtocol string
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
type OpenAIStruct struct {
	Model   string
	API_KEY string
}

/*
     	/\
       /  \    _____   _ _ __ ___
   	  / /\ \  |_  / | | | '__/ _ \
  	 / ____ \  / /| |_| | | |  __/
 	/_/    \_\/___|\__,_|_|  \___|
*/

//Azure structs here
