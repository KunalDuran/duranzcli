package process

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/KunalDuran/duranzcli/module/data"
)

func VenueMapper(venueName, city string) {

	if val, ok := data.MappedVenues[venueName]; ok && city == val {
		return
	}

	sqlStr := `INSERT INTO duranz_venue(venue, city) VALUES(? , ?)`

	_, err := data.SportsDb.Exec(sqlStr, venueName, city)
	if err != nil {
		panic(err.Error())
	}
	data.MappedVenues[venueName] = city

}

func TeamMapper(teams []string, teamType string) {

	for _, team := range teams {

		if _, ok := data.MappedTeams[team]; ok {
			continue
		}

		sqlStr := `INSERT INTO duranz_teams(team_name, team_type) VALUES(? , ?)`

		_, err := data.SportsDb.Exec(sqlStr, team, teamType)
		if err != nil {
			panic(err.Error())
		}
		data.MappedTeams[team] = teamType
	}
}

func PlayerMapper(players map[string]string) {

	var allTeamPlayers []string
	for _, playerID := range players {
		allTeamPlayers = append(allTeamPlayers, playerID)
	}

	allPlayerStr := strings.Join(allTeamPlayers, "','")

	// check for existing players
	sqlStr := `SELECT player_name, cricsheet_id FROM duranz_cricket_players WHERE cricsheet_id IN ('` + allPlayerStr + `')`
	rows, err := data.SportsDb.Query(sqlStr)
	if err != nil && err != sql.ErrNoRows {
		panic(err.Error())
	}

	existingPlayers := map[string]string{}
	for rows.Next() {
		var playerName, playerID string
		err = rows.Scan(
			&playerName,
			&playerID,
		)
		if err != nil {
			panic(err)
		}
		existingPlayers[playerID] = playerName
	}

	playerToInsert := map[string]string{}
	for playerName, playerID := range players {
		if _, ok := existingPlayers[playerID]; !ok {
			playerToInsert[playerName] = playerID
		}
	}

	if len(playerToInsert) > 0 {
		var valueStr []string
		valArgs := []interface{}{}
		for playerName, playerID := range playerToInsert {
			valueStr = append(valueStr, `(?, ?)`)
			valArgs = append(valArgs, playerName, playerID)
		}
		sqlStr := `INSERT INTO duranz_cricket_players(player_name, cricsheet_id) VALUES `

		sqlStr = sqlStr + strings.Join(valueStr, ",")
		_, err := data.SportsDb.Exec(sqlStr, valArgs...)
		if err != nil {
			panic(err)
		}
	}
}

func MatchMapper(match data.Match, fileName string) {

	fileName = strings.Replace(fileName, ".json", "", -1)
	var matchDates, startDate string
	leagueID := data.AllDuranzLeagues[match.Info.MatchType]

	if match.Info.Event.Name == "Indian Premier League" {
		leagueID = data.AllDuranzLeagues["ipl"]
	}
	venueID := data.GetVenueID(match.Info.Venue, match.Info.City)
	if len(match.Info.Dates) > 0 {
		startDate = match.Info.Dates[0]
		matchDates = strings.Join(match.Info.Dates, ";")
	}

	if len(match.Info.Teams) != 2 || venueID == 0 || leagueID == 0 || startDate == "" {
		fmt.Println("Error in match mapper process")
		if len(match.Info.Teams) != 2 {
			data.InsertErrorLog(data.DATABASE_ERROR, `More than 2 teams`, fileName)
		}
		if venueID == 0 {
			data.InsertErrorLog(data.DATABASE_ERROR, `Venue not found `+match.Info.Venue, fileName)
		}
		if leagueID == 0 {
			data.InsertErrorLog(data.LEAGUE_NOT_FOUND, `League not found `+match.Info.MatchType, fileName)
		}
		if startDate == "" {
			data.InsertErrorLog(data.DATETIME_ERROR, `Start date not found`, fileName)
		}
		return
	}

	home := match.Info.Teams[0]
	away := match.Info.Teams[1]
	homeTeamID := data.GetTeamID(home, match.Info.TeamType)
	awayTeamID := data.GetTeamID(away, match.Info.TeamType)

	var winningTeam, ManOfTheMatch, tossWinner int
	if home == match.Info.Outcome.Winner {
		winningTeam = homeTeamID
	} else if away == match.Info.Outcome.Winner {
		winningTeam = awayTeamID
	} else if match.Info.Outcome.Eliminator == home {
		winningTeam = homeTeamID
	} else if match.Info.Outcome.Eliminator == away {
		winningTeam = awayTeamID
	}

	if home == match.Info.Toss.Winner {
		tossWinner = homeTeamID
	} else if away == match.Info.Toss.Winner {
		tossWinner = awayTeamID
	}

	if len(match.Info.PlayerOfMatch) > 0 {
		peopleRegistery := match.Info.Register.People
		ManOfTheMatch = data.GetPlayerID(peopleRegistery[match.Info.PlayerOfMatch[0]])
	}

	tossDecision := match.Info.Toss.Decision

	if tossWinner == 0 {
		fmt.Println("Error in mapping match ", fileName)
		fmt.Println("TossWinner ", tossWinner)
		data.InsertErrorLog(data.DATABASE_ERROR, `Toss winner not found`, fileName)
		return
	}

	var resultStr string
	if match.Info.Outcome.Result == "no result" {
		resultStr = "no result"
	} else if match.Info.Outcome.Result == "tie" {
		resultStr = "tie"
	} else if match.Info.Outcome.Result == "draw" {
		resultStr = "draw"
	}

	if resultStr != "no result" && resultStr != "tie" && resultStr != "draw" {
		resultStr = match.Info.Outcome.Winner + " Won by "
		if match.Info.Outcome.By.Runs > 0 {
			resultStr += strconv.Itoa(match.Info.Outcome.By.Runs) + " Runs"
		} else if match.Info.Outcome.By.Wickets > 0 {
			resultStr += strconv.Itoa(match.Info.Outcome.By.Wickets) + " Wickets"
		}
	}
	if match.Info.Outcome.Method == "D/L" {
		resultStr += " (D/L Method)"
	}

	matchReferees := strings.Join(match.Info.MatchReferees, ";")
	reserveUmpires := strings.Join(match.Info.ReserveUmpires, ";")
	tvUmpires := strings.Join(match.Info.TvUmpires, ";")
	umpires := strings.Join(match.Info.Umpires, ";")

	sqlStr := `INSERT INTO duranz_cricket_matches(league_id, home_team_id, away_team_id, home_team_name, away_team_name, venue_id, match_date, match_date_multi, cricsheet_file_name,
				result, man_of_the_match, toss_winner, toss_decision, winning_team, gender, 
				match_refrees, reserve_umpires, tv_umpires, umpires) 
				VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := data.SportsDb.Exec(sqlStr, leagueID, homeTeamID, awayTeamID, home, away, venueID, startDate,
		matchDates, fileName, resultStr, ManOfTheMatch, tossWinner, tossDecision, winningTeam,
		match.Info.Gender, matchReferees, reserveUmpires, tvUmpires, umpires)
	if err != nil {
		fmt.Println(err.Error())
	}
}
