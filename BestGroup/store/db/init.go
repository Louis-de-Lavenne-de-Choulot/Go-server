package playerdb

import (
	"database/sql"

	"github.com/tuxago/go/BestGroup/store"
)

func InitDB(db *sql.DB) error {
	sqlStmt := `
	create table if not exists players (id integer not null primary key, name text, wins integer);
	delete from players;
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		panic(err)
	}

	for _, p := range initialPlayersToImport {
		err = AddPlayer(db, p.Name)
		if err != nil {
			return err
		}
	}
	return nil
}

var initialPlayersToImport = []store.Player{
	{
		Name: "John",
	},
	{
		Name: "Jane",
	},
	{
		Name: "Pepper",
	},
	{
		Name: "Salt",
	},
	{
		Name: "Louis",
	},
	{
		Name: "Maxime",
	},
	{
		Name: "Malo",
	},
	{
		Name: "Arthur",
	},
}
