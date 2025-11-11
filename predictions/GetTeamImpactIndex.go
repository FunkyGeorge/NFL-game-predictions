package predictions

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"slices"
	"guess-nfl-winners/config"
)

type NFLTeam struct {
	id string
	slug string
	abbreviation string
	displayName string
}

type Stat struct {
	Name string `json:"name"`
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
	Value float32 `json:"value"`
	Rank int `json:"rank"`
}

type Category struct {
	DisplayName string `json:"displayName"`
	Summary string `json:"summary"`
	Stats []Stat `json:"stats"`
}

type TeamStats struct {
	Statistics struct {
		Splits struct {
			Id string `json:"id"`
			Categories []Category `json:"categories"`
		} `json:"splits"`
	} `json:"statistics"`
}

var ImportantStats []string = []string{
	"turnOverDifferential",
	"passingBigPlays",
	"rushingBigPlays"}

func GetTeamImpactIndex(teamId string) float32 {
	req, _ := http.NewRequest("GET", fmt.Sprintf("https://nfl-api-data.p.rapidapi.com/nfl-team-statistics?id=%s&year=2025", teamId), nil)
	req.Header.Add("x-rapidapi-key", config.ApiKey)
	req.Header.Add("x-rapidapi-host", config.ApiHost)

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Fatalln(err)
	}
	
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var teamStats TeamStats
	err = json.Unmarshal(body, &teamStats)

	if err != nil {
		log.Fatalln(err)
	}


	defaultPlayed := config.Week
	if err != nil {
		log.Fatalln(err)
	}

	var gamesPlayed float32 = float32(defaultPlayed)
	var impactIndex float32 = 0

	for _, category := range teamStats.Statistics.Splits.Categories {
		for _, stat := range category.Stats {
			if slices.Contains(ImportantStats, stat.Name) {
				if stat.Name == "turnOverDifferential" {
					impactIndex += stat.Value * 0.8
				}
				if stat.Name == "passingBigPlays" {
					impactIndex += stat.Value * 0.6
				}
				if stat.Name == "rushingBigPlays" {
					impactIndex += stat.Value * 0.6
				}
			}

			if stat.Name == "gamesPlayed" {
				gamesPlayed = stat.Value
			}
		}
	}

	return impactIndex / gamesPlayed
}
