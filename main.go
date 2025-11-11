package main

import (
	"fmt"
	"flag"
	"os"
	"guess-nfl-winners/predictions"
	"guess-nfl-winners/collect"
	"guess-nfl-winners/config"
	"guess-nfl-winners/test"
)


func main() {
	flag.StringVar(&config.Mode, "mode", "normal", "run in normal, test, or collect mode")
	flag.IntVar(&config.Week, "week", 0, "Week to calculate evaluations")

	flag.Float64Var(&config.TO, "turnovers", 1, "Weight to apply to turnovers")
	flag.Float64Var(&config.Pass, "passing", 1, "Weight to apply to big passing plays")
	flag.Float64Var(&config.Run, "running", 1, "Weight to apply to big run plays")
	flag.Float64Var(&config.Conv, "conv", 1, "Weight to apply to 4th down conversions")
	flag.Float64Var(&config.Home, "home", 0.5, "Weight to apply to being the home team")

	flag.Parse()
	
	switch (config.Mode) {
		case ("normal"):
			if (config.Week == 0) {
				fmt.Println("Must enter a week")
				os.Exit(1)
			}
			var gameIds []string = predictions.GetGameIds(config.Week)
			for _, game := range gameIds {
				predictions.EvaluateGame(game)
			}
		case ("test"):
			test.Simulate()
		case ("collect"):
			if (config.Week == 0) {
				fmt.Println("Must enter a week")
				os.Exit(1)
			}
			collect.FillDb()
		default:
			fmt.Println("Not a valid mode")
	}
}

