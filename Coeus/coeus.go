package coeus

import (
	"fmt"
	"net/http"
)

func (c *Coeus) Init(Config SrvConfig, LLMEndpoint LLMConfig) {
	c.Config.HttpProtocol = Config.HttpProtocol
	c.Config.IPaddress = Config.IPaddress
	c.Config.Router = Config.Router
	c.Config.Port = Config.Port

	c.LLMEndpoint.provider = LLMEndpoint.provider
	c.LLMEndpoint.HttpProtocol = LLMEndpoint.HttpProtocol
	c.LLMEndpoint.SystemPrompt = LLMEndpoint.SystemPrompt
	c.LLMEndpoint.ServerIP = LLMEndpoint.ServerIP
	c.LLMEndpoint.Stream = LLMEndpoint.Stream
	c.LLMEndpoint.Model = LLMEndpoint.Model
	c.LLMEndpoint.Port = LLMEndpoint.Port
}

func (c *Coeus) Start() error {
	fmt.Printf("Running server at %s:%s:%s\n", c.Config.HttpProtocol, c.Config.IPaddress, c.Config.Port)
	err := http.ListenAndServe(":"+c.Config.Port, c.Config.Router)
	if err != nil {
		return err
	}
	return nil
}
