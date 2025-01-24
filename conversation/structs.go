package conversation

type Struct struct {
	llm    interface{}
	Memory func(args ...interface{}) string
}
