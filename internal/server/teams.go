package server

import (
	"net/http"
	"strconv"

	"github.com/devnull-twitch/nyooom-backend/pkg/jsondb"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetTeamsHandler(repo jsondb.JsonDatabase) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		teams, err := repo.ListTeams()
		if err != nil {
			logrus.WithError(err).Warn("unable to read teams")
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		events, err := repo.ListEvents()
		if err != nil {
			logrus.WithError(err).Warn("unable to read events")
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		teamsResp := convertTeamsToResponse(teams, events)

		ctx.JSON(http.StatusOK, teamsResp)
	}
}

func GetTeamHandler(repo jsondb.JsonDatabase) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		teamID, err := strconv.Atoi(ctx.Param("team_id"))
		if err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		existing, err := repo.GetTeam(uint64(teamID))
		if err != nil {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}

		respObj := convertTeamFlat(*existing)

		ctx.JSON(http.StatusOK, respObj)
	}
}

func AddTeamHandler(repo jsondb.JsonDatabase) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		newTeam := &jsondb.Team{}

		if err := ctx.BindJSON(newTeam); err != nil {
			logrus.WithError(err).Warn("unable to get user input for new team")
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		if err := repo.AddTeam(newTeam); err != nil {
			logrus.WithError(err).Warn("unable to add team")
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		ctx.Status(http.StatusCreated)
	}
}

func UpdateTeamHandler(repo jsondb.JsonDatabase) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userInputTeam := &jsondb.Team{}

		if err := ctx.BindJSON(userInputTeam); err != nil {
			logrus.WithError(err).Warn("unable to get user input for new team")
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		teamID, err := strconv.Atoi(ctx.Param("team_id"))
		if err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		existing, err := repo.GetTeam(uint64(teamID))
		if err != nil {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}

		existing.Name = userInputTeam.Name
		for index, driver := range userInputTeam.Drivers {
			if len(existing.Drivers) > index {
				existing.Drivers[index].Name = driver.Name
			} else {
				ctx.AbortWithStatus(http.StatusBadRequest)
				return
			}
		}

		if err := repo.UpdateTeam(existing); err != nil {
			logrus.WithError(err).Warn("unable to add team")
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		ctx.Status(http.StatusOK)
	}
}

func DeleteTeamHandler(repo jsondb.JsonDatabase) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		teamID, err := strconv.Atoi(ctx.Param("team_id"))
		if err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		if err := repo.DeleteTeam(uint64(teamID)); err != nil {
			logrus.WithError(err).Warn("unable to add team")
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		ctx.Status(http.StatusOK)
	}
}

func UpdateDriverHandler(repo jsondb.JsonDatabase) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userInputDriver := &jsondb.Driver{}

		teamID, err := strconv.Atoi(ctx.Param("team_id"))
		if err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		driverID, err := strconv.Atoi(ctx.Param("driver_id"))
		if err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		existing, err := repo.GetTeam(uint64(teamID))
		if err != nil {
			logrus.WithError(err).Warn("unable to read team before update")
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		found := false
		for index, existingDriver := range existing.Drivers {
			if existingDriver.ID == uint64(driverID) {
				existing.Drivers[index] = *userInputDriver
				found = true
				break
			}
		}

		if !found {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}

		if err := repo.UpdateTeam(existing); err != nil {
			logrus.WithError(err).Warn("unable to update team with new driver infos")
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		ctx.Status(http.StatusOK)
	}
}
