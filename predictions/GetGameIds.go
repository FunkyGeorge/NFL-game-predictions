package predictions

import (
	"encoding/json"
	"fmt"
	"guess-nfl-winners/config"
	"io"
	"log"
	"net/http"
	"strconv"
)

type GameItem struct {
	GameId string `json:"eventid"`
}

type GamesByWeekResponse struct {
	Items []GameItem `json:"items"`
}

func GetGameIds(week int) []string {
	// Currently set for postseason type 3
	req, _ := http.NewRequest("GET", fmt.Sprintf("https://nfl-api-data.p.rapidapi.com/nfl-weeks-events?year=2025&week=%s&type=3", strconv.Itoa(week)), nil)
	req.Header.Add("x-rapidapi-key", config.ApiKey)
	req.Header.Add("x-rapidapi-host", config.ApiHost)

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var games GamesByWeekResponse
	err = json.Unmarshal(body, &games)

	if err != nil {
		log.Fatalln(err)
	}

	var gameIds []string

	for _, gameItem := range games.Items {
		gameIds = append(gameIds, gameItem.GameId)
	}

	return gameIds
}
