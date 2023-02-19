package server

import "github.com/devnull-twitch/nyooom-backend/pkg/jsondb"

func convertTeamsToResponse(teams []jsondb.Team, events []jsondb.RaceEvent) []teamResponse {
	teamMap := make(map[uint64]*teamResponse)
	driverMap := make(map[uint64]*driverResponse)
	teamDriverIDs := make(map[uint64][]uint64)

	for _, t := range teams {
		teamMap[t.ID] = &teamResponse{
			ID:      t.ID,
			Name:    t.Name,
			Results: make([]teamResultResponse, 0),
			Drivers: make([]driverResponse, 2),
		}

		driverIDs := make([]uint64, 0)

		for _, d := range t.Drivers {
			driverMap[d.ID] = &driverResponse{
				ID:      d.ID,
				Name:    d.Name,
				Points:  0,
				Results: make([]driverResultResponse, 0),
			}

			driverIDs = append(driverIDs, d.ID)
		}

		teamDriverIDs[t.ID] = driverIDs
	}

	for _, e := range events {
		for _, result := range e.Results {
			if e.Type == jsondb.RaceEventType || e.Type == jsondb.SprintEventType {
				teamMap[result.TeamID].Points += result.Points
			}
			teamMap[result.TeamID].Results = append(teamMap[result.TeamID].Results, teamResultResponse{
				EventName:  e.Name,
				DriverName: driverMap[result.DriverID].Name,
				Points:     result.Points,
				Position:   result.Position,
			})
			if e.Type == jsondb.RaceEventType || e.Type == jsondb.SprintEventType {
				driverMap[result.DriverID].Points += result.Points
			}
			driverMap[result.DriverID].Results = append(driverMap[result.DriverID].Results, driverResultResponse{
				EventName: e.Name,
				Points:    result.Points,
				Position:  result.Position,
			})
		}
	}

	finalArray := make([]teamResponse, len(teamMap))
	index := 0
	for teamID, teamPtr := range teamMap {
		finalArray[index] = *teamPtr

		for driverIndex, id := range teamDriverIDs[teamID] {
			finalArray[index].Drivers[driverIndex] = *driverMap[id]
		}

		index++
	}

	return finalArray
}

func convertTeamFlat(team jsondb.Team) teamResponse {
	driverList := make([]driverResponse, len(team.Drivers))
	for i, d := range team.Drivers {
		driverList[i] = driverResponse{ID: d.ID, Name: d.Name}
	}

	return teamResponse{
		ID:      team.ID,
		Name:    team.Name,
		Results: []teamResultResponse{},
		Drivers: driverList,
	}
}
