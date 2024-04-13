package app

import (
	"context"
	"fest/config"
	dao "fest/db"
	ws "fest/handlers"
	"fmt"
	"sync"
	"time"
)

var (
	ticker *time.Ticker
	isBusy bool
	db     dao.Database
)

func StartCleaner(ctx context.Context, wg *sync.WaitGroup) {

	cfg, _ := config.LoadConfig("config.json")
	if !cfg.SessionSettings.CleanerEnable {
		wg.Done()
		return
	}

	go func() {
		select {
		case <-ctx.Done():
			ticker.Stop()
			fmt.Println("#################### SessionCleaner Ticker is Stopped #######################")
			wg.Done()
		}
	}()

	psqlconn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DbConnection.DBHost,
		cfg.DbConnection.DBPort,
		cfg.DbConnection.DBUser,
		cfg.DbConnection.DBPass,
		cfg.DbConnection.DBName,
	)
	db, _ = dao.NewDB(psqlconn)

	StartCleanTicker()
}

func StartCleanTicker() {
	fmt.Println("#################### SessionCleaner Ticker is Started ##########################")
	ticker = time.NewTicker(time.Second * 10)

	for now := range ticker.C {
		if !isBusy {
			cleanSession(now)
		}
	}
}

func cleanSession(t time.Time) {
	devices, err := db.ClearDieSession()
	if err != nil {
	}
	for _, devId := range devices {
		ws.SendMsgByDevId(devId, "{\"action\": \"closeSession\", \"device_id\":\""+devId+"\"}")
	}
}
