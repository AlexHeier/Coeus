package coeus

import (
	"time"

	"github.com/gorilla/mux"
)

type Coeus struct {
	Config      SrvConfig
	LLMEndpoint LLMConfig
	StartupTime time.Time
}

type SrvConfig struct {
	Router       *mux.Router
	HttpProtocol string
	IPaddress    string
	Port         string
}

// Used connecting a Coeus server instance to a Ollama server
type LLMConfig struct {
	provider     string // Specifies which provider to use. Etc ollama, Azure, OpenAI, AWS. Currently Only Ollama works.
	HttpProtocol string // Should be HTTP or HTTPS
	ServerIP     string // LLM server instance IP
	Port         string // LLM server port
	Model        string // Configures which LLM model to be used
	SystemPrompt string // Designates the base systempromt to be used by the LLM
	Stream       bool   // Decides if the response should be a single reply or multiple streamed ones.
}
