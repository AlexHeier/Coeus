# Coeus - LLM Library

Coeus is a Golang library designed for customizing existing Large Language Models (LLMs) to work according to user preferences. It provides tools for managing conversations, tool calls, memory, Retrieval-Augmented Generation (RAG), and more.

## Features
- **Multi-provider support**: Works with OpenAI, Azure, and Ollama.
- **System prompts**: Shape the LLM’s behavior using custom prompts.
- **LLM tools**: Extend the model’s capabilities with custom function calls.
- **Memory management**: Configurable conversation memory options.
- **Retrieval-Augmented Generation (RAG)**: Enhance responses with external knowledge.
- **Debug dashboard**: A web-based tool to test your setup and responses.

## Installation

Currently, Coeus is available only to users with repository access. To install:
```sh
# Configure Git to access the private repository
git remote set-url origin git@github.com:AlexHeier/Coeus.git

# Set GOPRIVATE to avoid authentication issues
go env -w GOPRIVATE=github.com/AlexHeier/Coeus

# Install Coeus
go get github.com/AlexHeier/Coeus@latest
```

When Coeus goes public, it will be installable with:
```sh
go get github.com/AlexHeier/Coeus
```

## Usage

The simplest way to use Coeus:
```go
err := coeus.Ollama("10.0.0.1", "11434", "llama3.1:8b")
if err != nil {
	log.Fatal(err)
}

conv := coeus.BeginConversation()

resp, err := conv.Prompt("Hello Ollama!")
if err != nil {
	log.Fatal(err)
}

fmt.Println(resp.Response)
```

### Supported LLM Providers
Coeus supports multiple LLM providers:
```go
err := coeus.OpenAI("gpt-4o-mini", "OPENAI_API_KEY")
err := coeus.Azure("AZURE_ENDPOINT", "AZURE_API_KEY", temperature float64, maxTokens int)
err := coeus.Ollama("10.0.0.1", "11434", "llama3.1:8b")
```

### System Prompts
To set a system prompt for shaping the LLM’s responses:
```go
coeus.SetSystemPrompt("You are an AI assistant")
```

### LLM Tools
Define functions that the LLM can call:
```go
coeus.NewTool("MULTIPLY", "Multiplies two integers", Multiply)

func Multiply(a, b int) int {
	return a * b
}
```

### Memory Management
The default memory mode is `MemoryAllMessage()`, which retains the entire conversation.
To change the memory mode:
```go
coeus.MemoryVersion(MemoryFunc)
```
Available memory modes:
```go
// Default: Uses all conversation messages.
coeus.MemoryVersion(MemoryAllMessage)

// Uses only the last X messages.
coeus.MemoryVersion(MemoryLastMessage, messageCount int)

// Uses messages within the last X minutes.
coeus.MemoryVersion(MemoryTime, messageAge int)
```

To implement custom memory management functions, ensure they follow to the following signature:

```go
func(c *Conversation) ([]HistoryStruct, error)
```

**Requirements:**
- Functions can accept any number of predefined input arguments found in the `memArgs` slice:
  
  ```go
  memArgs []interface{}
  ```

- The function should return a slice of `HistoryStruct` and an `error` value.

For reference, check how predefined memory functions are implemented in `coeus/llm_memory.go`.



### Retrieval-Augmented Generation (RAG)
Coeus includes a custom RAG model based on a skip-gram word2vec model (supports English only). The RAG training process is available [here](https://github.com/AlexHeier/word2vector).

To enable RAG:
```go
err := coeus.EnableRAG(host, dbname, user, password string, port int)
if err != nil {
	log.Fatal(err)
}
```
Configuration options:
```go
err := coeus.RAGConfig(context, chunkSize int, overlapRatio, multiplier float64, folder string)
```
**Parameters:**
- `context`: Number of closest results used as context (default: 2).
- `chunkSize`: Text chunk size (default: 300).
- `overlapRatio`: Overlapping ratio between chunks (default: 0.25).
- `multiplier`: Vector scaling multiplier (default: 2).
- `folder`: Storage folder for RAG files (default: "./RAG").

To test what the RAG would retrive given a user prompt:

```go
fmt.Println(coeus.GetRAG("Test prompt"))
```

### Debug Dashboard
To start a debug window as a local web interface:
```go
err := coeus.StartDashboard(port string)
if err != nil {
	log.Fatal(err)
}
```
Access it at `localhost:PORT`.

### Automatic Conversation Cleanup
To enable automatic conversation cleanup that will remove unactive conversations from memory:
```go
err := coeus.ConversationTTL(ttl int)
if err != nil {
	log.Fatal(err)
}
```
TTL is the number of minutes a conversation can be inactive before terminated.

## Running Tests
To run tests for Coeus:
```sh
go test -v ./...
```

