
package database

import "database/sql"

type NFLTeamsRepository struct {
	DB *sql.DB
}

type NFLTeam struct {
	TeamId string
	TurnOverDifferential float32
	FourthDownConvs float32
	PassingBigPlays float32
	RushingBigPlays float32
	GamesPlayed float32
}

func (r *NFLTeamsRepository) CreateTable() error {
 _, err := r.DB.Exec(`CREATE TABLE IF NOT EXISTS nflteams (
	teamId TEXT PRIMARY KEY,
	turnOverDifferential REAL,
	fourthDownConvs REAL,
	passingBigPlays REAL,
	rushingBigPlays REAL,
	gamesPlayed REAL
 )`)

 _, err = r.DB.Exec("DELETE FROM nflteams")

 return err
}

func (r *NFLTeamsRepository) Insert(row NFLTeam) error {
	_, err := r.DB.Exec(`INSERT INTO nflteams (
	teamId,
	turnOverDifferential,
	fourthDownConvs,
	passingBigPlays,
	rushingBigPlays,
	gamesPlayed) VALUES (?, ?, ?, ?, ?, ?)`,
	row.TeamId, row.TurnOverDifferential, row.FourthDownConvs,
	row.PassingBigPlays, row.RushingBigPlays, row.GamesPlayed)

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
		err := rows.Scan(&nflTeam.TeamId, &nflTeam.TurnOverDifferential, &nflTeam.FourthDownConvs,
		&nflTeam.PassingBigPlays, &nflTeam.RushingBigPlays, &nflTeam.GamesPlayed)

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
		&nflTeam.TurnOverDifferential, &nflTeam.FourthDownConvs,
		&nflTeam.PassingBigPlays, &nflTeam.RushingBigPlays, &nflTeam.GamesPlayed)

	if err != nil {
		return NFLTeam{}, err
	}

	return nflTeam, nil
}
