package process

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/KunalDuran/duranzcli/module/data"
)

// var DATASET_BASE = `/home/kunalduran/Desktop/duranz_api/develop/datasets/`
//var PWD string
var PWD, _ = os.Getwd()

var DATASET_BASE = PWD + `/datasets/`

var SLASH = `/`

func Activate(processName, fileName string) {
	switch processName {

	case "initial":
		data.CreateAllTables()

	case "delete":
		data.DeleteAllTableData()

	case "all", "venue", "players", "teams", "match", "matchstats", "playerstats":
		RunAllProcess(processName, fileName)
	}
}

func RunAllProcess(process, fileName string) {

	var newFiles []string
	if fileName != "" {
		newFiles = append(newFiles, fileName)
	} else {
		allFiles := ListFiles(DATASET_BASE)
		newFiles = GetNewFiles(allFiles)
	}

	for _, file := range newFiles {
		var mappingInfo data.MappingInfo
		match, err := GetCricsheetData(DATASET_BASE + file)
		if err != nil {
			fmt.Println(err.Error())
			data.InsertErrorLog(data.CRICSHEET_FILE_ERROR, `Error in fetching file `, file)
			return
		}

		mappingInfo.LeagueID = data.AllDuranzLeagues[match.Info.MatchType]

		// fmt.Println(string(match.Info.Season))
		if match.Info.TeamType == "club" && (match.Info.Event.Name == "Indian Premier League" || strings.ToLower(match.Info.Event.Name) == "ipl") {
			match.Info.TeamType = "ipl"
		}

		if process == "venue" || process == "all" {
			// Map the VENUES
			VenueMapper(match.Info.Venue, match.Info.City)
			mappingInfo.Venue = true
		}

		if process == "players" || process == "all" {
			// Map the Players
			PlayerMapper(match.Info.Register.People)
			mappingInfo.Players = true

		}

		if process == "teams" || process == "all" {
			// Map the Teams
			TeamMapper(match.Info.Teams, match.Info.TeamType)
			mappingInfo.Teams = true

		}

		if process == "match" || process == "all" {
			// Map the matches
			MatchMapper(match, file)
			mappingInfo.Match = true

		}

		if process == "matchstats" || process == "all" {
			// Map the match stats
			processMatchStats(match, file)
			mappingInfo.MatchStats = true

		}

		if process == "playerstats" || process == "all" {
			// Map the player stats
			processPlayerStats(match, file)
			mappingInfo.PlayerStats = true
		}

		data.InsertMappingInfo(file, mappingInfo)
	}
}

// GetCricsheetData : Reads the match json file
func GetCricsheetData(f_path string) (data.Match, error) {
	var matchData data.Match
	body, err := ioutil.ReadFile(f_path)
	if err != nil {
		return matchData, err
	}

	err = json.Unmarshal(body, &matchData)
	if err != nil {
		fmt.Println(err, "in file ", f_path)
		return matchData, err
	}
	return matchData, nil
}

// PENDING
func processScoreCard() {
	match, err := GetCricsheetData(`C:\Users\Kunal\Desktop\Duranz\duranz_api\matchdata\odis_json\433606.json`)
	if err != nil {
		fmt.Println(err.Error())
	}

	// Map the VENUES
	// VenueMapper(match.Info.Venue, match.Info.City)

	// MAP THE TEAMS

	// MAP THE PLAYERS

	var objScoreCard data.ScoreCard

	var AllInnings []data.Innings
	for _, inning := range match.Innings {
		fmt.Println("Scorecard process started innings for : ", inning.Team)

		var objInning data.Innings
		objInning.InningDetail = inning.Team

		var objExtra data.Extras
		var objBatsman = map[string]data.BattingStats{}
		var objBowler = map[string]data.BowlingStats{}
		var batsmanCount, bowlerCount, runningScore, wicketCnt int
		var fowArr []string
		var overRuns = map[int]int{}
		for _, over := range inning.Overs {
			overRuns[over.Over] = 0
			for _, delivery := range over.Deliveries {

				// Score Calculations
				runningScore += delivery.Runs.Total
				overRuns[over.Over] += delivery.Runs.Total

				// Batsman Init
				if _, exist := objBatsman[delivery.Batter]; !exist {
					batsmanCount++
					var tempBat data.BattingStats
					tempBat.BattingOrder = batsmanCount
					tempBat.Name = delivery.Batter
					objBatsman[delivery.Batter] = tempBat
				}
				batsman := objBatsman[delivery.Batter]
				batsman.Runs += delivery.Runs.Batter
				batsman.Balls += 1

				// Bowler Init
				if _, exist := objBowler[delivery.Bowler]; !exist {
					bowlerCount++
					var tempBowler data.BowlingStats
					tempBowler.BowlingOrder = bowlerCount
					tempBowler.Name = delivery.Bowler
					objBowler[delivery.Bowler] = tempBowler
				}
				bowler := objBowler[delivery.Bowler]
				bowler.Runs += delivery.Runs.Batter
				bowler.Balls += 1

				if delivery.Runs.Batter == 4 {
					batsman.Fours++
				} else if delivery.Runs.Batter == 6 {
					batsman.Sixes++
				}

				// Calculate Extras
				if delivery.Extras != (data.Extras{}) {
					if delivery.Extras.Byes > 0 {
						objExtra.Byes += delivery.Extras.Byes
						overRuns[over.Over] -= delivery.Extras.Byes
					} else if delivery.Extras.LegByes > 0 {
						objExtra.LegByes += delivery.Extras.LegByes
						overRuns[over.Over] -= delivery.Extras.LegByes
					} else if delivery.Extras.NoBall > 0 {
						// remove ball count if No Ball
						batsman.Balls -= 1
						bowler.Balls -= 1
						bowler.Runs += delivery.Extras.NoBall
						objExtra.NoBall += delivery.Extras.NoBall
					} else if delivery.Extras.Wides > 0 {
						// remove ball count if Wide Ball
						batsman.Balls -= 1
						bowler.Balls -= 1
						bowler.Runs += delivery.Extras.Wides
						objExtra.Wides += delivery.Extras.Wides
					}
				}

				// Check for Wicket
				for _, wicket := range delivery.Wickets {
					if wicket.Kind != "" && wicket.PlayerOut != "" {
						batsman.Out = wicket.Kind
						wicketCnt++
						fowStr := fmt.Sprint(wicketCnt, "-", runningScore, "(", wicket.PlayerOut, ")")
						fowArr = append(fowArr, fowStr)

						// bowler
						if wicket.Kind != "run out" {
							bowler.Wickets++
						}
					}
				}

				// bind all info and calculations
				objBatsman[delivery.Batter] = batsman
				objBowler[delivery.Bowler] = bowler
			}

			// check maiden over
			if val, ok := overRuns[over.Over]; ok && val == 0 {
				if len(over.Deliveries) > 0 {
					bowler := objBowler[over.Deliveries[0].Bowler]
					bowler.Maiden++
					objBowler[over.Deliveries[0].Bowler] = bowler
				}
			}
		}

		var allBatsman []data.BattingStats
		for _, batter := range objBatsman {
			if batter.Balls > 0 {
				batter.StrikeRate = math.Round((float64(batter.Runs)*100)/float64(batter.Balls)/0.01) * 0.01
			}
			if batter.Out == "" {
				batter.Out = "not out"
			}
			allBatsman = append(allBatsman, batter)
		}

		var allBowler []data.BowlingStats
		for _, bowler := range objBowler {
			if bowler.Balls > 0 {
				bowler.Economy = math.Round(float64(bowler.Runs)/(float64(bowler.Balls)/float64(6))/0.01) * 0.01
			}
			bowler.Overs = fmt.Sprint(bowler.Balls/6) + "." + fmt.Sprint(bowler.Balls%6)
			allBowler = append(allBowler, bowler)
		}

		objExtra.Total = objExtra.Byes + objExtra.LegByes + objExtra.Wides + objExtra.NoBall
		objInning.Extras = objExtra
		objInning.Batting = allBatsman
		objInning.FallOfWickets = strings.Join(fowArr, " , ")
		objInning.Bowling = allBowler

		AllInnings = append(AllInnings, objInning)
	}

	resultStr := match.Info.Outcome.Winner + " Won by "
	if match.Info.Outcome.By.Runs > 0 {
		resultStr += strconv.Itoa(match.Info.Outcome.By.Runs) + " Runs"
	} else if match.Info.Outcome.By.Wickets > 0 {
		resultStr += strconv.Itoa(match.Info.Outcome.By.Wickets) + " Wickets"
	}
	objScoreCard.Result = resultStr
	objScoreCard.Innings = AllInnings

	strScoreCard, err := json.MarshalIndent(objScoreCard, "", "  ")
	if err != nil {
		fmt.Println(err.Error())
	}

	ioutil.WriteFile(`C:\Users\Kunal\Desktop\Duranz\duranz_api\scoreCard.json`, strScoreCard, 0777)
}

func processPlayerStats(match data.Match, fileName string) {

	var hasErrors bool
	// extract year
	var seasonID int
	if len(match.Info.Dates) > 0 {
		matchDate, err := time.Parse("2006-01-02", match.Info.Dates[0])
		if err != nil {
			data.InsertErrorLog(data.DATETIME_ERROR, `Error in parsing date`+match.Info.Dates[0], fileName)
			panic(err)
		}
		seasonID = matchDate.Year()
	}
	// get match id, else log error and continue
	cricSheetID := strings.Replace(fileName, ".json", "", -1)
	matchID := data.GetMatchID(cricSheetID)

	// get team id
	home := match.Info.Teams[0]
	away := match.Info.Teams[1]
	homeTeamID := data.GetTeamID(home, match.Info.TeamType)
	awayTeamID := data.GetTeamID(away, match.Info.TeamType)

	allPlayerID := map[string]int{}
	for player, cricID := range match.Info.Register.People {
		playerID := data.GetPlayerID(cricID)
		if playerID == 0 {
			// log error
			data.InsertErrorLog(data.PLAYER_NOT_FOUND, `player not found`+player, fileName)
			hasErrors = true
			fmt.Println("ID of this player is zero ", player)
			continue
		}
		allPlayerID[player] = playerID
	}

	if hasErrors {
		fmt.Println("STOP!! Process have serious errors to be fixed")
		return
	}

	var teamInningPlayerStats = map[string]map[string]data.PlayerStats{}
	var innBatTeamMap = map[int]int{}

	for inningID, inning := range match.Innings {
		if inning.SuperOvers {
			continue
		}
		inningID++
		// fmt.Println("player Stats process started innings for : ", inning.Team)

		var battingTeamID, bowlingTeamID int
		if inning.Team == home {
			battingTeamID = homeTeamID
			bowlingTeamID = awayTeamID
		} else if inning.Team == away {
			battingTeamID = awayTeamID
			bowlingTeamID = homeTeamID
		}

		// for inning id in test also
		innBatTeamMap[battingTeamID] = inningID

		var objBatsman = map[string]data.BattingStats{}
		var objBowler = map[string]data.BowlingStats{}
		var objFielder = map[string]data.FieldingStats{}

		var battingOrder, bowlingOrder int
		var overRunsBowler = map[int]int{}

		for _, over := range inning.Overs {

			overRunsBowler[over.Over] = 0
			for _, delivery := range over.Deliveries {
				// init batsman
				batsman, existBat := objBatsman[delivery.Batter]
				if !existBat {
					battingOrder++
					batsman.BattingOrder = battingOrder
					batsman.Name = delivery.Batter
					batsman.IsBatted = true
				}

				// init non-striker
				nonBatsman, existBat2 := objBatsman[delivery.NonStriker]

				if !existBat2 {
					battingOrder++
					nonBatsman.BattingOrder = battingOrder
					nonBatsman.Name = delivery.NonStriker
					nonBatsman.IsBatted = true
					objBatsman[delivery.NonStriker] = nonBatsman
				}

				// init bowler
				bowler, existBowl := objBowler[delivery.Bowler]
				if !existBowl {
					bowlingOrder++
					bowler.BowlingOrder = bowlingOrder
					bowler.Name = delivery.Bowler
				}

				// init fielder/WK and calculate
				if len(delivery.Wickets) > 0 {
					for _, wicket := range delivery.Wickets {
						for _, f := range wicket.Fielders {
							if f.Substitute && f.Name == "" {
								continue
							}
							fielder, existField := objFielder[f.Name]
							if !existField {
								fielder.Name = f.Name
							}
							if wicket.Kind == "caught" {
								fielder.Catches++
							} else if wicket.Kind == "run out" {
								fielder.RunOuts++
							} else if wicket.Kind == "stumped" {
								fielder.Stumpings++
							}
							objFielder[f.Name] = fielder
						}
						if wicket.Kind == "caught and bowled" {
							fielder := objFielder[bowler.Name]
							fielder.Name = bowler.Name
							fielder.Catches++
							objFielder[bowler.Name] = fielder
						}
					}
				}

				// ================ bowler calculations
				bowlerRuns := delivery.Runs.Batter

				// Calculate Extras for Bowler
				if delivery.Extras != (data.Extras{}) {
					if delivery.Extras.NoBall > 0 {
						bowler.Extras += delivery.Runs.Total
						bowlerRuns += delivery.Extras.NoBall
					} else if delivery.Extras.Wides > 0 {
						bowler.Extras += delivery.Runs.Total
						bowlerRuns += delivery.Runs.Total
					}
				}

				overRunsBowler[over.Over] += bowlerRuns
				bowler.Runs += bowlerRuns

				if delivery.Extras.Wides == 0 && delivery.Extras.NoBall == 0 {
					bowler.Balls += 1
				}

				if delivery.Runs.Batter == 0 {
					bowler.Dots++
				}

				// ================  batter calculations
				batsman.Runs += delivery.Runs.Batter
				if delivery.Extras.Wides == 0 {
					batsman.Balls += 1
				}

				if delivery.Runs.Batter == 0 && delivery.Extras.Wides == 0 && delivery.Extras.NoBall == 0 {
					batsman.Dots++
				}

				// ==================== common calculations

				// Check for Wicket
				for _, wicket := range delivery.Wickets {
					if wicket.Kind != "" && wicket.PlayerOut != "" {

						var fielderName string

						if wicket.Kind != "caught and bowled" && wicket.Kind != "bowled" && wicket.Kind != "lbw" && len(wicket.Fielders) == 0 {
							fielderName = "NA"
						}
						if len(wicket.Fielders) > 0 {
							fielderName = wicket.Fielders[0].Name

							if wicket.Fielders[0].Substitute && fielderName == "" {
								fielderName = "substitute"
							}
						}

						// bowler
						if wicket.Kind != "run out" {
							bowler.Wickets++
							batsman.Out = wicket.Kind
							batsman.OutBowler = bowler.Name
							batsman.OutFielder = fielderName
						} else if wicket.Kind == "run out" {
							if batsman.Name == wicket.PlayerOut {
								batsman.Out = wicket.Kind
								batsman.OutFielder = fielderName
							} else if nonBatsman.Name == wicket.PlayerOut {
								nonBatsman.Out = wicket.Kind
								nonBatsman.OutFielder = fielderName
							}
						}
					}
				}

				// 4s/6s hit and conceded
				if delivery.Runs.Batter == 4 {
					batsman.Fours++
					bowler.FoursConceded++
				} else if delivery.Runs.Batter == 6 {
					batsman.Sixes++
					bowler.SixesConceded++
				} else if delivery.Runs.Batter == 1 {
					batsman.Singles++
				} else if delivery.Runs.Batter == 2 {
					batsman.Doubles++
				} else if delivery.Runs.Batter == 3 {
					batsman.Triples++
				}

				// ======== BIND ALL STATS
				objBatsman[delivery.Batter] = batsman
				objBatsman[delivery.NonStriker] = nonBatsman
				objBowler[delivery.Bowler] = bowler
			}

			// check maiden over
			if val, ok := overRunsBowler[over.Over]; ok && val == 0 {
				if len(over.Deliveries) > 0 {
					bowler := objBowler[over.Deliveries[0].Bowler]
					bowler.Maiden++
					objBowler[over.Deliveries[0].Bowler] = bowler
				}
			}
		}

		for name, batter := range objBatsman {
			if batter.Out == "" {
				batter.Out = "not out"
			}
			objBatsman[name] = batter
		}

		for name, bowler := range objBowler {
			bowler.Overs = fmt.Sprint(bowler.Balls/6) + "." + fmt.Sprint(bowler.Balls%6)
			objBowler[name] = bowler
		}

		// bind batter stats for the batting team on key teamID**inning
		// fmt.Println(objBatsman)
		battingKey := strconv.Itoa(battingTeamID) // + "**" + strconv.Itoa(inningID)
		bowlingKey := strconv.Itoa(bowlingTeamID) // + "**" + strconv.Itoa(inningID)
		battingStats := teamInningPlayerStats[battingKey]
		bowlingStats := teamInningPlayerStats[bowlingKey]

		// teamInningPlayerStats[bowlingTeamID] =
		battingStats = BindBattingPlayerStats(battingStats, objBatsman)
		bowlingStats = BindBowlingPlayerStats(bowlingStats, objBowler, objFielder)

		teamInningPlayerStats[battingKey] = battingStats
		teamInningPlayerStats[bowlingKey] = bowlingStats

		if inningID == 2 || inningID == 4 || (len(match.Innings) == 3 && inningID == 3) {
			// fmt.Println(teamInningPlayerStats)
			if len(match.Innings) == 3 && inningID == 3 {
				innBatTeamMap[bowlingTeamID] = inningID
			}
			data.InsertPlayerStats(matchID, seasonID, teamInningPlayerStats, allPlayerID, innBatTeamMap)

			// after second inning empty the playerstat for next inning
			teamInningPlayerStats = map[string]map[string]data.PlayerStats{}
		}
	}
}

func processMatchStats(match data.Match, fileName string) {

	var objMatchStats data.MatchStats
	cricSheetID := strings.Replace(fileName, ".json", "", -1)
	matchID := data.GetMatchID(cricSheetID)
	if matchID == 0 {
		data.InsertErrorLog(data.MATCH_NOT_FOUND, `matchID not found `+cricSheetID, fileName)
		return
	}

	// check for super over in a match
	for _, inning := range match.Innings {
		if inning.SuperOvers {
			objMatchStats.SuperOver = true
		}
	}

	for innID, inning := range match.Innings {
		innID++
		var fowArr []string
		var wicketCnt, runningScore, extras int

		if inning.SuperOvers {
			continue
		}
		for _, over := range inning.Overs {

			for _, delivery := range over.Deliveries {

				runningScore += delivery.Runs.Total

				// Check for Wicket
				for _, wicket := range delivery.Wickets {
					if wicket.Kind != "" && wicket.PlayerOut != "" {
						wicketCnt++
						fowStr := fmt.Sprint(wicketCnt, "-", runningScore, "(", wicket.PlayerOut, ")")
						fowArr = append(fowArr, fowStr)
					}
				}

				// Calculate Extras
				if delivery.Extras != (data.Extras{}) {
					if delivery.Extras.Byes > 0 {
						extras += delivery.Runs.Extras
					} else if delivery.Extras.LegByes > 0 {
						extras += delivery.Runs.Extras
					} else if delivery.Extras.NoBall > 0 {
						extras += delivery.Runs.Extras
					} else if delivery.Extras.Wides > 0 {
						extras += delivery.Runs.Extras
					}
				}
			}
		}
		objMatchStats.FOW = strings.Join(fowArr, " , ")
		objMatchStats.Score = runningScore
		objMatchStats.Extras = extras
		objMatchStats.Wickets = wicketCnt
		objMatchStats.InningsID = innID

		tempOvers := len(inning.Overs) - 1
		var lastoverBalls int
		for _, delivery := range inning.Overs[tempOvers].Deliveries {
			if delivery.Extras.Wides == 0 && delivery.Extras.NoBall == 0 {
				lastoverBalls++
			}
		}
		if lastoverBalls == 6 {
			tempOvers++
			lastoverBalls = 0
		}
		totalOvers := strconv.Itoa(tempOvers)
		if lastoverBalls > 0 {
			totalOvers = strconv.Itoa(tempOvers) + "." + strconv.Itoa(lastoverBalls)
		}
		objMatchStats.OversPlayed = totalOvers
		tempTeamID := data.GetTeamID(inning.Team, match.Info.TeamType)
		if tempTeamID == 0 {
			fmt.Println("team id not found")
		}
		objMatchStats.TeamID = tempTeamID

		data.InsertMatchStats(matchID, objMatchStats)
	}
}
