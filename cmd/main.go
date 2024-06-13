package main

import (
    "log"
    "strconv"
    "github.com/ihgazi/go-chat/db"
    "github.com/ihgazi/go-chat/internal/user"
    "github.com/ihgazi/go-chat/router"
    "github.com/ihgazi/go-chat/config"
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

    router.InitRouter(userHndlr)
    addr := conf.ServerConfig.Host + ":" + strconv.Itoa(conf.ServerConfig.Port)
    if err := router.Run(addr); err != nil {
        log.Fatalf("Error: %v", err)
    }
}
