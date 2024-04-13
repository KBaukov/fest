package main

import (
	"context"
	"fest/app"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func setupSignalHandler(cancel context.CancelFunc, wg *sync.WaitGroup) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		select {
		case <-c:
			fmt.Println("############ Приложение завершает работу, дождитесь полной остановки #########")
			wg.Done()
			cancel()
			return
		}
	}()
}

func main() {

	//db.InitDB()
	//dao.DaoInit()
	//dict.DictInit()

	// ticket, err := dao.GetTicketByLkId("SOLAR_#1702636398", 1051)
	// log.Println(ticket, err)
	// log.Fatal()

	var wg sync.WaitGroup
	wg.Add(3)

	// setup context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())

	// init signal handler
	go setupSignalHandler(cancel, &wg)

	// Init Common Api
	go app.ServApi(ctx, &wg)

	//Start SesssionCleaner
	go app.StartCleaner(ctx, &wg)

	fmt.Println("#################### Festival BackEnd App started: ########################")

	wg.Wait()

	fmt.Println("#################### Festival BackEnd App is stopped: ######################")

}
