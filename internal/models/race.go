package models

// ---------- Jolpica Ergast API response models ----------

// APIResponse is the top-level wrapper from the Ergast-compatible API.
type APIResponse struct {
	MRData MRData `json:"MRData"`
}

type MRData struct {
	Series        string         `json:"series"`
	URL           string         `json:"url"`
	Limit         string         `json:"limit"`
	Offset        string         `json:"offset"`
	Total         string         `json:"total"`
	RaceTable     *RaceTable     `json:"RaceTable,omitempty"`
	StandingsTable *StandingsTable `json:"StandingsTable,omitempty"`
}

// ---------- Race / Results ----------

type RaceTable struct {
	Season string `json:"season"`
	Round  string `json:"round,omitempty"`
	Races  []Race `json:"Races"`
}

type Race struct {
	Season   string   `json:"season"`
	Round    string   `json:"round"`
	URL      string   `json:"url"`
	RaceName string   `json:"raceName"`
	Circuit  Circuit  `json:"Circuit"`
	Date     string   `json:"date"`
	Time     string   `json:"time,omitempty"`
	Results  []Result `json:"Results,omitempty"`
}

type Circuit struct {
	CircuitID   string   `json:"circuitId"`
	URL         string   `json:"url"`
	CircuitName string   `json:"circuitName"`
	Location    Location `json:"Location"`
}

type Location struct {
	Lat      string `json:"lat"`
	Long     string `json:"long"`
	Locality string `json:"locality"`
	Country  string `json:"country"`
}

type Result struct {
	Number       string      `json:"number"`
	Position     string      `json:"position"`
	PositionText string      `json:"positionText"`
	Points       string      `json:"points"`
	Driver       Driver      `json:"Driver"`
	Constructor  Constructor `json:"Constructor"`
	Grid         string      `json:"grid"`
	Laps         string      `json:"laps"`
	Status       string      `json:"status"`
	Time         *ResultTime `json:"Time,omitempty"`
	FastestLap   *FastestLap `json:"FastestLap,omitempty"`
}

type ResultTime struct {
	Millis string `json:"millis"`
	Time   string `json:"time"`
}

type FastestLap struct {
	Rank    string      `json:"rank"`
	Lap     string      `json:"lap"`
	Time    *ResultTime `json:"Time,omitempty"`
	AverageSpeed *AverageSpeed `json:"AverageSpeed,omitempty"`
}

type AverageSpeed struct {
	Units string `json:"units"`
	Speed string `json:"speed"`
}

// ---------- Drivers & Constructors ----------

type Driver struct {
	DriverID        string `json:"driverId"`
	PermanentNumber string `json:"permanentNumber,omitempty"`
	Code            string `json:"code,omitempty"`
	URL             string `json:"url"`
	GivenName       string `json:"givenName"`
	FamilyName      string `json:"familyName"`
	DateOfBirth     string `json:"dateOfBirth"`
	Nationality     string `json:"nationality"`
}

type Constructor struct {
	ConstructorID string `json:"constructorId"`
	URL           string `json:"url"`
	Name          string `json:"name"`
	Nationality   string `json:"nationality"`
}

// ---------- Standings ----------

type StandingsTable struct {
	Season          string           `json:"season"`
	StandingsLists  []StandingsList  `json:"StandingsLists"`
}

type StandingsList struct {
	Season              string                `json:"season"`
	Round               string                `json:"round"`
	DriverStandings     []DriverStanding      `json:"DriverStandings,omitempty"`
	ConstructorStandings []ConstructorStanding `json:"ConstructorStandings,omitempty"`
}

type DriverStanding struct {
	Position     string        `json:"position"`
	PositionText string        `json:"positionText"`
	Points       string        `json:"points"`
	Wins         string        `json:"wins"`
	Driver       Driver        `json:"Driver"`
	Constructors []Constructor `json:"Constructors"`
}

type ConstructorStanding struct {
	Position     string      `json:"position"`
	PositionText string      `json:"positionText"`
	Points       string      `json:"points"`
	Wins         string      `json:"wins"`
	Constructor  Constructor `json:"Constructor"`
}

// ---------- DynamoDB Cache ----------

// CacheItem represents a cached API response stored in DynamoDB.
type CacheItem struct {
	CacheKey  string `dynamodbav:"cache_key"`
	Data      string `dynamodbav:"data"`
	TTL       int64  `dynamodbav:"ttl"`
	FetchedAt int64  `dynamodbav:"fetched_at"`
}
