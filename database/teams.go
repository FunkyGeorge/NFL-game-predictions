package database

import "database/sql"

type NFLTeamsRepository struct {
	DB *sql.DB
}

type NFLTeam struct {
	TeamId      string
	DisplayName string
}

func (r *NFLTeamsRepository) CreateTable() error {
	_, err := r.DB.Exec(`CREATE TABLE IF NOT EXISTS nflteams (
	teamId TEXT PRIMARY KEY, displayName TEXT)`)

	return err
}

func (r *NFLTeamsRepository) Insert(row NFLTeam) error {
	_, err := r.DB.Exec(`INSERT INTO nflteams (teamId, displayName) VALUES (?, ?)`,
		row.TeamId, row.DisplayName)

	return err
}

func (r *NFLTeamsRepository) GetAll() ([]NFLTeam, error) {
	rows, err := r.DB.Query("SELECT * FROM nflteams")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var nflTeams []NFLTeam

	for rows.Next() {
		var nflTeam NFLTeam
		err := rows.Scan(&nflTeam.TeamId, &nflTeam.DisplayName)

		if err != nil {
			return nil, err
		}
		nflTeams = append(nflTeams, nflTeam)
	}

	return nflTeams, nil
}

func (r *NFLTeamsRepository) GetById(id string) (NFLTeam, error) {
	var nflTeam NFLTeam
	err := r.DB.QueryRow("SELECT * FROM nflteams WHERE teamId = ?", id).Scan(&nflTeam.TeamId,
		&nflTeam.DisplayName)

	if err != nil {
		return NFLTeam{}, err
	}

	return nflTeam, nil
}
