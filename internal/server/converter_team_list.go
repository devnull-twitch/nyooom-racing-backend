package server

import (
	"time"

	"github.com/devnull-twitch/nyooom-backend/pkg/jsondb"
)

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

	latestEventIDs := getLatestEventIDs(events)

	for _, e := range events {
		for _, result := range e.Results {
			var prevPoints uint64 = 0
			if !IDisInList(latestEventIDs, e.ID) {
				prevPoints = result.Points
			}

			if e.Type == jsondb.RaceEventType || e.Type == jsondb.SprintEventType {
				teamMap[result.TeamID].Points += result.Points
				teamMap[result.TeamID].PrevPoints += prevPoints
			} else {
				teamMap[result.TeamID].PreSeasonPoints += result.Points
				teamMap[result.TeamID].PrevPreSeasonPoints += prevPoints
			}
			teamMap[result.TeamID].Results = append(teamMap[result.TeamID].Results, teamResultResponse{
				EventName:  e.Name,
				DriverName: driverMap[result.DriverID].Name,
				Points:     result.Points,
				Position:   result.Position,
			})
			if e.Type == jsondb.RaceEventType || e.Type == jsondb.SprintEventType {
				driverMap[result.DriverID].Points += result.Points
				driverMap[result.DriverID].PrevPoints += prevPoints
			} else {
				driverMap[result.DriverID].PreSeasonPoints += result.Points
				driverMap[result.DriverID].PrevPreSeasonPoints += prevPoints
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

func IDisInList(ids []uint64, checkID uint64) bool {
	for _, id := range ids {
		if checkID == id {
			return true
		}
	}

	return false
}

func getLatestEventIDs(events []jsondb.RaceEvent) []uint64 {
	IDs := make([]uint64, 0)
	var latestTS time.Time

	for _, event := range events {
		orgTime := time.Unix(event.Date, 0)
		checkTS := time.Date(
			orgTime.Year(),
			orgTime.Month(),
			orgTime.Day(),
			0, 0, 0, 0, time.Local,
		)
		if checkTS.After(latestTS) {
			IDs = make([]uint64, 1)
			IDs[0] = event.ID
			latestTS = checkTS
		} else if latestTS.Sub(checkTS) < time.Minute {
			IDs = append(IDs, event.ID)
		}
	}

	return IDs
}
