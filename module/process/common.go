package process

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/KunalDuran/duranzcli/module/data"
)

func ListFiles(path string) []string {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println(err)
		data.InsertErrorLog(data.CRICSHEET_FILE_ERROR, `Error in ListFiles : `, path)
		return []string{}
	}

	var fileList []string
	for _, f := range files {
		if strings.Contains(f.Name(), "json") {
			fileList = append(fileList, f.Name())
		}
	}
	return fileList
}

func BindBattingPlayerStats(objPlayers map[string]data.PlayerStats, objBatsman map[string]data.BattingStats) map[string]data.PlayerStats {

	if len(objPlayers) == 0 {
		objPlayers = map[string]data.PlayerStats{}
	}
	for batsmanName, batsman := range objBatsman {
		player := objPlayers[batsmanName]
		player.Name = batsman.Name
		player.BattingOrder = batsman.BattingOrder
		player.RunsScored = batsman.Runs
		player.BallsPlayed = batsman.Balls
		player.Singles = batsman.Singles
		player.Doubles = batsman.Doubles
		player.Triples = batsman.Triples
		player.FoursHit = batsman.Fours
		player.SixesHit = batsman.Sixes
		player.StrikeRate = batsman.StrikeRate
		player.OutType = batsman.Out
		player.OutBowler = batsman.OutBowler
		player.OutFielder = batsman.OutFielder
		player.DotsPlayed = batsman.Dots
		player.NotOut = batsman.NotOut
		player.IsBatted = batsman.IsBatted
		objPlayers[batsmanName] = player
	}
	return objPlayers
}

func BindBowlingPlayerStats(objPlayers map[string]data.PlayerStats, objBowler map[string]data.BowlingStats, objFielder map[string]data.FieldingStats) map[string]data.PlayerStats {
	if len(objPlayers) == 0 {
		objPlayers = map[string]data.PlayerStats{}
	}
	for bowlerName, bowler := range objBowler {
		player := objPlayers[bowlerName]
		player.Name = bowler.Name
		player.BowlingOrder = bowler.BowlingOrder
		player.OversBowled = bowler.Overs
		player.MaidenOvers = bowler.Maiden
		player.RunsConceded = bowler.Runs
		player.WicketsTaken = bowler.Wickets
		player.Economy = bowler.Economy
		player.BallsBowled = bowler.Balls
		player.DotsBowled = bowler.Dots
		player.FoursConceded = bowler.FoursConceded
		player.SixesConceded = bowler.SixesConceded
		player.ExtrasConceded = bowler.Extras
		objPlayers[bowlerName] = player
	}

	for fielderName, fielder := range objFielder {
		player := objPlayers[fielderName]
		player.Name = fielder.Name
		player.Catches = fielder.Catches
		player.RunOuts = fielder.RunOuts
		player.Stumpings = fielder.Stumpings
		objPlayers[fielderName] = player

	}

	return objPlayers
}

func GetNewFiles(allFiles []string) []string {
	objMappingDetails := data.GetMappingDetails()

	var newFiles []string

	for _, file := range allFiles {
		if _, exist := objMappingDetails[file]; !exist {
			newFiles = append(newFiles, file)
		}
	}

	return newFiles
}
