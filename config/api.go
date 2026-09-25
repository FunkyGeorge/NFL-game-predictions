package config

import "os"

var ApiKey string = os.Getenv("RAPID_API_KEY")
var ApiHost string = "nfl-api-data.p.rapidapi.com"
