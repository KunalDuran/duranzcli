package data

import "encoding/json"

type Match struct {
	Meta struct {
		DataVersion string `json:"data_version"`
		Created     string `json:"created"`
		Revision    int    `json:"revision"`
	} `json:"meta"`
	Info struct {
		BallsPerOver int      `json:"balls_per_over"`
		City         string   `json:"city"`
		Dates        []string `json:"dates"`
		Event        struct {
			Name        string `json:"name"`
			MatchNumber int    `json:"match_number"`
		} `json:"event"`
		Gender          string `json:"gender"`
		MatchType       string `json:"match_type"`
		MatchTypeNumber int    `json:"match_type_number"`
		Officials       `json:"officials"`
		Outcome         struct {
			By struct {
				Runs    int `json:"runs"`
				Wickets int `json:"wickets"`
			} `json:"by"`
			Winner     string `json:"winner"`
			Eliminator string `json:"eliminator"`
			Result     string `json:"result"`
			Method     string `json:"method"`
		} `json:"outcome"`
		Overs         int                 `json:"overs"`
		PlayerOfMatch []string            `json:"player_of_match"`
		Players       map[string][]string `json:"players"`
		Register      Registry            `json:"registry"`
		Season        json.RawMessage     `json:"season"`
		TeamType      string              `json:"team_type"`
		Teams         []string            `json:"teams"`
		Toss          struct {
			Decision string `json:"decision"`
			Winner   string `json:"winner"`
		} `json:"toss"`
		Venue string `json:"venue"`
	} `json:"info"`
	Innings []struct {
		Team       string        `json:"team"`
		Overs      []OverDetails `json:"overs"`
		SuperOvers bool          `json:"super_over"`
		Powerplays []struct {
			From float64 `json:"from"`
			To   float64 `json:"to"`
			Type string  `json:"type"`
		} `json:"powerplays"`
		Target struct {
			Overs json.RawMessage `json:"overs"`
			Runs  int             `json:"runs"`
		} `json:"target,omitempty"`
	} `json:"innings"`
}

type Registry struct {
	People map[string]string `json:"people"`
}

// SCORE CARD STRUCTS
type ScoreCard struct {
	Innings []Innings `json:"innings"`
	Result  string    `json:"result"`
}

type Innings struct {
	InningID      int            `json:"innings_id"`
	InningDetail  string         `json:"innings_detail"`
	Bowling       []BowlingStats `json:"bowling"`
	Batting       []BattingStats `json:"batting"`
	Extras        `json:"extras"`
	FallOfWickets string `json:"fall_of_wickets"`
}

type BowlingStats struct {
	BowlingOrder  int     `json:"bowling_order"`
	Name          string  `json:"name"`
	Overs         string  `json:"overs"`
	Maiden        int     `json:"maiden"`
	Runs          int     `json:"runs"`
	Wickets       int     `json:"wickets"`
	Economy       float64 `json:"economy"`
	Balls         int     `json:"-"`
	Dots          int     `json:"-"`
	IsBowled      bool    `json:"-"`
	FoursConceded int     `json:"-"`
	SixesConceded int     `json:"-"`
	Extras        int     `json:"-"`
}

type BattingStats struct {
	BattingOrder int     `json:"batting_order"`
	Name         string  `json:"name"`
	Runs         int     `json:"runs"`
	Balls        int     `json:"balls"`
	Fours        int     `json:"fours"`
	Sixes        int     `json:"sixes"`
	StrikeRate   float64 `json:"strike_rate"`
	Out          string  `json:"out"`
	OutBowler    string  `json:"-"`
	OutFielder   string  `json:"-"`
	Singles      int     `json:"-"`
	Doubles      int     `json:"-"`
	Triples      int     `json:"-"`
	Dots         int     `json:"-"`
	NotOut       bool    `json:"-"`
	IsBatted     bool    `json:"-"`
}

type Extras struct {
	Wides   int `json:"wides"`
	NoBall  int `json:"noballs"`
	Byes    int `json:"byes"`
	LegByes int `json:"legbyes"`
	Total   int `json:"total"`
}

type PlayerStats struct {
	PlayerID int
	Name     string

	// bowling stats
	BowlingOrder   int
	OversBowled    string
	MaidenOvers    int
	RunsConceded   int
	WicketsTaken   int
	Economy        float64
	BallsBowled    int
	DotsBowled     int
	FoursConceded  int
	SixesConceded  int
	ExtrasConceded int

	// batting stats
	BattingOrder int
	RunsScored   int
	BallsPlayed  int
	Singles      int
	Doubles      int
	Triples      int
	FoursHit     int
	SixesHit     int
	StrikeRate   float64
	OutType      string
	OutBowler    string
	OutFielder   string
	DotsPlayed   int
	NotOut       bool
	IsBatted     bool

	// fielding stats
	RunOuts   int
	Catches   int
	Stumpings int
}

type FieldingStats struct {
	Name      string
	RunOuts   int
	Catches   int
	Stumpings int
}

type TeamStats struct {
}

type MatchStats struct {
	MatchID     int
	TeamID      int
	InningsID   int
	Captain     string
	FOW         string
	Extras      int
	Wickets     int
	OversPlayed string
	Score       int
	SuperOver   bool
}

type Venue struct {
}

type PlayerStatsBind struct {
	Batting  map[string]BattingStats
	Bowling  map[string]BowlingStats
	Fielding map[string]FieldingStats
}

type Officials struct {
	MatchReferees  []string `json:"match_referees"`
	ReserveUmpires []string `json:"reserve_umpires"`
	TvUmpires      []string `json:"tv_umpires"`
	Umpires        []string `json:"umpires"`
}

type OverDetails struct {
	Over       int `json:"over"`
	Deliveries []struct {
		Batter     string `json:"batter"`
		Bowler     string `json:"bowler"`
		Extras     `json:"extras"`
		NonStriker string `json:"non_striker"`
		Runs       struct {
			Batter int `json:"batter"`
			Extras int `json:"extras"`
			Total  int `json:"total"`
		} `json:"runs"`
		Wickets []struct {
			Kind      string `json:"kind"`
			PlayerOut string `json:"player_out"`
			Fielders  []struct {
				Name       string `json:"name"`
				Substitute bool   `json:"substitute"`
			} `json:"fielders"`
		}
	} `json:"deliveries"`
}

type MappingInfo struct {
	FileName    string
	LeagueID    int
	Teams       bool
	Players     bool
	Venue       bool
	Match       bool
	MatchStats  bool
	PlayerStats bool
}
