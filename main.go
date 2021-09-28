package main

import (
	"context"
	"log"
	"samwang0723/jarvis/helper"
	"samwang0723/jarvis/server/handlers"
)

func main() {
	d := helper.GetDate(0, 0, -5)
	log.Println("Request Date: ", d)

	resp, err := handlers.DownloadDailyCloses(context.Background(), d)
	if err != nil {
		log.Fatalf("DownloadDailyClose error: %v\n", err)
	}

	for k, v := range resp {
		log.Printf("%s: %v\n", k, v)
	}

}
