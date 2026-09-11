package database

import (
	"database/sql"
)

type GamesByWeekRepository struct {
	DB *sql.DB
}

type GameByWeek struct {
	EventId    string
	Week       string
	SeasonType string
}

func (r *GamesByWeekRepository) CreateTable() error {
	_, err := r.DB.Exec(`CREATE TABLE IF NOT EXISTS gamesByWeek (
	 eventId TEXT PRIMARY KEY,
	 week TEXT,
	 seasonType TEXT
 )`)

	return err
}

func (r *GamesByWeekRepository) Insert(row GameByWeek) error {
	_, err := r.DB.Exec(`INSERT INTO gamesByWeek (
		eventId,
		week,
		seasonType) VALUES (?, ?, ?)`,
		row.EventId, row.Week, row.SeasonType)

	return err
}

func (r *GamesByWeekRepository) GetAllByWeek(week int, seasonType int) ([]GameByWeek, error) {
	rows, err := r.DB.Query("SELECT * FROM gamesByWeek WHERE week = ? AND seasonType = ?", week, seasonType)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var games []GameByWeek

	for rows.Next() {
		var game GameByWeek
		err := rows.Scan(&game.EventId, &game.Week, &game.SeasonType)
		if err != nil {
			return nil, err
		}
		games = append(games, game)
	}

	return games, nil
}
