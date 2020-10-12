package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	cfg "github.com/rustedturnip/media-mapper/config"
	"github.com/rustedturnip/media-mapper/controller"
	"github.com/rustedturnip/media-mapper/dbs"
	"github.com/rustedturnip/media-mapper/filing"
)

const (
	version = "v0.1.0-dev"
)

var (
	AuthConfigs string //base 64 encoded, initialised at build time

	versionFlag bool
	database    string
	auth        string
	location    string
)

func init() {
	flag.BoolVar(&versionFlag, "version", false, "database to extract data from")

	flag.StringVar(&database, "database", "TMDB", "database to extract data from")
	flag.StringVar(&auth, "auth", "", "location of auth")
	flag.StringVar(&location, "location", "", "location of files to be formatted")
}

func main() {

	if versionFlag {
		fmt.Print(fmt.Sprintf("media-mapper version: %s", version))
		return
	}

	//create DB instance
	db, ok := dbs.API_value[database]
	if !ok {
		log.Fatalf("Unssupported network specified: %s", database)
	}

	authReader, err := getAuthReader()
	if err != nil {
		log.Fatalf(err.Error())
	}

	api, err := cfg.GetInstance(authReader, db)
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

func getAuthReader() (io.Reader, error) {

	//expects base64 to be encoded
	if AuthConfigs != "" {
		return base64.NewDecoder(base64.StdEncoding, strings.NewReader(AuthConfigs)), nil
	}

	if auth != "" {
		return os.Open(auth)
	}

	return nil, fmt.Errorf("failed to find database credentials")
}
