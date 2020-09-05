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
	Page    int              `json:"page"`
	TotalResults int         `json:"total_results"`
	TotalPages   int         `json:"total_pages"`
	Results []tvSearchResult `json:"results"`
}

type tvSearchResult struct {
	ID               int      `json:"id"`
	Name             string   `json:"name"`
	OriginalName     string   `json:"original_name"`
}

type tv struct {
	ID               int      `json:"id"`
	Name             string      `json:"name"`
	NumberOfEpisodes    int      `json:"number_of_episodes"`
	NumberOfSeasons     int      `json:"number_of_seasons"`
	Seasons []tvSeries `json:"seasons"`
}

type tvSeries struct {
	AirDate      string `json:"air_date"`
	EpisodeCount int    `json:"episode_count"`
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Overview     string `json:"overview"`
	PosterPath   string `json:"poster_path"`
	SeasonNumber int    `json:"season_number"`
}

type series struct {
	AirDate  string `json:"air_date"`
	Episodes []seriesEpisode `json:"episodes"`
	Name         string `json:"name"`
	Overview     string `json:"overview"`
	ID           int    `json:"id"`
	PosterPath   string `json:"poster_path"`
	SeasonNumber int    `json:"season_number"`
}

type seriesEpisode struct {
	AirDate        string  `json:"air_date"`
	EpisodeNumber  int     `json:"episode_number"`
	ID             int     `json:"id"`
	Name           string  `json:"name"`
	Overview       string  `json:"overview"`
	ProductionCode string  `json:"production_code"`
	SeasonNumber   int     `json:"season_number"`
}