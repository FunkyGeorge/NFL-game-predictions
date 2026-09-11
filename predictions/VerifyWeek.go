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
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type ScheduleWeek struct {
	Label     string `json:"label"`
	Value     string `json:"value"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

type ScheduleSection struct {
	Label   string          `json:"label"`
	Value   string          `json:"value"`
	Entries *[]ScheduleWeek `json:"entries"`
}

type ScheduleResponse struct {
	Sections []ScheduleSection `json:"sections"`
}

func VerifyWeek() {
	conn, err := sql.Open("sqlite3", "./nfldata.db")

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	scheduleRepository := &database.ScheduleRepository{DB: conn}
	err = scheduleRepository.CreateTable()

	if err != nil {
		fmt.Println(err)
	}

	if config.Week == 0 {
		// If no week input, should default to last week's results +1
		fmt.Println("No week was specified, will try to automatically predict week")

		isSeeded := scheduleRepository.IsSeeded()
		if !isSeeded {
			fmt.Println("Weeks table is empty, will use NFL api to populate")
			GetSchedule(scheduleRepository)
		}

		fmt.Printf("Searching schedule for current week %s\n", time.Now().Format(time.RFC3339))

		week, err := scheduleRepository.FindWeekFromTime(time.Now().Format(time.RFC3339))

		if err != nil {
			fmt.Println("Error querying current week")
		}
		fmt.Printf("Proceeding with season type %s and week %s\n", week.Type, week.Value)
		weekInt, err := strconv.Atoi(week.Value)
		if err != nil {
			fmt.Println("Invalid week value")
		}
		typeInt, err := strconv.Atoi(week.Type)
		if err != nil {
			fmt.Println("Invalid season type value")
		}
		config.Week = weekInt
		config.SeasonType = typeInt
	}
}

func GetSchedule(repo *database.ScheduleRepository) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("https://%s/nfl-whitelist", config.ApiHost), nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-rapidapi-key", config.ApiKey)
	req.Header.Add("x-rapidapi-host", config.ApiHost)

	resp, err := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Printf("Schedule API returned with unexpected status: %i\n", resp.StatusCode)
	}

	if err != nil {
		log.Fatalln(err)
	}

	body, _ := io.ReadAll(resp.Body)

	var jsonResp ScheduleResponse
	err = json.Unmarshal(body, &jsonResp)

	if err != nil {
		log.Fatalln(err)
	}

	acceptedSeasons := []string{"Preseason", "Regular Season", "Postseason"}

	for _, section := range jsonResp.Sections {
		if section.Entries != nil && slices.Contains(acceptedSeasons, section.Label) {
			for _, entry := range *section.Entries {
				fmt.Printf("Inserting %s\n", entry.Label)
				weekInsert := database.Week{
					Type:      section.Value,
					Value:     entry.Value,
					StartDate: entry.StartDate,
					EndDate:   entry.EndDate}
				repo.Insert(weekInsert)
			}
		}
	}
	fmt.Println("Populated the schedule table")
}
