package tmdb

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/rustedturnip/media-mapper/database"
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
	apiKey     string
	httpClient *http.Client
}

func New(key string) database.Database {
	return &TMDB{
		apiKey:     key,
		httpClient: &http.Client{},
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
		movies = append(movies, buildMovie(movie))
	}

	return movies
}

func (db *TMDB) searchMovies(title string) (*movieSearch, error) {

	searchQuery := url.QueryEscape(title)

	resp, err := db.httpClient.Get(fmt.Sprintf(apiMovieSearch, db.apiKey, searchQuery))
	if err != nil {
		return nil, err
	}

	var searchResults *movieSearch
	err = database.ReadJsonToStruct(resp.Body, &searchResults)

	if err != nil {
		return nil, err
	}

	return searchResults, nil
}

func buildMovie(result movieSearchResult) *types.Movie {

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

		showData := db.fetchTVShow(&show)
		bShow := buildTV(showData)
		if bShow != nil {
			shows = append(shows, bShow)
		}
	}

	return shows
}

func (db *TMDB) searchTV(title string) (*tvSearch, error) {

	searchQuery := url.QueryEscape(title)

	resp, err := db.httpClient.Get(fmt.Sprintf(apiTVSearch, db.apiKey, searchQuery))
	if err != nil {
		return nil, err
	}

	var searchResults *tvSearch
	err = database.ReadJsonToStruct(resp.Body, &searchResults)

	if err != nil {
		return nil, err
	}

	return searchResults, nil
}

//fetches all show specific data and returns as tvShow
func (db *TMDB) fetchTVShow(result *tvSearchResult) *tvShow {
	errScope := "TV Build error: %s"

	tvResp, err := db.httpClient.Get(fmt.Sprintf(apiTVByID, result.ID, db.apiKey))
	if err != nil {
		log.Println(fmt.Sprintf(errScope, err))
		return nil
	}

	var tvObj *tvShow
	err = database.ReadJsonToStruct(tvResp.Body, &tvObj)
	if err != nil {
		log.Println(fmt.Sprintf(errScope, err))
		return nil
	}

	for _, s := range tvObj.Seasons {
		seriesResp, err := db.httpClient.Get(fmt.Sprintf(apiSeriesByNumber, result.ID, s.SeasonNumber, db.apiKey))
		if err != nil {
			log.Println(fmt.Sprintf(errScope, err))
			return nil
		}

		var sObj *tvShowSeriesData
		err = database.ReadJsonToStruct(seriesResp.Body, &sObj)
		if err != nil {
			log.Println(fmt.Sprintf(errScope, err))
			return nil
		}

		s.SeasonData = sObj
	}

	return tvObj
}

//builds tvShow into types.TV
func buildTV(show *tvShow) *types.TV {
	//Build TV Series
	tvBuilder := builder.NewTVBuilder()

	for _, sInfo := range show.Seasons {
		//Build Series
		seriesBuilder := builder.NewSeriesBuilder()
		for _, e := range sInfo.SeasonData.Episodes {

			//Build Episode
			episodeBuilder := builder.NewEpisodeBuilder()
			episodeBuilder.
				WithTitle(e.Name).
				WithNumber(e.EpisodeNumber)

			seriesBuilder.WithEpisode(episodeBuilder)
		}

		seriesBuilder.
			WithTitle(sInfo.Name).
			WithNumber(sInfo.SeasonNumber)

		tvBuilder.WithSeries(seriesBuilder)
	}

	tvBuilder.
		WithTitle(show.Name).
		WithSeriesCount(show.NumberOfSeasons)

	return tvBuilder.Build()
}
