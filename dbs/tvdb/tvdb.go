package tvdb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/rustedturnip/media-mapper/dbs"
	"github.com/rustedturnip/media-mapper/types"
	"github.com/rustedturnip/media-mapper/types/builder"
)

const (
	apiBase               = "https://api.thetvdb.com"
	apiLogin              = "/login"
	apiSeriesSearch       = "/search/series"
	apiEpisodesBySeriesID = "/series/%d/episodes"

	httpHeaderAuth = "Authorization"
)

type TVDB struct {
	auth            auth //details used to get token
	token           string
	requestTemplate http.Request //used so only URL needs adding in future requests
	httpCli         http.Client
}

func New(apiKey, username, userkey string) dbs.Database {

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
		httpCli: http.Client{},
	}

	body, err := json.Marshal(tvdb.auth)
	if err != nil {
		return nil
	}

	//get JWT (token)
	resp, err := http.Post(fmt.Sprintf("%s%s", apiBase, apiLogin), "application/json", bytes.NewReader(body))
	if err != nil {
		return nil
	}

	var token *token
	err = dbs.ReadJsonToStruct(resp.Body, &token)
	if err != nil {
		return nil
	}

	tvdb.token = token.Token
	tvdb.requestTemplate.Header.Add(httpHeaderAuth, fmt.Sprintf("Bearer %s", tvdb.token))

	return tvdb
}

//v3 of the TVDB API doesn't support movie search
//TODO - implement v4 when available
func (db *TVDB) SearchMovies(title string) []*types.Movie {

	return nil
}

func (db *TVDB) SearchTV(title string) []*types.TV {

	searchResults, err := db.searchTV(title)
	if err != nil {
		log.Println(fmt.Sprintf("Failed getting TV results with error: %s", err.Error()))
		return []*types.TV{} //empty slice
	}

	var shows []*types.TV
	for _, show := range searchResults.Results {
		if bShow := db.buildTV(show); bShow != nil {
			shows = append(shows, bShow)
		}
	}

	return shows
}

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

	client := http.Client{}

	resp, err := client.Do(req)
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

func (db *TVDB) buildTV(result tvSearchResult) *types.TV {

	//build episodes and group by season/seriesNum
	episodes := db.getEpisodes(result.ID)
	groupedEpisodes := make(map[int][]*builder.EpisodeBuilder)

	for _, episode := range episodes {
		seriesNum := episode.AiredSeason

		eb := builder.NewEpisodeBuilder()
		eb.
			WithTitle(episode.EpisodeName).
			WithNumber(episode.AiredEpisodeNumber)

		if _, ok := groupedEpisodes[seriesNum]; !ok {
			groupedEpisodes[seriesNum] = []*builder.EpisodeBuilder{}
		}

		groupedEpisodes[seriesNum] = append(groupedEpisodes[seriesNum], eb)
	}

	//start tv build
	tvb := builder.NewTVBuilder()
	tvb.
		WithTitle(result.SeriesName).
		WithSeriesCount(len(groupedEpisodes))

	//build series based on grouped episodes
	for seriesNum, episodes := range groupedEpisodes {

		sb := builder.NewSeriesBuilder()
		sb.
			WithNumber(seriesNum).
			WithTitle(fmt.Sprintf("Season %d", seriesNum))

		for _, e := range episodes {
			sb.WithEpisode(e)
		}

		tvb.WithSeries(sb)
	}

	return tvb.Build()
}

func (db *TVDB) getEpisodes(seriesID uint64) []episode {

	var results []episode

	nextPage := 1
	req := &db.requestTemplate

	if link, err := url.Parse(fmt.Sprintf("%s%s", apiBase, fmt.Sprintf(apiEpisodesBySeriesID, seriesID))); err == nil {
		req.URL = link
	} else {
		return nil
	}

	for {
		if nextPage == 0 {
			break
		}

		q := req.URL.Query()
		q.Set("page", strconv.Itoa(nextPage))
		req.URL.RawQuery = q.Encode()

		resp, err := db.httpCli.Do(req)
		if err != nil {
			return nil //if error, discard all
		}

		var episodeResults *tvSeriesEpisodes = &tvSeriesEpisodes{}
		err = dbs.ReadJsonToStruct(resp.Body, &episodeResults)
		if err != nil {
			fmt.Println(err.Error())
			return nil //if error, discard all
		}

		results = append(results, episodeResults.Episodes...)
		nextPage = episodeResults.Links.Next
	}

	return results
}
