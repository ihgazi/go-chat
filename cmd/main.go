package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/ihgazi/go-chat/config"
	"github.com/ihgazi/go-chat/db"
	"github.com/ihgazi/go-chat/internal/user"
	"github.com/ihgazi/go-chat/internal/ws"
	"github.com/ihgazi/go-chat/router"
)

func main() {
	dbConn, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	userRep := user.NewRepository(dbConn.GetDB())
	userSvc := user.NewService(userRep)
	userHndlr := user.NewHandler(userSvc)

	wsRep := ws.NewRepository(dbConn.GetDB())
	hub := ws.NewHub(wsRep)
	wsSvc := ws.NewService(wsRep, hub)
	wsHndlr := ws.NewHandler(wsSvc)
	go hub.Run()

	conf := config.LoadConfig()
	r := router.Init(userHndlr, wsHndlr)
	port, _ := strconv.Atoi(conf.ServerPort)
	addr := fmt.Sprintf("%s:%d", conf.ServerHost, port)

	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// Run server in a goroutine so it doesn't block
	go func() {
		log.Printf("Starting server on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish existing requests
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Setting connected users to offline and closing websockets...")
	hub.Shutdown()

	log.Println("Closing database connection...")
	dbConn.Close()

	log.Println("Server successfully exited")
}
