package main

import (
	"backend/config"
	"backend/database"
	"backend/router"
	"backend/services"
	"backend/utils"

	"github.com/rs/zerolog/log"

	"github.com/gin-gonic/gin"
)

func main() {
	database.Init()
	services.InitializeCompoundData()
	r := gin.Default()
	router.Init(r)
	services.InitRdkit()

	err := r.Run(config.Config.GetString("adress_port"))
	if err != nil {
		utils.LogError(err)
		log.Fatal().Err(err).Msg("Failed to start server")
	}

}
