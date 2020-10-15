package tvdb

//login structs
type auth struct {
	APIKey   string `json:"apikey"`
	Username string `json:"username"`
	UserKey  string `json:"userkey"`
}

type token struct {
	Token string `json:"token"`
}

//tv structs
type tvSearch struct {
	Results []tvSearchResult `json:"data"`
}
type tvSearchResult struct {
	Aliases    []string `json:"aliases"`
	Banner     string   `json:"banner"`
	FirstAired string   `json:"firstAired"`
	ID         uint64   `json:"id"`
	Image      string   `json:"image"`
	Network    string   `json:"network"`
	Overview   string   `json:"overview"`
	Poster     string   `json:"poster"`
	SeriesName string   `json:"seriesName"`
	Slug       string   `json:"slug"`
	Status     string   `json:"status"`
}

type tv struct {
	Show   tvShow `json:"data"`
	Errors Errors `json:"errors"`
}
type tvShow struct {
	Added           string   `json:"added"`
	AirsDayOfWeek   string   `json:"airsDayOfWeek"`
	AirsTime        string   `json:"airsTime"`
	Aliases         []string `json:"aliases"`
	Banner          string   `json:"banner"`
	FirstAired      string   `json:"firstAired"`
	Genre           []string `json:"genre"`
	ID              uint64   `json:"id"`
	ImdbID          string   `json:"imdbId"`
	LastUpdated     int      `json:"lastUpdated"`
	Network         string   `json:"network"`
	NetworkID       string   `json:"networkId"`
	Overview        string   `json:"overview"`
	Rating          string   `json:"rating"`
	Runtime         string   `json:"runtime"`
	SeriesID        string   `json:"seriesId"`
	SeriesName      string   `json:"seriesName"`
	SiteRating      int      `json:"siteRating"`
	SiteRatingCount int      `json:"siteRatingCount"`
	Slug            string   `json:"slug"`
	Status          string   `json:"status"`
	Zap2ItID        string   `json:"zap2itId"`
}

type tvSeriesEpisodes struct {
	Episodes []episode `json:"data"`
	Errors   Errors    `json:"errors"`
	Links    Links     `json:"links"`
}

type episode struct {
	AbsoluteNumber     int      `json:"absoluteNumber"`
	AiredEpisodeNumber int      `json:"airedEpisodeNumber"`
	AiredSeason        int      `json:"airedSeason"`
	AirsAfterSeason    int      `json:"airsAfterSeason"`
	AirsBeforeEpisode  int      `json:"airsBeforeEpisode"`
	AirsBeforeSeason   int      `json:"airsBeforeSeason"`
	Director           string   `json:"director"`
	Directors          []string `json:"directors"`
	DvdChapter         int      `json:"dvdChapter"`
	DvdDiscid          string   `json:"dvdDiscid"`
	DvdEpisodeNumber   int      `json:"dvdEpisodeNumber"`
	DvdSeason          int      `json:"dvdSeason"`
	EpisodeName        string   `json:"episodeName"`
	Filename           string   `json:"filename"`
	FirstAired         string   `json:"firstAired"`
	GuestStars         []string `json:"guestStars"`
	ID                 uint64   `json:"id"`
	ImdbID             string   `json:"imdbId"`
	LastUpdated        int      `json:"lastUpdated"`
	LastUpdatedBy      uint64   `json:"lastUpdatedBy"`
	Overview           string   `json:"overview"`
	ProductionCode     string   `json:"productionCode"`
	SeriesID           uint64   `json:"seriesId"`
	ShowURL            string   `json:"showUrl"`
	SiteRating         float64  `json:"siteRating"`
	SiteRatingCount    int      `json:"siteRatingCount"`
	ThumbAdded         string   `json:"thumbAdded"`
	ThumbAuthor        int      `json:"thumbAuthor"`
	ThumbHeight        string   `json:"thumbHeight"`
	ThumbWidth         string   `json:"thumbWidth"`
	Writers            []string `json:"writers"`
}
type Errors struct {
	InvalidFilters     []string `json:"invalidFilters"`
	InvalidLanguage    string   `json:"invalidLanguage"`
	InvalidQueryParams []string `json:"invalidQueryParams"`
}
type Links struct {
	First    int `json:"first"`
	Last     int `json:"last"`
	Next     int `json:"next"`
	Previous int `json:"previous"`
}
