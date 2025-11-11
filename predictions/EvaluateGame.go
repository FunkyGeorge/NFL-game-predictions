package predictions

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"guess-nfl-winners/config"
)

type Team struct {
	Id string `json:"id"`
	DisplayName string `json:"displayName"`
}

type Competitor struct {
	Id string `json:"id"`
	CompetitorTeam Team `json:"team"`
	HomeAway string `json:"homeAway"`
	Winner bool `json:"winner"`
}

type Competition struct {
	Competitors []Competitor `json:"competitors"`
}

type EventResponse struct {
	Name string `json:"name"`
	ShortName string `json:"shortName"`
	Competitions []Competition `json:"competitions"`
}

func EvaluateGame(gameId string) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("https://nfl-api-data.p.rapidapi.com/nfl-single-events?id=%s", gameId), nil)
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

	resultString := ""
	teamIndex1 := GetTeamImpactIndex(event.Competitions[0].Competitors[0].Id) + 0.2 // add 0.2 to home team
	teamIndex2 := GetTeamImpactIndex(event.Competitions[0].Competitors[1].Id)

	if teamIndex1 > teamIndex2 {
		resultString = fmt.Sprintf("%s%s Wins!;", resultString, event.Competitions[0].Competitors[0].CompetitorTeam.DisplayName)
	} else if teamIndex1 < teamIndex2 {
		resultString = fmt.Sprintf("%s%s Wins!;", resultString, event.Competitions[0].Competitors[1].CompetitorTeam.DisplayName)
	} else {
		resultString = fmt.Sprintf("%s Tie;", resultString)
	}

	resultString = fmt.Sprintf("%s%s - ", resultString, event.Competitions[0].Competitors[0].CompetitorTeam.DisplayName)
	resultString = fmt.Sprintf("%s%f; %s - ", resultString, teamIndex1, event.Competitions[0].Competitors[1].CompetitorTeam.DisplayName)
	resultString = fmt.Sprintf("%s%f\n", resultString, teamIndex2)

	fmt.Println(resultString)
}
