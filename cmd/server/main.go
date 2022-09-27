package main

import (
	"context"
	"database/sql"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"

	"fakedating/cmd/server/internal/handler"
	"fakedating/pkg/middleware"
	"fakedating/pkg/repository"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var router = mux.NewRouter()

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	db, dbOpenErr := sql.Open("mysql", "root:example@tcp(127.0.0.1:3306)/fakedating")
	if dbOpenErr != nil {
		log.Fatalf("Failed to open database: %v", dbOpenErr)
		return
	}
	defer db.Close()

	authRepository := repository.NewAuth(db)
	userRepository := repository.NewUser(db)

	// Init routes
	h := handler.New(authRepository, userRepository)
	router.HandleFunc("/user/create", h.CreateUser).Methods(http.MethodPost)
	router.HandleFunc("/profiles", h.ListProfiles).Methods(http.MethodGet)
	router.HandleFunc("/swipe", h.Swipe).Methods(http.MethodPost)
	router.HandleFunc("/login", h.Login).Methods(http.MethodPost)
	router.NotFoundHandler = handler.InvalidRoute{}
	router.MethodNotAllowedHandler = handler.InvalidRoute{}

	// Init server
	srv := &http.Server{
		Handler:      middleware.AuthenticateRequest(authRepository, router),
		Addr:         "0.0.0.0:8000",
		WriteTimeout: time.Second,
		ReadTimeout:  5 * time.Second,
	}

	// Run our server in a goroutine so that it doesn't block listening for shutdown
	c := make(chan os.Signal, 1)
	go func() {
		log.Println("HTTP Server starting")
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("HTTP server failed: %v", err)
			os.Exit(1)
		}
	}()

	// Block until a shutdown signal received (CTRL+C)
	signal.Notify(c, os.Interrupt)
	<-c

	// Shut down within context deadline (will shutdown immediately if no active connections)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)

	log.Println("HTTP Server shutdown")
	os.Exit(0)
}
