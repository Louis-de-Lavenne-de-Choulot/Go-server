package store

type Player struct {
	Name string
	Wins int
}

type PlayerStorage interface {
	GetPlayer(name string) (Player, error)
	IncWins(name string) (int, error)
	RemovePlayer(name string) error
	AddPlayer(name string) error
	GetAllPlayers() ([]Player, error)
}
