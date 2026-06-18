# NFL weekly winners guesser
Using API data and a low effort data calculation, this CLI tool guesses who the winners of an NFL week will be.
The main statistics that are tracked are turnover differential, big plays, and home team advantage. Weights have been
tweaked to try to come out with the best estimates. The current accuracy of this low effort algorithm is about 74% which
is only slightly lower than other standard outcome predictors like Yahoo.

### Usage
Check the config/api.go file to see what environment variables you need to set. This tool relies one having a rapidapi account and 
API key for the NFL API Data.

```
go build
./guess-nfl-winners --mode=normal --week=<Enter Week>
```
This will pull the latest data from the NFL api and run the calculations each time

### Testing
At some point I needed to check how accurate my algorithm was with historical data.
I created some commands to pull data and store it so that I can continue running a validation calculation
with stored data to prevent running up API costs

```
./guess-nfl-winners --mode=collect
```
This will pull and store the required NFL data in sqlite

```
./guess-nfl-winners --mode=test --week=<Enter Week>
```
This runs the calculation against stored data from the collect command and compares results to the actual winners.
It will give a percentage accuracy. This allows for tweaking the weights of the collected data to try to find 
optimal weights for our data

### Tweaking Weights
Settings as listed in config/flags give options to override the weights of our data
