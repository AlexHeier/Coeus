package llm

import (
	"Coeus/llm/tool"
	"Coeus/provider"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"
)

// Struct for containing all the conversations
var ConvAll struct {
	Mutex         sync.Mutex
	Conversations []*Conversation
}

// Struct for a single conversation.
type Conversation struct {
	Mutex      sync.Mutex
	MainPrompt string
	ToolsResp  []interface{}
	History    []string
	Summary    string
	UserPrompt string
	LastActive time.Time
}

// Appends a prompt and section to the history within the conversation
func (c *Conversation) AppendHistory(user, llm string) {
	newHistory := "[USER:] " + user + "[LLM:] " + llm
	c.History = append(c.History, newHistory)
}

func (c *Conversation) Prompt(userPrompt string) (provider.ResponseStruct, error) {
	var toolDesc string
	var response provider.ResponseStruct
	var err error
	//var toolResponse []interface{}
	var toolUsed bool = false

	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	c.LastActive = time.Now()

	if len(tool.Tools) > 0 {
		for _, t := range tool.Tools {
			toolDesc += t.Name + ": " + t.Desc + "\n"
		}
	}

	for {
		toolUsed = false

		// Memory(append([]interface{}{c}, MemArgs...)...) sends the conversation and the arguments to the memory function if the user defined some.
		/* Build the prompt (Structure should look like this)
		- Systemprompt
		- Tools Description
		- Tool results
		- Prior History (According to memory module)
		- New userprompt
		*/

		//systemprompt := "[SYSTEMPROMPT]\n" + c.MainPrompt + tool.GetToolsDescription()

		// Adds the tool responses if there exists any.
		//if len(toolResponse) > 0 {
		//	systemprompt += "About tool response: contains the results from previous tool calls. Always use these before the history.\n[BEGIN TOOL RESPONSE]\n" + fmt.Sprintf("%v", toolResponse) + "\n[END TOOL RESPONSE]\n"
		//}

		// Adds the history created by a specified memory function
		//systemprompt += Memory(append([]interface{}{c}, MemArgs...)...) + "[SYSTEMPROMPT]\n\n"

		//fmt.Println(systemprompt + userPrompt)
		//fmt.Printf("Systemprompt tokens: %d\n", len(systemprompt)/4)
		//fmt.Printf("Userprompt tokens: %d\n", len(userPrompt)/4)
		//fmt.Printf("Total tokens: %d\n", (len(systemprompt)+len(userPrompt))/4)

		response, err = provider.Send(c.BuildPrompt(userPrompt))
		if err != nil {
			return response, err
		}
		c.AppendHistory(userPrompt, response.Response)

		splitString := strings.Split(response.Response, " ")

		// Check for if the response contains a summary and extract it
		for i, w := range splitString {
			if strings.Contains(w, "SUMMARY") {
				c.Summary = strings.Join(splitString[i+1:], " ")
				response.Response = strings.Join(splitString[:i], " ")
				break
			}
		}

		// TO-DO: FÅ LLM TIL Å IKKE SENDE NAVN PÅ KOMMANDOER I FULL CAPS DER DEN IKKE SKAL BRUKE DEN

		for _, t := range tool.Tools {
			// Gets the index of the first letter in the tool name
			index := strings.Index(response.Response, t.Name)
			if index >= 0 {

				resArray := strings.Split(response.Response[index:], " ")

				// LLM might include newline chars which can give weird arguments if not handled
				for arrIndex, arrString := range resArray {
					newLineIndex := strings.Index(arrString, "\n")
					if newLineIndex >= 0 {
						resArray[arrIndex] = arrString[:newLineIndex]
						resArray = resArray[:arrIndex]
						break
					}
				}

				ft := reflect.ValueOf(t.Function)
				argCount := ft.Type().NumIn()

				args := resArray[1 : argCount+1]

				fmt.Printf("Arguments: %s\n", args)

				callArgs := make([]interface{}, argCount)
				for i := 0; i < argCount; i++ {
					callArgs[i] = args[i]
				}
				tr, err := t.Run(callArgs[0:]...)
				if err != nil {
					return response, err
				}

				c.ToolsResp = append(c.ToolsResp, fmt.Sprintf("%v %v = %v", t.Name, args, tr))

				response, err = provider.Send(c.BuildPrompt("[SYSTEM:] TOOL RESPONSE UPDATED. USE THIS TO ANSWER THE USERS PREVIOUS PROMPT"))
				if err != nil {
					return response, err
				}
			}
		}
		if !toolUsed {
			fmt.Println("Breaking")
			break
		}
		fmt.Println("What")
	}

	fmt.Println(c.BuildPrompt(""))

	return response, err
}

func (c *Conversation) DumpConversation() string {
	var temp string
	//c.Mutex.Lock()
	//defer c.Mutex.Unlock()

	for _, h := range c.History {
		temp += h
	}

	return temp
}

func BeginConversation() *Conversation {
	ConvAll.Mutex.Lock()
	defer ConvAll.Mutex.Unlock()
	newCon := Conversation{
		MainPrompt: Persona,
	}

	ConvAll.Conversations = append(ConvAll.Conversations, &newCon)
	return ConvAll.Conversations[len(ConvAll.Conversations)-1]
}

func (c *Conversation) BuildPrompt(prompt string) string {
	systemprompt := "[SYSTEMPROMPT]\n" + c.MainPrompt + tool.GetToolsDescription()

	if len(c.ToolsResp) > 0 {
		systemprompt += "About tool response: contains the results from previous tool calls. Always use these before the history.\n[BEGIN TOOL RESPONSE]\n" + fmt.Sprintf("%v", c.ToolsResp) + "\n[END TOOL RESPONSE]\n"
	}

	systemprompt += Memory(append([]interface{}{c}, MemArgs...)...) + "[SYSTEMPROMPT]\n\n"

	return systemprompt + prompt
}
