package data

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
)

// SportsDb is the pointer to the iSports database resource.
var SportsDb *sql.DB

const PLAYER_NOT_FOUND = "1"
const TEAM_NOT_FOUND = "2"
const MATCH_NOT_FOUND = "3"
const VENUE_NOT_FOUND = "4"
const DATABASE_ERROR = "5"
const CRICSHEET_FILE_ERROR = "6"
const DATETIME_ERROR = "7"
const LEAGUE_NOT_FOUND = "8"

var AllDuranzLeagues = map[string]int{
	"ODI":                   1,
	"Test":                  2,
	"T20":                   3,
	"ipl":                   4,
	"Indian Premier League": 4,
}

var GamePath = map[string]string{
	"odi":  "odis_json",
	"test": "tests_json",
	"t20":  "t20s_json",
	"ipl":  "ipl_json",
}

// cache simulation
var MappedTeams = map[string]string{}
var MappedVenues = map[string]string{}

// InitDuranzDB :
func InitDuranzDB(host, port, user, password string) (sportsDb *sql.DB, err error) {
	SportsDb, err = sql.Open("mysql", user+":"+password+"@tcp("+host+":"+port+")/duranz")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	if err = SportsDb.Ping(); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return SportsDb, nil
}

func InsertMatchStats(matchID int, objMatchStats MatchStats) {

	sqlStr := `INSERT INTO duranz_match_stats(
	match_id ,	
	team_id, 
	fall_of_wickets, 
	extras, 
	score, 
	super_over,  
	wickets, 
	overs_played, 
	innings
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := SportsDb.Exec(sqlStr,
		matchID,
		objMatchStats.TeamID,
		objMatchStats.FOW,
		objMatchStats.Extras,
		objMatchStats.Score,
		objMatchStats.SuperOver,
		objMatchStats.Wickets,
		objMatchStats.OversPlayed,
		objMatchStats.InningsID,
	)
	if err != nil {
		panic(err)
	}
}

func GetVenueID(venueName, city string) int {
	var venueID int
	sqlStr := `SELECT venue_id FROM duranz_venue WHERE venue = ? AND city = ?`

	err := SportsDb.QueryRow(sqlStr, venueName, city).Scan(&venueID)
	if err != nil && err != sql.ErrNoRows {
		panic(err)
	}
	return venueID
}

func GetTeamID(teamName, team_type string) int {
	var teamID int
	sqlStr := `SELECT team_id FROM duranz_teams WHERE team_name = ? AND team_type = ?`

	err := SportsDb.QueryRow(sqlStr, teamName, team_type).Scan(&teamID)
	if err != nil && err != sql.ErrNoRows {
		panic(err)
	}
	return teamID
}

func GetPlayerID(cricsheetID string) int {
	var playerID int
	sqlStr := `SELECT player_id FROM duranz_cricket_players WHERE cricsheet_id = ?`

	err := SportsDb.QueryRow(sqlStr, cricsheetID).Scan(&playerID)
	if err != nil && err != sql.ErrNoRows {
		panic(err)
	}
	return playerID
}

func CheckTable(tableName, whereClause string) (int, error) {
	var countVal int
	// sqlCheck := `SELECT COUNT(1) FROM ? WHERE ?`
	sqlCheck := `SELECT COUNT(1) FROM ` + tableName + ` WHERE ` + whereClause
	// err := data.SportsDb.QueryRow(sqlCheck, tableName, whereClause).Scan(&countVal)
	// fmt.Println(sqlCheck)
	err := SportsDb.QueryRow(sqlCheck).Scan(&countVal)
	if err != nil {
		return countVal, err
	}

	return countVal, nil
}

func GetMatchID(cricsheetID string) int {
	var matchID int
	sqlStr := `SELECT match_id FROM duranz_cricket_matches WHERE cricsheet_file_name = ?`

	err := SportsDb.QueryRow(sqlStr, cricsheetID).Scan(&matchID)
	if err != nil && err != sql.ErrNoRows {
		panic(err)
	}
	return matchID
}

//InsertErrorLog : Insert All Error log record according to alert ID, file name
func InsertErrorLog(alertid, errormsg, fileName string) {
	sqlStr := `INSERT INTO duranz_errorlog (alert_id, error_msg, file_name) VALUES (?, ?, ?)`

	_, err := SportsDb.Exec(sqlStr, alertid, errormsg, fileName)
	if err != nil {
		panic(err)
	}

}

func InsertMappingInfo(fileName string, mappingInfo MappingInfo) {
	sqlStr := `INSERT INTO duranz_file_mappings (file_name, league_id, teams, players, venue, matches, match_stats, player_stats) 
				VALUES (?, ?, ?, ?, ?, ?, ?, ?) 
				ON DUPLICATE KEY UPDATE file_name=values(file_name), teams=values(teams), players=values(players), venue=values(venue), matches=values(matches), match_stats=values(match_stats), player_stats = values(player_stats)`

	_, err := SportsDb.Exec(sqlStr, fileName, mappingInfo.LeagueID, mappingInfo.Teams, mappingInfo.Players, mappingInfo.Venue, mappingInfo.Match, mappingInfo.MatchStats, mappingInfo.PlayerStats)
	if err != nil {
		panic(err)
	}
}

func DeleteAllTableData() {

	tables := []string{
		"duranz_errorlog",
		"duranz_teams",
		"duranz_venue",
		"duranz_cricket_players",
		"duranz_cricket_matches",
		"duranz_match_stats",
		"duranz_player_match_stats",
		"duranz_file_mappings",
	}

	for _, table := range tables {
		sqlStr := `DELETE FROM ` + table
		_, err := SportsDb.Exec(sqlStr)
		if err != nil {
			panic(err)
		}

		sqlStr2 := `ALTER TABLE ` + table + ` AUTO_INCREMENT = 1`
		_, err = SportsDb.Exec(sqlStr2)
		if err != nil {
			panic(err)
		}
	}
}

func GetMappingDetails() map[string]MappingInfo {

	var objMappingInfo = map[string]MappingInfo{}

	sqlStr := `SELECT file_name, league_id, teams, players, venue, matches, match_stats, player_stats FROM duranz_file_mappings`
	rows, err := SportsDb.Query(sqlStr)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var mappingInfo MappingInfo
		var fileName string
		err := rows.Scan(
			&fileName,
			&mappingInfo.LeagueID,
			&mappingInfo.Teams,
			&mappingInfo.Players,
			&mappingInfo.Venue,
			&mappingInfo.Match,
			&mappingInfo.MatchStats,
			&mappingInfo.PlayerStats,
		)

		if err != nil {
			panic(err)
		}

		objMappingInfo[fileName] = mappingInfo
	}

	return objMappingInfo
}

func InsertPlayerStats(matchID, seasonID int, teamInningPlayerStats map[string]map[string]PlayerStats, allPlayerID map[string]int, innBatTeamMap map[int]int) {

	var valueStr []string
	valArgs := []interface{}{}
	sqlStr := `INSERT INTO duranz_player_match_stats(
	match_id ,
	season_id ,
	innings_id ,
	team_id ,
	player_id ,
	batting_order ,
	runs_scored ,
	balls_faced ,
	dot_balls_played ,
	singles ,
	doubles ,
	triples ,
	fours_hit ,
	sixes_hit ,
	out_type ,
	out_bowler ,
	out_fielder ,
	is_batted ,
	overs_bowled ,
	runs_conceded ,
	balls_bowled ,
	dots_bowled ,
	wickets_taken ,
	fours_conceded ,
	sixes_conceded ,
	extras_conceded ,
	maiden_over ,
	bowling_order,
	run_out ,
	catches ,
	stumpings
) VALUES `
	for teamID, allPlayerStats := range teamInningPlayerStats {
		for _, playerStats := range allPlayerStats {
			tempPlayerID := allPlayerID[playerStats.Name]
			if tempPlayerID == 0 {
				fmt.Println(playerStats.Name)
				// continue
			}
			teamIDint, _ := strconv.Atoi(teamID)
			innID := innBatTeamMap[teamIDint]
			playerStats.PlayerID = tempPlayerID
			valueStr = append(valueStr, `(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)

			valArgs = append(valArgs,
				matchID,
				seasonID,
				innID,
				teamID,
				playerStats.PlayerID,
				playerStats.BattingOrder,
				playerStats.RunsScored,
				playerStats.BallsPlayed,
				playerStats.DotsPlayed,
				playerStats.Singles,
				playerStats.Doubles,
				playerStats.Triples,
				playerStats.FoursHit,
				playerStats.SixesHit,
				playerStats.OutType,
				playerStats.OutBowler,
				playerStats.OutFielder,
				playerStats.IsBatted,

				playerStats.OversBowled,
				playerStats.RunsConceded,
				playerStats.BallsBowled,
				playerStats.DotsBowled,
				playerStats.WicketsTaken,
				playerStats.FoursConceded,
				playerStats.SixesConceded,
				playerStats.ExtrasConceded,
				playerStats.MaidenOvers,
				playerStats.BowlingOrder,

				playerStats.RunOuts,
				playerStats.Catches,
				playerStats.Stumpings,
			)

		}
	}

	sqlStr = sqlStr + strings.Join(valueStr, ",")
	_, err := SportsDb.Exec(sqlStr, valArgs...)
	if err != nil {
		panic(err)
	}
}

func PseudoCacheLayer(teamType string) {
	// load Teams
	if teamType != "ipl" {
		teamType = "international"
	}

	sqlStr := `SELECT team_name FROM duranz_teams WHERE team_type = ?`

	rows, err := SportsDb.Query(sqlStr, teamType)
	if err != nil && err != sql.ErrNoRows {
		panic(err)
	}
	for rows.Next() {
		var teamName string
		err = rows.Scan(&teamName)
		if err != nil {
			panic(err)
		}
		MappedTeams[teamName] = teamType
	}

	// load Venues
	sqlStr = `SELECT venue, city FROM duranz_venue`

	rows, err = SportsDb.Query(sqlStr)
	if err != nil && err != sql.ErrNoRows {
		panic(err)
	}
	for rows.Next() {
		var venue, city string
		err = rows.Scan(&venue, &city)
		if err != nil {
			panic(err)
		}
		MappedVenues[venue] = city
	}
}
