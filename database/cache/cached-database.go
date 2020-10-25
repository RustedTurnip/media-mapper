package cache

import (
	"github.com/rustedturnip/media-mapper/database"
	"github.com/rustedturnip/media-mapper/types"
)

//Implements cached database to prevent excessive calls to API
type DatabaseCache struct {
	movieSearchCache map[string][]*types.Movie
	tvSearchCache    map[string][]*types.TV

	dbClient database.Database
}

func New(client database.Database) *DatabaseCache {
	return &DatabaseCache{
		movieSearchCache: make(map[string][]*types.Movie),
		tvSearchCache:    make(map[string][]*types.TV),
		dbClient:         client,
	}
}

func (db *DatabaseCache) SearchMovies(title string) []*types.Movie {
	if _, ok := db.movieSearchCache[title]; ok {
		return db.movieSearchCache[title]
	}

	results := db.dbClient.SearchMovies(title)
	db.movieSearchCache[title] = results

	return results
}

func (db *DatabaseCache) SearchTV(title string) []*types.TV {
	if _, ok := db.tvSearchCache[title]; ok {
		return db.tvSearchCache[title]
	}

	results := db.dbClient.SearchTV(title)
	db.tvSearchCache[title] = results

	return results
}
