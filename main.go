package main

import (
	"flag"
	"fmt"
	"guess-nfl-winners/config"
	"guess-nfl-winners/predictions"
	// "guess-nfl-winners/test"
	"os"
)

func main() {
	if config.ApiKey == "" {
		fmt.Println("Couldn't find ApiKey, make sure to set ApiKey in .env")
		os.Exit(1)
	}

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
		predictions.VerifyWeek()
		gameIds := predictions.GetGameIds(config.Week)
		fmt.Println("\nStarting Evaluations...")
		for _, game := range gameIds {
			predictions.EvaluateGame(game)
		}
		os.Exit(1)
	case ("test"):
		fmt.Println("Not implemented yet")
	default:
		fmt.Println("Not a valid mode")
	}
}
