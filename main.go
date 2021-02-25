package main

import (
	"flag"

	"github.com/babolivier/london-covid-vaccination/api"
	"github.com/babolivier/london-covid-vaccination/config"
	"github.com/babolivier/london-covid-vaccination/miner"
	"github.com/babolivier/london-covid-vaccination/storage"

	"github.com/sirupsen/logrus"
)

var (
	configPath = flag.String("c", "config.yaml", "Path to the configuration file")
)

func main() {
	// Parse the command-line arguments.
	flag.Parse()

	// Configure the logger.
	logrus.SetFormatter(
		&logrus.TextFormatter{
			TimestampFormat: "2006-02-01 15:04:05.999",
			FullTimestamp:   true,
		},
	)

	// Load the configuration file.
	cfg, err := config.NewConfig(*configPath)
	if err != nil {
		panic(err)
	}

	// Connect to the database.
	db, err := storage.NewDatabase(cfg.Database)
	if err != nil {
		panic(err)
	}

	// Instantiate and start the miner.
	m := miner.NewMiner(db)
	m.Start()

	// Start the API HTTP server.
	api.StartApiServer(cfg.Api, db)
}
