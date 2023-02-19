package jsondb

type Driver struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

type EventType uint

const (
	RaceEventType EventType = iota + 1
	SprintEventType
	PreSeason
	PreSeasonSprintType
)

func (e EventType) Name() string {
	switch e {
	case RaceEventType:
		return "Race"
	case SprintEventType:
		return "Sprint"
	case PreSeason:
		return "Pre-Season race"
	case PreSeasonSprintType:
		return "Pre-Season sprint"
	default:
		return "Unknown"
	}
}

type RaceEvent struct {
	ID           uint64         `json:"id"`
	Name         string         `json:"name"`
	Date         int64          `json:"date_unix"`
	Type         EventType      `json:"race_type"`
	StartingGrid []RacePosition `json:"starting"`
	Results      []RacePosition `json:"results"`
}

type RacePosition struct {
	Position uint64 `json:"position"`
	Points   uint64 `json:"points"`
	DriverID uint64 `json:"driver_id"`
	TeamID   uint64 `json:"team_id"`
}

type Team struct {
	ID      uint64   `json:"id"`
	Name    string   `json:"name"`
	Drivers []Driver `json:"drivers"`
}
