package tvdb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/rustedturnip/media-mapper/database"
	"github.com/rustedturnip/media-mapper/types"
	"github.com/rustedturnip/media-mapper/types/builder"
)

const (
	apiBase               = "https://api.thetvdb.com"
	apiLogin              = "/login"
	apiSeriesSearch       = "/search/series"
	apiSeriesByID         = "/series/%d"
	apiEpisodesBySeriesID = "/series/%d/episodes"

	httpHeaderAuth = "Authorization"

	specialEpisodes = 0
)

type TVDB struct {
	auth            auth //details used to get token
	token           string
	requestTemplate http.Request //used so only URL needs adding in future requests
	httpClient      *http.Client
}

func New(apiKey, username, userkey string) (database.Database, error) {

	tvdb := &TVDB{
		auth: auth{
			APIKey:   apiKey,
			Username: username,
			UserKey:  userkey,
		},
		requestTemplate: http.Request{
			Method: http.MethodGet,
			Header: http.Header{},
		},
		httpClient: &http.Client{},
	}

	body, err := json.Marshal(tvdb.auth)
	if err != nil {
		return nil, err
	}

	//get JWT (token)
	resp, err := tvdb.httpClient.Post(fmt.Sprintf("%s%s", apiBase, apiLogin), "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var token *token
	err = database.ReadJsonToStruct(resp.Body, &token)
	if err != nil {
		return nil, err
	}

	tvdb.token = token.Token
	tvdb.requestTemplate.Header.Add(httpHeaderAuth, fmt.Sprintf("Bearer %s", tvdb.token))

	return tvdb, nil
}

//v3 of the TVDB API doesn't support movie search
//TODO - implement v4 when available
func (db *TVDB) SearchMovies(title string) []*types.Movie {

	return nil
}

func (db *TVDB) SearchTV(title string) []*types.TV {

	searchResults, err := db.searchTV(title)
	if err != nil {
		return []*types.TV{} //empty slice
	}

	//compile list of shows (built to *type.TV)
	var shows []*types.TV
	for _, show := range searchResults.Results {

		data, err := db.fetchTV(show)
		if err != nil {
			continue
		}

		if bShow := buildTV(data); bShow != nil {
			shows = append(shows, bShow)
		}
	}

	return shows
}

//queries search endpoint with specified title
func (db *TVDB) searchTV(title string) (*tvSearch, error) {

	req := &db.requestTemplate

	if link, err := url.Parse(fmt.Sprintf("%s%s", apiBase, apiSeriesSearch)); err == nil {
		req.URL = link
	} else {
		return nil, err
	}

	//add api call parameters
	q := req.URL.Query()
	q.Set("name", title)
	req.URL.RawQuery = q.Encode()

	resp, err := db.httpClient.Do(req)
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

//queries series by ID to get series data
func (db *TVDB) fetchTV(result *tvSearchResult) (*tvShow, error) {

	//fetch show data
	req := &db.requestTemplate
	url, _ := url.Parse(fmt.Sprintf("%s%s", apiBase, fmt.Sprintf(apiSeriesByID, result.ID)))
	req.URL = url

	resp, err := db.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error requesting tv show - %s", err.Error())
	}

	var tv *tv = &tv{}
	err = database.ReadJsonToStruct(resp.Body, &tv)

	if err != nil {
		return nil, fmt.Errorf("error reading tv response - %s", err.Error())
	}

	//fetch show episodes
	episodes, err := db.getEpisodes(result.ID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving episodes - %s", err.Error())
	}

	tv.Show.Series = &tvSeriesEpisodes{
		Episodes: episodes,
	}

	return tv.Show, nil
}

//queries for episodes pertaining to series (by series ID)
func (db *TVDB) getEpisodes(seriesID uint64) ([]*episode, error) {

	var results []*episode

	nextPage := 1
	req := &db.requestTemplate

	if link, err := url.Parse(fmt.Sprintf("%s%s", apiBase, fmt.Sprintf(apiEpisodesBySeriesID, seriesID))); err == nil {
		req.URL = link
	} else {
		return nil, err
	}

	for {
		if nextPage == 0 {
			break
		}

		q := req.URL.Query()
		q.Set("page", strconv.Itoa(nextPage))
		req.URL.RawQuery = q.Encode()

		resp, err := db.httpClient.Do(req)
		if err != nil {
			return nil, err //if error, discard all
		}

		var episodeResults *tvSeriesEpisodes = &tvSeriesEpisodes{}
		err = database.ReadJsonToStruct(resp.Body, &episodeResults)
		if err != nil {
			return nil, err //if error, discard all
		}

		results = append(results, episodeResults.Episodes...)
		nextPage = episodeResults.Links.Next
	}

	return results, nil
}

//builds types.TV object based on json structs built from api responses
func buildTV(show *tvShow) *types.TV {

	//group episodes by series number
	groupedEpisodes := make(map[int][]*builder.EpisodeBuilder)

	for _, episode := range show.Series.Episodes {
		seriesNum := episode.AiredSeason

		//build episode
		eb := builder.NewEpisodeBuilder()
		eb.
			WithTitle(episode.EpisodeName).
			WithNumber(episode.AiredEpisodeNumber)

		if _, ok := groupedEpisodes[seriesNum]; !ok {
			groupedEpisodes[seriesNum] = []*builder.EpisodeBuilder{}
		}

		groupedEpisodes[seriesNum] = append(groupedEpisodes[seriesNum], eb)
	}

	//Get show's number of series
	seriesCount := len(groupedEpisodes)
	if _, ok := groupedEpisodes[specialEpisodes]; ok {
		seriesCount -= 1 //Ignore series with number 0 as reserved for special episodes
	}

	//start tv build
	tvb := builder.NewTVBuilder()
	tvb.
		WithTitle(show.SeriesName).
		WithSeriesCount(seriesCount)

	//build series based on grouped episodes
	for seriesNum, episodes := range groupedEpisodes {

		//Build season/specific series name
		seasonName := fmt.Sprintf("Season %d", seriesNum)
		if seriesNum == 0 {
			seasonName = "Specials"
		}

		//build series
		sb := builder.NewSeriesBuilder()
		sb.
			WithNumber(seriesNum).
			WithTitle(seasonName)

		for _, e := range episodes {
			sb.WithEpisode(e)
		}

		tvb.WithSeries(sb)
	}

	return tvb.Build()
}
