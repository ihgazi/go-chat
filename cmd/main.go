package main

import (
	"fmt"
	"log"

	"github.com/ihgazi/go-chat/config"
	"github.com/ihgazi/go-chat/db"
	"github.com/ihgazi/go-chat/internal/user"
	"github.com/ihgazi/go-chat/router"
)

func main() {
	dbConn, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Repository is injected with dbConn, takes User struct and updates database
	// Service injected with Repository, takes CreateUserReq and creates User Struct
	// Handler injected with Service, parses the Json data and creates the CreateUserReq
	userRep := user.NewRepository(dbConn.GetDB())
	userSvc := user.NewService(userRep)
	userHndlr := user.NewHandler(userSvc)

	// Pulling in server config
	conf := config.LoadConfig("config.toml")

	// Initializing GIN router and running the server
	r := router.Init(userHndlr)
	addr := fmt.Sprintf("%s:%d", conf.ServerConfig.Host, conf.ServerConfig.Port)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
