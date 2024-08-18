package structs

type Intent struct {
	Tag       string   `json:"tag"`
	Patterns  []string `json:"patterns"`
	Responses []string `json:"responses"`
}

type Intents struct {
	Intents []Intent `json:"intents"`
}
