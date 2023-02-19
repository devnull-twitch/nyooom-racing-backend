package server

import "github.com/devnull-twitch/nyooom-backend/pkg/jsondb"

func buildNameMaps(repo jsondb.JsonDatabase) (
	teamNameMap map[uint64]string,
	driverNameMap map[uint64]string,
	err error,
) {
	var teams []jsondb.Team
	teams, err = repo.ListTeams()
	if err != nil {
		return
	}

	teamNameMap = make(map[uint64]string)
	driverNameMap = make(map[uint64]string)

	for _, t := range teams {
		teamNameMap[t.ID] = t.Name
		for _, d := range t.Drivers {
			driverNameMap[d.ID] = d.Name
		}
	}

	return
}

func convertEventsToResponse(
	events []jsondb.RaceEvent,
	teamNameMap map[uint64]string,
	driverNameMap map[uint64]string,
) []eventResponse {
	finalResp := make([]eventResponse, 0, len(events))
	for _, event := range events {
		finalResp = append(finalResp, convertEventToResponse(event, teamNameMap, driverNameMap))
	}

	return finalResp
}

func convertEventToResponse(
	event jsondb.RaceEvent,
	teamNameMap map[uint64]string,
	driverNameMap map[uint64]string,
) eventResponse {
	grid := make([]eventGridResponse, 0)
	for _, gridPos := range event.StartingGrid {
		grid = append(grid, eventGridResponse{
			DriverID:   gridPos.DriverID,
			DriverName: driverNameMap[gridPos.DriverID],
			TeamName:   teamNameMap[gridPos.TeamID],
			Position:   gridPos.Position,
		})
	}

	result := make([]eventResultResponse, 0)
	for _, eventRes := range event.Results {
		result = append(result, eventResultResponse{
			DriverName: driverNameMap[eventRes.DriverID],
			DriverID:   eventRes.DriverID,
			TeamName:   teamNameMap[eventRes.TeamID],
			Position:   eventRes.Position,
			Points:     eventRes.Points,
		})
	}

	return eventResponse{
		ID:           event.ID,
		Type:         event.Type.Name(),
		UnixDate:     event.Date,
		Name:         event.Name,
		StartingGrid: grid,
		Results:      result,
	}
}
