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
	conn, err := sql.Open("sqlite3", "./nfldata.db")

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	// TODO: This can probably be optimized not opening new connections for sqlite
	nflEventsRepository := &database.NFLEventsRepository{DB: conn}
	nflEventsRepository.CreateTable()

	nflTeamsRepository := &database.NFLTeamsRepository{DB: conn}
	nflTeamsRepository.CreateTable()

	nflTeamStatsRepository := &database.NFLTeamStatsRepository{DB: conn}
	nflTeamStatsRepository.CreateTable()

	// Build NFL teams map
	allTeams, err := nflTeamsRepository.GetAll()

	teamMap := make(map[string]string)

	for _, t := range allTeams {
		teamMap[t.TeamId] = t.DisplayName
	}

	if err != nil {
		fmt.Println("Error querying all NFL teams")
	}

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

		if resp.StatusCode != 200 {
			fmt.Printf("Skipping game %s. Likely due to API error. Try again later for result",
				gameId)
			return
		}

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

		eventRecord = database.NFLEvent{GameId: gameId}

		fmt.Println("Going through competitions", event.Competitions)
		// TODO: Currently issue here when api fails with exceeded quota. Seems like api side error
		for index, team := range event.Competitions[0].Competitors {
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

			if _, exists := teamMap[team.Id]; !exists {
				var newTeam database.NFLTeam
				newTeam.TeamId = team.Id
				newTeam.DisplayName = team.CompetitorTeam.DisplayName
				err = nflTeamsRepository.Insert(newTeam)

				if err != nil {
					fmt.Println("Could not insert new NFL team")
				}
				teamMap[team.Id] = team.CompetitorTeam.DisplayName
			}
		}

		err = nflEventsRepository.Insert(eventRecord)
		if err != nil {
			log.Fatal(err)
		}
	}

	resultString := ""
	teamIndex1 := GetTeamImpactIndex(eventRecord.HomeTeam)
	teamIndex1 = teamIndex1 + float32(config.Home) // add 0.2 to home team

	teamIndex2 := GetTeamImpactIndex(eventRecord.AwayTeam)

	if teamIndex1 > teamIndex2 {
		resultString = fmt.Sprintf("%s%s Wins!;", resultString, teamMap[eventRecord.HomeTeam])
	} else if teamIndex1 < teamIndex2 {
		resultString = fmt.Sprintf("%s%s Wins!;", resultString, teamMap[eventRecord.AwayTeam])
	} else {
		resultString = fmt.Sprintf("%s Tie;", resultString)
	}

	resultString = fmt.Sprintf("%s%s - ", resultString, teamMap[eventRecord.HomeTeam])
	resultString = fmt.Sprintf("%s%f; %s - ", resultString, teamIndex1, teamMap[eventRecord.AwayTeam])
	resultString = fmt.Sprintf("%s%f\n", resultString, teamIndex2)

	fmt.Println(resultString)
}
