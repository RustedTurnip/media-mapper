package config

import (
	"encoding/json"
	"io"
	"log"

	dbPackage "github.com/rustedturnip/media-mapper/database"
	"github.com/rustedturnip/media-mapper/database/cache"
	"github.com/rustedturnip/media-mapper/database/databases/tmdb"
	"github.com/rustedturnip/media-mapper/database/databases/tvdb"
)

type database struct {
	API  string            `json:"database"`
	Auth map[string]string `json:"auth"`
}

type config struct {
	Databases []*database `json:"databases"`
}

func GetInstance(authReader io.Reader, api dbPackage.API) (dbPackage.Database, error) {

	configs, err := getConfigs(authReader)
	if err != nil {
		return nil, err
	}

	switch api {
	case dbPackage.TMDB:
		if db, ok := configs[dbPackage.API_name[int(api)]]; ok {
			return cache.New(tmdb.New(db.Auth["apikey"])), nil
		}
		return nil, err

	case dbPackage.TVDB:
		if db, ok := configs[dbPackage.API_name[int(api)]]; ok {
			log.Println("Warning: TVDB only supports TV lookup currently")

			if impl, err := tvdb.New(db.Auth["apikey"], db.Auth["username"], db.Auth["userkey"]); err != nil {
				return nil, err
			} else {
				return cache.New(impl), nil
			}
		}
		return nil, err
	default:
		return nil, nil
	}
}

func getConfigs(reader io.Reader) (map[string]*database, error) {
	//parse json
	var config *config
	decoder := json.NewDecoder(reader)
	err := decoder.Decode(&config)

	if err != nil {
		return nil, err
	}

	mapped := make(map[string]*database)
	for _, db := range config.Databases {
		mapped[db.API] = db
	}

	return mapped, nil
}
