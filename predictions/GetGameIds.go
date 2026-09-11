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
	"strconv"
)

type GameItem struct {
	GameId string `json:"eventid"`
}

type GamesByWeekResponse struct {
	Items []GameItem `json:"items"`
}

func GetGameIds(week int) []string {
	// Check db for saved weeks first
	var gameIds []string
	conn, err := sql.Open("sqlite3", "./nfldata.db")

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	gbwRepo := &database.GamesByWeekRepository{DB: conn}
	err = gbwRepo.CreateTable()
	if err != nil {
		fmt.Println("Error creating GamesByWeek Table")
	}

	events, err := gbwRepo.GetAllByWeek(config.Week, config.SeasonType)
	if err != nil {
		fmt.Println("Error querying events ids from table")
	}

	if len(events) == 0 {
		fmt.Println("Don't have events information yet... querying API")

		req, _ := http.NewRequest("GET", fmt.Sprintf(
			"https://%s/nfl-weeks-events?year=2025&week=%s&type=%s",
			config.ApiHost,
			strconv.Itoa(week),
			strconv.Itoa(config.SeasonType)), nil)
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

		fmt.Println("Inserting games by week")
		for _, gameItem := range games.Items {
			gameByWeek := database.GameByWeek{
				EventId:    gameItem.GameId,
				Week:       strconv.Itoa(week),
				SeasonType: strconv.Itoa(config.SeasonType)}
			fmt.Println(gameByWeek)
			err = gbwRepo.Insert(gameByWeek)
			if err != nil {
				fmt.Printf("Could not insert game %s\n", gameItem.GameId)
			}
			gameIds = append(gameIds, gameItem.GameId)
		}
	} else {
		fmt.Println("Found events in db already")
		for _, event := range events {
			gameIds = append(gameIds, event.EventId)
		}
	}
	return gameIds
}
