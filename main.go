package main

import (
	"flag"
	"fmt"
	// "guess-nfl-winners/collect"
	"guess-nfl-winners/config"
	"guess-nfl-winners/predictions"
	// "guess-nfl-winners/test"
	"os"
	"time"
)

func main() {
	flag.StringVar(&config.Mode, "mode", "normal", "run in normal, test, or collect mode")
	flag.IntVar(&config.Week, "week", 0, "Week to calculate evaluations")

	flag.Float64Var(&config.TO, "turnovers", 0.8, "Weight to apply to turnovers")
	flag.Float64Var(&config.Pass, "passing", 0.6, "Weight to apply to big passing plays")
	flag.Float64Var(&config.Run, "running", 0.6, "Weight to apply to big run plays")
	flag.Float64Var(&config.Home, "home", 0.2, "Weight to apply to being the home team")

	flag.Parse()
	fmt.Printf("Running in %s mode\n", config.Mode)

	switch config.Mode {
	case ("normal"):
		if config.Week == 0 {
			// If no week input, should default to last week's results +1
			fmt.Println("Normal mode requires a week flag")
			fmt.Println(time.Now())
			os.Exit(1)
		}

		// Check if prediction has required data for week

		var gameIds []string = predictions.GetGameIds(config.Week)
		for _, game := range gameIds {
			predictions.EvaluateGame(game)
		}
	// case ("test"):
	// 	test.Simulate()
	// case ("collect"):
	// 	if config.Week == 0 {
	// 		fmt.Println("Must enter a week")
	// 		os.Exit(1)
	// 	}
	// 	collect.FillDb()
	default:
		fmt.Println("Not a valid mode")
	}
}
