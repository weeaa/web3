package main

import (
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/weeaa/nft/api"
	"github.com/weeaa/nft/db"
	"github.com/weeaa/nft/discord"
	"github.com/weeaa/nft/internal/services"
)

func main() {
	c := make(chan struct{})
	var err error

	if err = godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	pg, err := db.New()
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()
	initModules(router, nil, pg)

	//ft, _ := friendtech.NewClient(nil, true, os.Getenv("NODE_WSS_URL"), os.Getenv("NODE_HTTP_URL"))
	//ft.MonitorFriendTechLogs()

	if err = router.Run(); err != nil {
		gin.DefaultErrorWriter.Write([]byte(fmt.Sprintf("failed starting application: %v", err)))
	}
	<-c
}

func initModules(router *gin.Engine, bot *discord.Bot, db *db.DB) {

	api.InitRoutes(router, services.NewUserService(db))
}
