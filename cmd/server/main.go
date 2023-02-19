package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/devnull-twitch/nyooom-backend/internal/server"
	"github.com/devnull-twitch/nyooom-backend/pkg/jsondb"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	repo := jsondb.CreateFileDatabase()

	editors := make([]*server.EditorLogin, 0)
	editorConfigStr := os.Getenv("EDITORS")
	for _, credentials := range strings.Split(editorConfigStr, ";") {
		credentialParts := strings.Split(credentials, "=")
		if len(credentialParts) != 2 {
			panic(fmt.Errorf("invalid editor credentials. must be user=pass"))
		}

		editors = append(editors, &server.EditorLogin{Username: credentialParts[0], Password: credentialParts[1]})
	}
	editorCheckMW := server.GetEditorMiddleware(editors)

	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AddAllowHeaders("Authorization")
	r.Use(cors.New(corsConfig))

	r.GET("/user-check", editorCheckMW, server.GetNoopHandler())

	r.GET("/team", server.GetTeamsHandler(repo))
	r.GET("/team/:team_id", server.GetTeamHandler(repo))
	r.GET("/race", server.GetEventsHandler(repo))
	r.GET("/race/:race_id", server.GetEventHandler(repo))

	r.POST("/team", editorCheckMW, server.AddTeamHandler(repo))
	r.PUT("/team/:team_id", editorCheckMW, server.UpdateTeamHandler(repo))
	r.DELETE("/team/:team_id", editorCheckMW, server.DeleteTeamHandler(repo))
	r.PUT("/team/:team_id/:driver_id", editorCheckMW, func(ctx *gin.Context) {})
	r.POST("/race", editorCheckMW, server.CreateRaceEventHandler(repo))
	r.PUT("/race/:race_id", editorCheckMW, server.UpdateRaceEventHandler(repo))
	r.DELETE("/race/:race_id", editorCheckMW, server.DeleteRaceEventHandler(repo))

	r.Run(os.Getenv("WEBSERVER_ADDRESS"))
}
