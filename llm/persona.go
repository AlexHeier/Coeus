package llm

var Persona string

func SetPersona(s string) {
	Persona = "[SYSTEMPROMPT BEGIN]\n" + s + "\n[SYSTEMPROMPT END]\n"
}
