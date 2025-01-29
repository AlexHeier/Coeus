package llm

type Tool struct {
	desc     string
	function interface{}
}

var Tools []Tool

func NewTool(d string, f interface{}) {

	var newTool Tool
	newTool.desc = d
	newTool.function = f

	Tools = append(Tools, newTool)

}
