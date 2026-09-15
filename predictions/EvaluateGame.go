package predictions

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"guess-nfl-winners/config"
	"guess-nfl-winners/database"
	"io"
	"log"
	"net/http"
	"os"
	"slices"
)

type Team struct {
	Id          string `json:"id"`
	DisplayName string `json:"displayName"`
}

type Competitor struct {
	Id             string `json:"id"`
	CompetitorTeam Team   `json:"team"`
	HomeAway       string `json:"homeAway"`
	Winner         bool   `json:"winner"`
}

type Competition struct {
	Competitors []Competitor `json:"competitors"`
}

type EventResponse struct {
	Name         string        `json:"name"`
	ShortName    string        `json:"shortName"`
	Competitions []Competition `json:"competitions"`
}

func EvaluateGame(gameId string) {
	// TODO: cache games in db and try to calculate per week relevant stats

	conn, err := sql.Open("sqlite3", "./nfldata.db")

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	nflEventsRepository := &database.NFLEventsRepository{DB: conn}
	nflTeamsRepository := &database.NFLTeamsRepository{DB: conn}

	nflEventsRepository.CreateTable()
	nflTeamsRepository.CreateTable()

	// Check db for event record
	eventRecord, err := nflEventsRepository.Find(gameId)
	foundRecord := err == nil

	if !foundRecord {
		fmt.Println("Missing game record. Querying API and storing")
		var bothTeams []string
		req, _ := http.NewRequest("GET", fmt.Sprintf("https://%s/nfl-single-events?id=%s",
			config.ApiHost, gameId), nil)
		req.Header.Add("x-rapidapi-key", config.ApiKey)
		req.Header.Add("x-rapidapi-host", config.ApiHost)

		resp, err := http.DefaultClient.Do(req)

		if err != nil {
			log.Fatalln(err)
		}

		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)

		var event EventResponse
		err = json.Unmarshal(body, &event)

		if err != nil {
			log.Fatalln(err)
		}

		// Write to db
		var teams []string
		eventRecord = database.NFLEvent{GameId: gameId}

		for index, team := range event.Competitions[0].Competitors {
			teams = append(teams, team.Id)

			if !slices.Contains(bothTeams, team.Id) {
				bothTeams = append(bothTeams, team.Id)
			}

			if index == 0 {
				eventRecord.Team1 = team.Id
			}

			if index == 1 {
				eventRecord.Team2 = team.Id
			}

			if team.HomeAway == "home" {
				eventRecord.HomeTeam = team.Id
			}

			if team.HomeAway == "away" {
				eventRecord.AwayTeam = team.Id
			}

			if team.Winner {
				eventRecord.Winner = team.Id
			}
		}

		err = nflEventsRepository.Insert(eventRecord)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Print(eventRecord)
	fmt.Println("End of test")
	os.Exit(1)
	resultString := ""
	teamName1, teamIndex1 := GetTeamImpactIndex(eventRecord.HomeTeam)
	teamIndex1 = teamIndex1 + float32(config.Home) // add 0.2 to home team

	teamName2, teamIndex2 := GetTeamImpactIndex(eventRecord.AwayTeam)

	if teamIndex1 > teamIndex2 {
		resultString = fmt.Sprintf("%s%s Wins!;", resultString, teamName1)
	} else if teamIndex1 < teamIndex2 {
		resultString = fmt.Sprintf("%s%s Wins!;", resultString, teamName2)
	} else {
		resultString = fmt.Sprintf("%s Tie;", resultString)
	}

	resultString = fmt.Sprintf("%s%s - ", resultString, teamName1)
	resultString = fmt.Sprintf("%s%f; %s - ", resultString, teamIndex1, teamName2)
	resultString = fmt.Sprintf("%s%f\n", resultString, teamIndex2)

	fmt.Println(resultString)
}
