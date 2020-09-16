package main

import (
	"flag"
	"log"

	cfg "github.com/rustedturnip/media-mapper/config"
	"github.com/rustedturnip/media-mapper/controller"
	"github.com/rustedturnip/media-mapper/dbs"
	"github.com/rustedturnip/media-mapper/filing"
)

type config struct {
	DB  string   `json:"database"`
	cfg []string `json:"config"`
}

const (
	version = "v0.0.0-dev"
)

var (
	database string
	auth     string
	location string
)

func init() {
	//todo use different flag package
	flag.StringVar(&database, "database", "TMDB", "database to extract data from")
	flag.StringVar(&auth, "authentication", "/Users/samuel/go/src/github.com/rustedturnip/media-mapper/configs.json", "location of auth")
	flag.StringVar(&location, "location", "/Users/samuel/go/src/github.com/rustedturnip/media-mapper/tmp-test", "location of files to be formatted")
}

func main() {

	//create DB instance
	db, ok := dbs.API_value[database]
	if !ok {
		log.Fatalf("Unssupported network specified: %s", database)
	}

	api, err := cfg.GetInstance(auth, db)
	if err != nil {
		log.Fatalf("Unable to create network instance for %s", database)
	}

	//create Filer instance
	var filer *filing.Filer
	if filer, err = filing.New(location); err != nil {
		log.Fatalf("File handler failed to initialise: %s", err.Error())
	}

	worker := controller.New(api, filer)
	worker.Do()
}
