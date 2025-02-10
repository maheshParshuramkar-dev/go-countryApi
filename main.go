package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/go-countryApi/initializers"
)

func main() {
	// intialize all the connections and other ops
	initializers.InitializeOps()

	notify := make(chan os.Signal, 1)
	signal.Notify(notify, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-notify
		log.Println("Received Graceful shutdown signal")
		initializers.StopServices()
	}()

	// setup and start server
	initializers.SetupAndStartSrv()

	// wait for graceful shutdown
	wg.Wait()

	log.Println("Exiting ...")
}
