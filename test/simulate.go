package test

import (
	"database/sql"
	// "flag"
	"fmt"
	"log"
	"guess-nfl-winners/config"
	"guess-nfl-winners/database"
	_ "github.com/mattn/go-sqlite3"
)

func Simulate() {
	// flag.Float64Var(&config.TO, "turnovers", 1, "Weight to apply to turnovers")
	// flag.Float64Var(&config.Pass, "passing", 1, "Weight to apply to big passing plays")
	// flag.Float64Var(&config.Run, "running", 1, "Weight to apply to big run plays")
	// flag.Float64Var(&config.Conv, "conv", 1, "Weight to apply to 4th down conversions")
	// flag.Float64Var(&config.Home, "home", 0.5, "Weight to apply to being the home team")
	//
	// flag.Parse()

	conn, err := sql.Open("sqlite3", "./nfldata.db")

	var totalCorrectAssessments int = 0
	var totalGames int = 0

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	nflEventsRepository := &database.NFLEventsRepository{DB: conn}
	nflTeamsRepository := &database.NFLTeamsRepository{DB: conn}

	allEvents, err := nflEventsRepository.GetAll()

	if err != nil {
		log.Fatal(err)
	}

	for _, event := range allEvents {
		totalGames++
		team1, err := nflTeamsRepository.GetById(event.Team1)
		if err != nil {
			log.Fatal(err)
		}
		team2, err := nflTeamsRepository.GetById(event.Team2)
		if err != nil {
			log.Fatal(err)
		}

		var score1 float32 = (
			team1.TurnOverDifferential * float32(config.TO) +
			team1.PassingBigPlays * float32(config.Pass) +
			team1.RushingBigPlays * float32(config.Run) +
			team1.FourthDownConvs * float32(config.Conv)) / team1.GamesPlayed



		var score2 float32 = (
			team2.TurnOverDifferential * float32(config.TO) +
			team2.PassingBigPlays * float32(config.Pass) +
			team2.RushingBigPlays * float32(config.Run) +
			team2.FourthDownConvs * float32(config.Conv)) / team2.GamesPlayed

		if (event.HomeTeam == team1.TeamId) {
			score1 = score1 + float32(config.Home)
		} else {
			score2 = score2 + float32(config.Home)
		}

		winningTeam := ""
		if score1 > score2 {
			winningTeam = team1.TeamId
		} else if score1 < score2 {
			winningTeam = team2.TeamId
		}

		if winningTeam == event.Winner {
			totalCorrectAssessments++
		}
	}

	fmt.Printf("Results: %.2f - TO: %.2f, Pass: %.2f, Rush: %.2f, Conv: %.2f, Home: %.2f\n",
		float32(totalCorrectAssessments)/float32(totalGames) * 100,
		config.TO, config.Pass, config.Run, config.Conv, config.Home)
}
