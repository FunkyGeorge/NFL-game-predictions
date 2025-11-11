package collect

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"slices"
	"strconv"
	"guess-nfl-winners/config"
	"guess-nfl-winners/predictions"
	"guess-nfl-winners/database"
	_ "github.com/mattn/go-sqlite3"
)

func FillDb () {
	var allTeams []string

	conn, err := sql.Open("sqlite3", "./nfldata.db")

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	nflEventsRepository := &database.NFLEventsRepository{DB: conn}
	nflTeamsRepository := &database.NFLTeamsRepository{DB: conn}

	nflEventsRepository.CreateTable()
	nflTeamsRepository.CreateTable()

	// Get game events for past 3 weeks
	for week := config.Week; week > 0 && week > config.Week - 3; week-- {
		req, _ := http.NewRequest("GET", fmt.Sprintf("https://nfl-api-data.p.rapidapi.com/nfl-weeks-events?year=2025&week=%s&type=2", strconv.Itoa(week)), nil)
		req.Header.Add("x-rapidapi-key", config.ApiKey)
		req.Header.Add("x-rapidapi-host", config.ApiHost)


		resp, err := http.DefaultClient.Do(req)

		if err != nil {
			log.Fatalln(err)
		}
		
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)

		var games predictions.GamesByWeekResponse
		err = json.Unmarshal(body, &games)

		if err != nil {
			log.Fatalln(err)
		}


		for _, gameItem := range games.Items {
			req, _ := http.NewRequest("GET", fmt.Sprintf("https://nfl-api-data.p.rapidapi.com/nfl-single-events?id=%s", gameItem.GameId), nil)
			req.Header.Add("x-rapidapi-key", config.ApiKey)
			req.Header.Add("x-rapidapi-host", config.ApiHost)

			resp, err := http.DefaultClient.Do(req)

			if err != nil {
				log.Fatalln(err)
			}
			
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)

			var event predictions.EventResponse
			err = json.Unmarshal(body, &event)

			if err != nil {
				log.Fatalln(err)
			}

			var teams []string
			eventRecord := database.NFLEvent{ GameId: gameItem.GameId}
			
			for index, team := range event.Competitions[0].Competitors {
				teams = append(teams, team.Id)
				
				if (!slices.Contains(allTeams, team.Id)) {
					allTeams = append(allTeams, team.Id)
				}
				
				if (index == 0) {
					eventRecord.Team1 = team.Id
				}

				if (index == 1) {
					eventRecord.Team2 = team.Id
				}

				if (team.HomeAway == "home") {
					eventRecord.HomeTeam = team.Id
				}

				if (team.Winner) {
					eventRecord.Winner = team.Id
				}
			}

			nflEventsRepository.Insert(eventRecord)
		}
	}

	// Get team statistics for important stats for all teams
	for _, teamId := range allTeams {
		
		req, _ := http.NewRequest("GET", fmt.Sprintf("https://nfl-api-data.p.rapidapi.com/nfl-team-statistics?id=%s&year=2025", teamId), nil)
		req.Header.Add("x-rapidapi-key", config.ApiKey)
		req.Header.Add("x-rapidapi-host", config.ApiHost)

		resp, err := http.DefaultClient.Do(req)

		if err != nil {
			log.Fatalln(err)
		}
		
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)

		var teamStats predictions.TeamStats
		err = json.Unmarshal(body, &teamStats)

		if err != nil {
			log.Fatalln(err)
		}


		defaultPlayed := config.Week
		if err != nil {
			log.Fatalln(err)
		}

		var gamesPlayed float32 = float32(defaultPlayed)
		var savedStats database.NFLTeam

		savedStats.TeamId = teamId

		for _, category := range teamStats.Statistics.Splits.Categories {
			for _, stat := range category.Stats {
				switch (stat.Name) {
					case "turnOverDifferential":
						savedStats.TurnOverDifferential = stat.Value
					case "fourthDownConvs":
						savedStats.FourthDownConvs = stat.Value
					case "passingBigPlays":
						savedStats.PassingBigPlays = stat.Value
					case "rushingBigPlays":
						savedStats.RushingBigPlays = stat.Value
				}

				if stat.Name == "gamesPlayed" {
					gamesPlayed = stat.Value
				}
			}
		}
		savedStats.GamesPlayed = gamesPlayed

		nflTeamsRepository.Insert(savedStats)
	}
}
