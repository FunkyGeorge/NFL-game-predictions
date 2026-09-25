package database

import "database/sql"

type NFLTeamStatsRepository struct {
	DB *sql.DB
}

type NFLTeamStats struct {
	TeamId               string
	week                 int
	TurnOverDifferential float32
	FourthDownConvs      float32
	PassingBigPlays      float32
	RushingBigPlays      float32
	GamesPlayed          float32
}

func (r *NFLTeamStatsRepository) CreateTable() error {
	_, err := r.DB.Exec(`CREATE TABLE IF NOT EXISTS nflteamstats (
	teamId TEXT PRIMARY KEY,
	week INTEGER,
	fourthDownConvs REAL,
	passingBigPlays REAL,
	rushingBigPlays REAL,
	gamesPlayed REAL
 )`)

	return err
}

func (r *NFLTeamStatsRepository) Insert(row NFLTeamStats) error {
	_, err := r.DB.Exec(`INSERT INTO nflteamstats (
	teamId,
	week,
	turnOverDifferential,
	fourthDownConvs,
	passingBigPlays,
	rushingBigPlays,
	gamesPlayed) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		row.TeamId, row.week,
		row.TurnOverDifferential, row.FourthDownConvs,
		row.PassingBigPlays, row.RushingBigPlays, row.GamesPlayed)

	return err
}

func (r *NFLTeamStatsRepository) GetAll() ([]NFLTeamStats, error) {
	rows, err := r.DB.Query("SELECT * FROM nflteamstats")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var nflTeams []NFLTeamStats

	for rows.Next() {
		var nflTeam NFLTeamStats
		err := rows.Scan(&nflTeam.TeamId, &nflTeam.TurnOverDifferential, &nflTeam.FourthDownConvs,
			&nflTeam.PassingBigPlays, &nflTeam.RushingBigPlays, &nflTeam.GamesPlayed)

		if err != nil {
			return nil, err
		}
		nflTeams = append(nflTeams, nflTeam)
	}

	return nflTeams, nil
}

func (r *NFLTeamStatsRepository) GetById(id string) (NFLTeamStats, error) {
	var nflTeam NFLTeamStats
	err := r.DB.QueryRow("SELECT * FROM nflteamstats WHERE teamId = ?", id).Scan(&nflTeam.TeamId,
		&nflTeam.TurnOverDifferential, &nflTeam.FourthDownConvs,
		&nflTeam.PassingBigPlays, &nflTeam.RushingBigPlays, &nflTeam.GamesPlayed)

	if err != nil {
		return NFLTeamStats{}, err
	}

	return nflTeam, nil
}
