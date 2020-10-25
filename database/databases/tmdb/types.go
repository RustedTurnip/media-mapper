package tmdb

type movieSearch struct {
	Page         int                 `json:"page"`
	TotalResults int                 `json:"total_results"`
	TotalPages   int                 `json:"total_pages"`
	Results      []movieSearchResult `json:"results"`
}

type movieSearchResult struct {
	Popularity       float64 `json:"popularity"`
	VoteCount        int     `json:"vote_count"`
	Video            bool    `json:"video"`
	PosterPath       string  `json:"poster_path"`
	ID               int     `json:"id"`
	Adult            bool    `json:"adult"`
	BackdropPath     string  `json:"backdrop_path"`
	OriginalLanguage string  `json:"original_language"`
	OriginalTitle    string  `json:"original_title"`
	GenreIds         []int   `json:"genre_ids"`
	Title            string  `json:"title"`
	VoteAverage      float64 `json:"vote_average"`
	Overview         string  `json:"overview"`
	ReleaseDate      string  `json:"release_date"`
}

type tvSearch struct {
	Page         int              `json:"page"`
	TotalResults int              `json:"total_results"`
	TotalPages   int              `json:"total_pages"`
	Results      []tvSearchResult `json:"results"`
}

type tvSearchResult struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	OriginalName string `json:"original_name"`
}

type tvShow struct {
	ID               int                 `json:"id"`
	Name             string              `json:"name"`
	NumberOfEpisodes int                 `json:"number_of_episodes"`
	NumberOfSeasons  int                 `json:"number_of_seasons"`
	Seasons          []*tvShowSeriesInfo `json:"seasons"`
}

//info about the series as a whole
type tvShowSeriesInfo struct {
	AirDate      string `json:"air_date"`
	EpisodeCount int    `json:"episode_count"`
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Overview     string `json:"overview"`
	PosterPath   string `json:"poster_path"`
	SeasonNumber int    `json:"season_number"`
	SeasonData   *tvShowSeriesData
}

//data of series, i.e. episodes etc.
type tvShowSeriesData struct {
	AirDate      string                 `json:"air_date"`
	Episodes     []*tvShowSeriesEpisode `json:"episodes"`
	Name         string                 `json:"name"`
	Overview     string                 `json:"overview"`
	ID           int                    `json:"id"`
	PosterPath   string                 `json:"poster_path"`
	SeasonNumber int                    `json:"season_number"`
}

type tvShowSeriesEpisode struct {
	AirDate        string `json:"air_date"`
	EpisodeNumber  int    `json:"episode_number"`
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Overview       string `json:"overview"`
	ProductionCode string `json:"production_code"`
	SeasonNumber   int    `json:"season_number"`
}
