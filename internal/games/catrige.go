package games

type Game struct {
	Emulator    string `json:"emulator"`
	URL         string `json:"url"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type GameList struct {
	Games []Game `json:"games"`
}
