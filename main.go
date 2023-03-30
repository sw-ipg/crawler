package main

import (
	"context"
	"crawler/app"
	"crawler/config"
	"crawler/httpserver"
	"log"
	"os"
	"os/signal"
)

func main() {
	s, err := config.ReadConfig()
	if err != nil {
		log.Fatalf("FATAL: unable to run app due to read config err: %s", err)
	}

	a := app.NewApp(s)
	crawlCtx, crawlCtxCancel := context.WithCancel(context.TODO())
	defer crawlCtxCancel()
	a.RunCrawlPipeline(crawlCtx)
	log.Printf("INFO: crawl pipeline has been started")

	serverCtx, serverCtxCancel := context.WithCancel(context.TODO())
	defer serverCtxCancel()
	server := httpserver.New(serverCtx, a, s.ServerPort)
	if err = server.ListenAndServe(); err != nil {
		log.Fatalf("ERROR: cannot run http server: %s", err)
	}

	doneSignal := make(chan os.Signal)
	signal.Notify(doneSignal, os.Interrupt, os.Kill)
	<-doneSignal
	log.Printf("INFO: app gracefully stopped. Bye :)")
}
