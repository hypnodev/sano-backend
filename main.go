package main

import (
	"context"
	"github.com/robfig/cron/v3"
	"log"
	"sano/api"
	"sano/config"
	"sano/database"
	"sano/services"
	"time"
)

func main() {
	config.CheckConfig()
	cfg := config.LoadConfig()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	database.ConnectToDatabase(ctx, cfg.App.Database.Url)
	defer cancel()
	defer func() {
		if err := database.Client.Disconnect(ctx); err != nil {
			log.Panicln(err)
		}
	}()
	log.Printf("Connected to MongoDB database [%s]", cfg.App.Database.Url)

	c := cron.New()
	defer c.Stop()

	services.RunLookup(cfg.Services, cfg.HealthCheck.Cron, c)
	c.Start()
	log.Printf("[%d] cron started for [%d] services", len(c.Entries()), len(cfg.Services))

	api.Start(cfg.App.Port)
}
