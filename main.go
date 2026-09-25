package main

import (
	"flag"
	"fmt"
	"guess-nfl-winners/config"
	"guess-nfl-winners/predictions"
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
		complete := true
		output := ""
		for _, game := range gameIds {
			if result, errored := predictions.EvaluateGame(game); errored {
				complete = false
			} else {
				output = fmt.Sprintf("%s%s\n", output, result)
			}
		}
		if !complete {
			fmt.Println("[WARNING]: Not all games have evaluated, results incomplete")
		}
		fmt.Print(output)
	case ("test"):
		fmt.Println("Not implemented yet")
	default:
		fmt.Println("Not a valid mode")
	}
}
