package dbs

import (
	"encoding/json"
	"github.com/rustedturnip/media-mapper/types"
	"io"
	"io/ioutil"
)

type Database interface {
	SearchMovies(string) []*types.Movie
	SearchTV(string) []*types.TV
}

type API int

const (
	TMDB API = iota
	TVDB
)

var API_value = map[string]API{
	"TMDB": TMDB,
	"TVDB": TVDB,
}

var API_name = map[int]string{
	0: "TMDB",
	1: "TVDB",
}

func ReadJsonToStruct(reader io.ReadCloser, obj interface{}) error {

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, obj)
	if err != nil {
		return err
	}

	return nil
}
