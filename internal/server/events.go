package server

import (
	"net/http"
	"strconv"

	"github.com/devnull-twitch/nyooom-backend/pkg/jsondb"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetEventsHandler(repo jsondb.JsonDatabase) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		events, err := repo.ListEvents()
		if err != nil {
			logrus.WithError(err).Warn("unable to read events")
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		teamNameMap, driverNameMap, err := buildNameMaps(repo)
		if err != nil {
			logrus.WithError(err).Warn("unable to generate name maps")
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		eventResp := convertEventsToResponse(events, teamNameMap, driverNameMap)

		ctx.JSON(http.StatusOK, eventResp)
	}
}

func GetEventHandler(repo jsondb.JsonDatabase) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("race_id"))
		if err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		teamNameMap, driverNameMap, err := buildNameMaps(repo)
		if err != nil {
			logrus.WithError(err).Warn("unable to generate name maps")
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		event, err := repo.GetEvent(uint64(id))
		if err != nil {
			logrus.WithError(err).Warn("unable to load single event")
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		eventResp := convertEventToResponse(*event, teamNameMap, driverNameMap)
		ctx.JSON(http.StatusOK, eventResp)
	}
}

type raceEventRequest struct {
	Name         string           `json:"name"`
	Date         int64            `json:"race_date_unix"`
	Type         jsondb.EventType `json:"type"`
	StartingGrid []uint64         `json:"starting_grid"`
	Results      []uint64         `json:"results"`
}

func CreateRaceEventHandler(repo jsondb.JsonDatabase) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userInput := &raceEventRequest{}

		if err := ctx.BindJSON(userInput); err != nil {
			logrus.WithError(err).Warn("unable to get user input for new event")
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		newRaceEvent := &jsondb.RaceEvent{
			Name:         userInput.Name,
			Date:         userInput.Date,
			Type:         userInput.Type,
			StartingGrid: make([]jsondb.RacePosition, 0),
			Results:      make([]jsondb.RacePosition, 0),
		}

		teams, err := repo.ListTeams()
		if err != nil {
			logrus.WithError(err).Warn("unable to read teams for adding event")
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		driverToTeamMap := make(map[uint64]uint64)
		for _, t := range teams {
			for _, d := range t.Drivers {
				driverToTeamMap[d.ID] = t.ID
			}
		}

		// overwrite team IDs based on driver ID to make user input easier
		for index, driverID := range userInput.StartingGrid {
			newRaceEvent.StartingGrid = append(newRaceEvent.StartingGrid, jsondb.RacePosition{
				Position: uint64(index + 1),
				DriverID: driverID,
				TeamID:   driverToTeamMap[driverID],
			})
		}

		newRaceEvent.Results = buildResults(userInput.Results, newRaceEvent.Type, driverToTeamMap)

		if err := repo.AddEvent(newRaceEvent); err != nil {
			logrus.WithError(err).Warn("unable to add race event")
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		ctx.Status(http.StatusCreated)
	}
}

func UpdateRaceEventHandler(repo jsondb.JsonDatabase) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userInput := &raceEventRequest{}

		if err := ctx.BindJSON(userInput); err != nil {
			logrus.WithError(err).Warn("unable to get user input for new event")
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		raceID, err := strconv.Atoi(ctx.Param("race_id"))
		if err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		newRaceEvent := &jsondb.RaceEvent{
			ID:           uint64(raceID),
			Name:         userInput.Name,
			Date:         userInput.Date,
			Type:         userInput.Type,
			StartingGrid: make([]jsondb.RacePosition, 0),
			Results:      make([]jsondb.RacePosition, 0),
		}

		teams, err := repo.ListTeams()
		if err != nil {
			logrus.WithError(err).Warn("unable to read teams for adding event")
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		driverToTeamMap := make(map[uint64]uint64)
		for _, t := range teams {
			for _, d := range t.Drivers {
				driverToTeamMap[d.ID] = t.ID
			}
		}

		// overwrite team IDs based on driver ID to make user input easier
		for index, driverID := range userInput.StartingGrid {
			newRaceEvent.StartingGrid = append(newRaceEvent.StartingGrid, jsondb.RacePosition{
				Position: uint64(index + 1),
				DriverID: driverID,
				TeamID:   driverToTeamMap[driverID],
			})
		}

		newRaceEvent.Results = buildResults(userInput.Results, newRaceEvent.Type, driverToTeamMap)

		if err := repo.UpdateEvent(newRaceEvent); err != nil {
			logrus.WithError(err).Warn("unable to update event")
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		ctx.Status(http.StatusOK)
	}
}

func DeleteRaceEventHandler(repo jsondb.JsonDatabase) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		raceID, err := strconv.Atoi(ctx.Param("race_id"))
		if err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		if err := repo.DeleteEvent(uint64(raceID)); err != nil {
			logrus.WithError(err).Warn("unable to delete event")
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		ctx.Status(http.StatusOK)
	}
}

func buildResults(input []uint64, eventType jsondb.EventType, driverToTeamMap map[uint64]uint64) []jsondb.RacePosition {
	res := make([]jsondb.RacePosition, 0, len(input))
	for index, driverID := range input {
		var points uint64 = 0
		if eventType == jsondb.SprintEventType || eventType == jsondb.PreSeasonSprintType {
			points = getSprintPointsByIndex(index)
		} else if eventType == jsondb.RaceEventType || eventType == jsondb.PreSeason {
			points = getRacePointsByIndex(index)
		}
		res = append(res, jsondb.RacePosition{
			Position: uint64(index + 1),
			Points:   points,
			DriverID: driverID,
			TeamID:   driverToTeamMap[driverID],
		})
	}

	return res
}

func getSprintPointsByIndex(index int) uint64 {
	points := 8 - index
	if points > 0 {
		return uint64(points)
	}

	return 0
}

func getRacePointsByIndex(index int) uint64 {
	switch index {
	case 0:
		return 25
	case 1:
		return 18
	case 2:
		return 15
	case 3:
		return 12
	case 4:
		return 10
	case 5:
		return 8
	case 6:
		return 6
	case 7:
		return 4
	case 8:
		return 2
	case 9:
		return 1
	default:
		return 0
	}
}
