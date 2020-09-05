package config

import (
	"encoding/json"
	"os"

	"github.com/rustedturnip/media-mapper/dbs"
	"github.com/rustedturnip/media-mapper/dbs/tmdb"
)

type database struct {
	API  string            `json:"database"`
	Auth map[string]string `json:"auth"`
}

type config struct {
	Databases []*database `json:"databases"`
}

func GetInstance(authLocation string, api dbs.API) (dbs.Database, error) {

	configs, err := getConfigs(authLocation)
	if err != nil {
		return nil, err
	}

	switch api {
	case dbs.TMDB:
		if db, ok := configs[dbs.API_name[int(api)]]; ok {
			return tmdb.New(db.Auth["apikey"]), nil
		}
		return nil, err

	case dbs.TVDB:
		//TODO
		return nil, nil
	default:
		return nil, nil
	}
}

func getConfigs(location string) (map[string]*database, error) {
	//read in json auth
	data, err := os.Open(location)
	if err != nil {
		return nil, err
	}

	//parse json
	var config *config
	decoder := json.NewDecoder(data)
	err = decoder.Decode(&config)

	if err != nil {
		return nil, err
	}

	mapped := make(map[string]*database)
	for _, db := range config.Databases {
		mapped[db.API] = db
	}

	return mapped, nil
}
