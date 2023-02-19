package server

type teamResultResponse struct {
	EventName  string `json:"event_name"`
	DriverName string `json:"driver_name"`
	Points     uint64 `json:"points"`
	Position   uint64 `json:"position"`
}

type driverResultResponse struct {
	EventName string `json:"event_name"`
	Points    uint64 `json:"points"`
	Position  uint64 `json:"position"`
}

type driverResponse struct {
	ID      uint64                 `json:"id"`
	Name    string                 `json:"name"`
	Points  uint64                 `json:"points"`
	Results []driverResultResponse `json:"results"`
}

type teamResponse struct {
	ID      uint64               `json:"id"`
	Name    string               `json:"name"`
	Points  uint64               `json:"points"`
	Results []teamResultResponse `json:"results"`
	Drivers []driverResponse     `json:"drivers"`
}

type eventGridResponse struct {
	DriverID   uint64 `json:"driver_id"`
	DriverName string `json:"driver_name"`
	TeamName   string `json:"team_name"`
	Position   uint64 `json:"position"`
}

type eventResultResponse struct {
	DriverName string `json:"driver_name"`
	DriverID   uint64 `json:"driver_id"`
	TeamName   string `json:"team_name"`
	Position   uint64 `json:"position"`
	Points     uint64 `json:"points"`
}

type eventResponse struct {
	ID           uint64                `json:"id"`
	Name         string                `json:"name"`
	Type         string                `json:"type"`
	UnixDate     int64                 `json:"race_date_unix"`
	StartingGrid []eventGridResponse   `json:"starting_grid"`
	Results      []eventResultResponse `json:"results"`
}
