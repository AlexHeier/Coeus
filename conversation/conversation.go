package conversation

func (c *Struct) Setup(llm interface{}, memoryFunc func(args ...interface{}) string) {
	c.llm = llm
	c.Memory = memoryFunc
}
