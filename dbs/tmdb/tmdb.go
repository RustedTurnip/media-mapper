package tmdb

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/rustedturnip/media-mapper/dbs"
	"github.com/rustedturnip/media-mapper/types"
	"github.com/rustedturnip/media-mapper/types/builder"
)

const (
	//movie calls
	apiMovieSearch = "https://api.themoviedb.org/3/search/movie?api_key=%s&language=en-GB&query=%s&page=1&include_adult=true"

	//tv calls
	apiTVSearch       = "https://api.themoviedb.org/3/search/tv?api_key=%s&language=en-GB&query=%s&page=1&include_adult=true"
	apiTVByID         = "https://api.themoviedb.org/3/tv/%d?api_key=%s&language=en-GB"
	apiSeriesByNumber = "https://api.themoviedb.org/3/tv/%d/season/%d?api_key=%s&language=en-GB"

	apiDateFormat = "2006-01-02"
)

type TMDB struct {
	apiKey string
}

func New(key string) dbs.Database {
	return &TMDB{
		apiKey: key,
	}
}

func (db *TMDB) SearchMovies(title string) []*types.Movie {

	results, err := db.searchMovies(title)

	if err != nil {
		log.Println(fmt.Sprintf("Failed getting Movie results with error: %s", err.Error()))
		return []*types.Movie{} //empty slice
	}

	var movies []*types.Movie
	for _, movie := range results.Results {
		movies = append(movies, db.buildMovie(movie))
	}

	return movies
}

func (db *TMDB) searchMovies(title string) (*movieSearch, error) {

	searchQuery := url.QueryEscape(title)

	resp, err := http.Get(fmt.Sprintf(apiMovieSearch, db.apiKey, searchQuery))
	if err != nil {
		return nil, err
	}

	var searchResults *movieSearch
	err = dbs.ReadJsonToStruct(resp.Body, &searchResults)

	if err != nil {
		return nil, err
	}

	return searchResults, nil
}

func (db *TMDB) buildMovie(result movieSearchResult) *types.Movie {

	movieBuilder := builder.NewMovieBuilder()

	date, err := time.Parse(apiDateFormat, result.ReleaseDate)
	if err != nil {
		log.Println(fmt.Sprintf("Movie build error: %s", err.Error()))
	}

	movie := movieBuilder.
		WithTitle(result.Title).
		WithReleaseDate(date).
		Build()

	return movie
}

func (db *TMDB) SearchTV(title string) []*types.TV {
	results, err := db.searchTV(title)

	if err != nil {
		log.Println(fmt.Sprintf("Failed getting TV results with error: %s", err.Error()))
		return []*types.TV{} //empty slice
	}

	var shows []*types.TV
	for _, show := range results.Results {

		bShow := db.buildTV(show)
		if bShow != nil {
			shows = append(shows, bShow)
		}
	}

	return shows
}

func (db *TMDB) searchTV(title string) (*tvSearch, error) {

	searchQuery := url.QueryEscape(title)

	resp, err := http.Get(fmt.Sprintf(apiTVSearch, db.apiKey, searchQuery))
	if err != nil {
		return nil, err
	}

	var searchResults *tvSearch
	err = dbs.ReadJsonToStruct(resp.Body, &searchResults)

	if err != nil {
		return nil, err
	}

	return searchResults, nil
}

func (db *TMDB) buildTV(result tvSearchResult) *types.TV {
	errScope := "TV Build error: %s"

	tvResp, err := http.Get(fmt.Sprintf(apiTVByID, result.ID, db.apiKey))
	if err != nil {
		log.Println(fmt.Sprintf(errScope, err))
		return nil
	}

	var tvObj *tv
	err = dbs.ReadJsonToStruct(tvResp.Body, &tvObj)
	if err != nil {
		log.Println(fmt.Sprintf(errScope, err))
		return nil
	}

	//Build TV Series
	tvBuilder := builder.NewTVBuilder()
	for _, s := range tvObj.Seasons {
		seriesResp, err := http.Get(fmt.Sprintf(apiSeriesByNumber, result.ID, s.SeasonNumber, db.apiKey))
		if err != nil {
			log.Println(fmt.Sprintf(errScope, err))
			return nil
		}

		var sObj *series
		err = dbs.ReadJsonToStruct(seriesResp.Body, &sObj)
		if err != nil {
			log.Println(fmt.Sprintf(errScope, err))
			return nil
		}

		//Build Series
		seriesBuilder := builder.NewSeriesBuilder()
		for _, e := range sObj.Episodes {

			//Build Episode
			episodeBuilder := builder.NewEpisodeBuilder()
			episodeBuilder.
				WithTitle(e.Name).
				WithNumber(e.EpisodeNumber)

			seriesBuilder.WithEpisode(episodeBuilder)
		}

		seriesBuilder.
			WithTitle(sObj.Name).
			WithNumber(sObj.SeasonNumber)

		tvBuilder.WithSeries(seriesBuilder)
	}

	tvBuilder.
		WithTitle(tvObj.Name).
		WithSeriesCount(tvObj.NumberOfSeasons)

	return tvBuilder.Build()
}
