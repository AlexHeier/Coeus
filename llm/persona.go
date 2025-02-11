package llm

var Persona string

func SetPersona(s string) {
	Persona = "[MAIN BEGIN]\n" + s + "\n[MAIN END]\n"
}
