package database

import (
	"database/sql"
)

type NFLEventsRepository struct {
	DB *sql.DB
}

type NFLEvent struct {
	GameId string
	Team1 string
	Team2 string
	HomeTeam string
	Winner string
}

func (r *NFLEventsRepository) CreateTable() error {
 _, err := r.DB.Exec(`CREATE TABLE IF NOT EXISTS nflevents (
	 gameId TEXT PRIMARY KEY,
	 team1 TEXT,
	 team2 TEXT,
	 homeTeam TEXT,
	 winner TEXT
 )`)

 _, err = r.DB.Exec("DELETE FROM nflevents")

 return err
}

func (r *NFLEventsRepository) Insert(row NFLEvent) error {
	_, err := r.DB.Exec(`INSERT INTO nflevents (
		gameId,
		team1,
		team2,
		homeTeam,
		winner) VALUES (?, ?, ?, ?, ?)`,
	row.GameId, row.Team1, row.Team2, row.HomeTeam, row.Winner)

	return err
}

func (r *NFLEventsRepository) GetAll() ([]NFLEvent, error) {
	rows, err := r.DB.Query("SELECT * FROM nflevents")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var nflEvents []NFLEvent

	for rows.Next() {
		var nflEvent NFLEvent
		err := rows.Scan(&nflEvent.GameId, &nflEvent.Team1, &nflEvent.Team2, &nflEvent.HomeTeam, &nflEvent.Winner)
		if err != nil {
			return nil, err
		}
		nflEvents = append(nflEvents, nflEvent)
	}

	return nflEvents, nil
}
