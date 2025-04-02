# Coeus - LLM Library

Coeus is a Golang library for sculpting exsisting LLM to work as the user wants. Coeus manages conversations, tool calls, memory, RAG and more.

## Usage

The most bearbone setup of Coeus is the following:

```golang
err := coeus.Ollama("10.0.0.1", "11434", "llama3.1:8b")
if err != nil {
	log.Fatal(err.Error())
}

conv := coeus.BeginConversation()

resp, err := conv.Prompt("Hello Ollama!")
if err != nil {
	log.Fatal(err.Error())
}

fmt.Println(resp.Response)
```
### Providers
Coeus supports also OpenAI and Azure as LLM providers:
```golang
err := coeus.OpenAI("gpt-4o-mini", "OPENAI_API_KEY")

err := coeus.Azure("AZURE_ENDPOINT", "AZURE_API_KEY", Temperature float64, MaxTokens int)

err := coeus.Ollama("10.0.0.1", "11434", "llama3.1:8b")
```

### System Prompt
If you want to send a custom systemprompt to the LLM to help shape the LLM to your needs. Use the following command:

```golang
coeus.SetSystemPrompt("You are an AI assistant")
```

### LLM Tools
If you want to LLM to be able to call functions that run on your computer then you can define the function as a tool.
```golang
coeus.NewTool("TOOL NAME", "TOOL DESCRIPTION", ToolFunc)
```
Example:
```golang
coeus.NewTool("MULTIPLY", "Multiply takes in two ints a and b and returns the resoult of a times b", Multiply)

func Multiply(a, b int) int (
    return a * b
)
```

### LLM Memory
The defoult memory version is MemoryAllMessage(). MemoryAllMessage() uses the intire conversation as memory for the conversations.

To change the memory version:
```golang
coeus.MemoryVersion(MemoryFunc)
```

There are three premade memory versions avalible:
```golang
coeus.MemoryVersion(MemoryAllMessage) // Default


// MemoryLastMessage is a function that will use the last int x messages as memory.
// @extra param: The number of last user messages to use as memory.
coeus.MemoryVersion(MemoryLastMessage, MessageCount int)

//MemoryTime is a function that will use the messages within the last int x minutes as memory.
//@extra param: The number of last minutes to use as memory.
coeus.MemoryVersion(MemoryTime, MessageAge int)
```

### Retrieval-Augmented Generation (RAG)
RAG is a way to ground a LLM by giving it access to predefined information. Coeus has its own custom made RAG model that uses a simple word2vec model made with skip-gram. This model only supports english. The RAG training can be found [here](https://github.com/AlexHeier/word2vector) 

To enable RAG use the following command:
```golang
err := coeus.EnableRAG(host, dbname, user, password string, port int)
if err != nil {
	log.Fatal(err.Error())
}
```
The current solution uses a PostgrSQL database to store the vectors in. If time will be changed to getting the vectors for external source and managing them itself.

To configure the RAG model you use:
```golang
err := coeus.RAGConfig(context, chunkSize int, overlapRatio, multiplier float64, folder string)
if err != nil {
	log.Fatal(err.Error())
}
```
**Params:**
- context: The number of closest results to use as context for the model. Default is 2.
- chunkSize: The size of the chunks to split the text into. Default is 300.
- overlapRatio: The ratio of overlap between chunks. Default is 0.25.
- multiplier: The multiplier for the vector scaling. Default is 2.
- folder: The folder where the RAG files are stored. Default is "./RAG".

## Installation

Coeus is currently only avalible to users who has access to the reposetory by running the following commands within a terminal within the program:

```powershell
git remote set-url origin git@github.com:AlexHeier/Coeus.git
go env -w GOPRIVATE=github.com/AlexHeier/Coeus
go get github.com/AlexHeier/Coeus@latest
```

However, when Coeus goes public, it will be accessable running the following in a terminal within your program:

```powershell
go get github.com/AlexHeier/Coeus
```

## Test

To run Coeus's test. Run the following command in a terminal at the root of Coeus.

```powershell
go test -v ./...
```

