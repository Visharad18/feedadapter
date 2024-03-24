package main

import (
	"context"

	"github.com/Visharad18/feedadapter/api"
	"github.com/Visharad18/feedadapter/app"
	"github.com/Visharad18/feedadapter/config"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg, err := config.NewConfig()
	lg := logrus.New()
	if err != nil {
		logrus.Fatalf("error in parsing config %s", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app := app.NewApp(cfg, lg)
	go app.Run(ctx)

	if _, err := api.NewServer(cfg, app); err != nil {
		lg.Fatalf("error in creating server: %s", err)
	}

}
