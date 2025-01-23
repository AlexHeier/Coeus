package coeus

// Creates a new config for a Ollama Endpoint
func (c *Coeus) NewOllamaConfig(Protocol string, IP string, Port string, Model string, SysPrompt string, Streamed bool) LLMConfig {
	var conf LLMConfig

	conf.provider = "ollama"
	conf.HttpProtocol = Protocol
	conf.ServerIP = IP
	conf.Port = Port
	conf.Model = Model
	conf.SystemPrompt = SysPrompt
	conf.Stream = Streamed

	return conf
}
