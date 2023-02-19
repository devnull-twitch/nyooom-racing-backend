package jsondb

type TeamSchema struct {
	Teams        []Team `json:"teams"`
	NextTeamID   uint64 `json:"next_team_id"`
	NextDriverID uint64 `json:"next_driver_id"`
}

type EventSchema struct {
	Events      []RaceEvent `json:"events"`
	NextEventID uint64      `json:"next_event_id"`
}
