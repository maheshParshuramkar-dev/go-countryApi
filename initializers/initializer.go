package initializers

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/go-countryApi/api"
	"github.com/go-countryApi/config"
)

var server *http.Server

func InitializeOps() {
	// get config and create in mem cache
	config.GetConfig()
	api.CreateInMemCache()
}

// setting up server and start it
func SetupAndStartSrv() *http.Server {
	http.HandleFunc(config.AppConfig.Prefix+"/ping", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "pong")
	})

	http.HandleFunc(config.AppConfig.Prefix+"/search", api.GetCountryData)

	server = &http.Server{
		Addr:         ":" + config.AppConfig.Server.Port,
		ReadTimeout:  time.Duration(config.AppConfig.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(config.AppConfig.Server.WriteTimeout) * time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("error while listen :", err)
		}
	}()

	log.Printf("Server is up and running on http://localhost:%v%v/ping\n", config.AppConfig.Server.Port, config.AppConfig.Prefix)
	return server
}

// stop server and close server
func StopServices() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Error while shutting down server", err.Error())
	}

	log.Println("Server gracefully stopped...")
}
