package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"samwang0723/jarvis/db"
	"samwang0723/jarvis/db/dal"
	"samwang0723/jarvis/handlers"
	"samwang0723/jarvis/services"
	"syscall"
	"time"
)

func main() {
	// service initialization
	config := &db.Config{
		User:     "jarvis",
		Password: "password",
		Host:     "tcp(localhost:3306)",
		Database: "jarvis",
	}
	dalService := dal.New(dal.WithDB(db.GormFactory(config)))
	dataService := services.New(services.WithDAL(dalService))
	handler := handlers.New(dataService)

	// graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// executing batch download
	handler.BatchingDownload(context.Background(), -6, 5000)
	//	req := &dto.ListDailyCloseRequest{
	//		Offset: 0,
	//		Limit:  10,
	//		SearchParams: &dto.ListDailyCloseSearchParams{
	//			StockIDs: &[]string{"2330", "3035", "3707"},
	//			Start:    "20211007",
	//		},
	//	}
	//	resp, err := handler.ListDailyClose(context.Background(), req)
	//	if err != nil {
	//		log.Printf("listing DailyClose failed: %s\n", err)
	//	}
	//	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	//	data, _ := json.Marshal(&resp)
	//	log.Printf("json response: %s\n", string(data))

	<-done
	log.Println("server stopped")

	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// extra handling here
		cancel()
	}()
	log.Println("server exited properly")
}
