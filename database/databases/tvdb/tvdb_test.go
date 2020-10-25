package tvdb

import (
	"bytes"
	"github.com/rustedturnip/media-mapper/database"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/kylelemons/godebug/pretty"
	"github.com/rustedturnip/media-mapper/types"
)

func TestTVDB_SearchTV(t *testing.T) {
	var tests = []struct {
		name       string
		titleInput string
		expected   []*types.TV
		responses  map[string]*http.Response //map[expectedURL]response
	}{
		{
			name:       "Normal TV Search - 1 Season, Multi-Episodes (1 Page)",
			titleInput: "Taboo",
			expected: []*types.TV{
				{
					Title:       "Taboo (2017)",
					SeriesCount: 1,
					Series: map[int]*types.Series{
						1: {
							Title:  "Season 1",
							Number: 1,
							Episodes: map[int]*types.Episode{
								1: {
									Title:  "Episode 1",
									Number: 1,
								},
								2: {
									Title:  "Episode 2",
									Number: 2,
								},
							},
						},
					},
				},
			},
			responses: map[string]*http.Response{
				//search response
				"https://api.thetvdb.com/search/series?name=Taboo": {
					StatusCode: http.StatusOK,
					Body: ioutil.NopCloser(bytes.NewBufferString(`{
    "data": [
        {
            "aliases": [
                "Taboo"
            ],
            "banner": "/banners/graphical/292157-g3.jpg",
            "firstAired": "2017-1-7",
            "id": 292157,
            "image": "/banners/posters/292157-1.jpg",
            "network": "BBC One",
            "overview": "James Keziah Delaney has been to the ends of the earth and comes back irrevocably changed. Believed to be long dead, he returns home to London from Africa to inherit what is left of his father's shipping empire and rebuild a life for himself. But his father's legacy is a poisoned chalice, and with enemies lurking in every dark corner, James must navigate increasingly complex territories to avoid his own death sentence. Encircled by conspiracy, murder and betrayal, a dark family mystery unfolds in a combustible tale of love and treachery.",
            "poster": "/banners/posters/292157-1.jpg",
            "seriesName": "Taboo (2017)",
            "slug": "taboo-2017",
            "status": "Continuing"
        }
    ]
}`)),
				},
				//seriesByID response
				"https://api.thetvdb.com/series/292157": {
					StatusCode: http.StatusOK,
					Body: ioutil.NopCloser(bytes.NewBufferString(`{
    "data": {
        "id": 292157,
        "seriesId": "",
        "seriesName": "Taboo (2017)",
        "aliases": [],
        "season": "1",
        "poster": "posters/292157-1.jpg",
        "banner": "graphical/292157-g3.jpg",
        "fanart": "fanart/original/292157-2.jpg",
        "status": "Continuing",
        "firstAired": "2017-01-07",
        "network": "BBC One",
        "networkId": "37",
        "runtime": "60",
        "language": "en",
        "genre": [
            "Action",
            "Drama"
        ],
        "overview": "James Keziah Delaney has been to the ends of the earth and comes back irrevocably changed. Believed to be long dead, he returns home to London from Africa to inherit what is left of his father's shipping empire and rebuild a life for himself. But his father's legacy is a poisoned chalice, and with enemies lurking in every dark corner, James must navigate increasingly complex territories to avoid his own death sentence. Encircled by conspiracy, murder and betrayal, a dark family mystery unfolds in a combustible tale of love and treachery.",
        "lastUpdated": 1598923282,
        "airsDayOfWeek": "Saturday",
        "airsTime": "9:15 PM",
        "rating": "TV-14",
        "imdbId": "tt3647998",
        "zap2itId": "",
        "added": "2015-02-16 08:07:08",
        "addedBy": 114544,
        "siteRating": 8.8,
        "siteRatingCount": 887,
        "slug": "taboo-2017"
    }
}`)),
				},

				//Episodes By Series ID response
				"https://api.thetvdb.com/series/292157/episodes?page=1": {
					StatusCode: http.StatusOK,
					Body: ioutil.NopCloser(bytes.NewBufferString(`{
    "links": {
        "first": 1,
        "last": 1,
        "next": null,
        "prev": null
    },
    "data": [
        {
            "id": 5840491,
            "airedSeason": 1,
            "airedSeasonID": 682219,
            "airedEpisodeNumber": 1,
            "episodeName": "Episode 1",
            "firstAired": "2017-01-07",
            "guestStars": [],
            "directors": [
                "Kristoffer Nyholm"
            ],
            "writers": [
                "Steven Knight"
            ],
            "overview": "In the series opener, James Delaney returns to 1814 London after ten years in Africa to claim a mysterious legacy left to him by his father.",
            "language": {
                "episodeName": "en",
                "overview": "en"
            },
            "productionCode": "b088rzbt",
            "showUrl": "",
            "lastUpdated": 1526904215,
            "dvdDiscid": "",
            "dvdSeason": 1,
            "dvdEpisodeNumber": 1,
            "dvdChapter": null,
            "absoluteNumber": 1,
            "filename": "episodes/292157/5840491.jpg",
            "seriesId": 292157,
            "lastUpdatedBy": 1,
            "airsAfterSeason": null,
            "airsBeforeSeason": null,
            "airsBeforeEpisode": null,
            "imdbId": "tt4700596",
            "contentRating": "TV-14",
            "thumbAuthor": 1,
            "thumbAdded": "2019-11-13 11:41:13",
            "thumbWidth": "640",
            "thumbHeight": "360",
            "siteRating": 8.8,
            "siteRatingCount": 419,
            "isMovie": 0
        },
        {
            "id": 5840501,
            "airedSeason": 1,
            "airedSeasonID": 682219,
            "airedEpisodeNumber": 2,
            "episodeName": "Episode 2",
            "firstAired": "2017-01-14",
            "guestStars": [],
            "directors": [
                "Kristoffer Nyholm"
            ],
            "writers": [
                "Steven Knight"
            ],
            "overview": "As James Delaney assembles his league of the damned; an unexpected arrival threatens to disrupt his plans. ",
            "language": {
                "episodeName": "en",
                "overview": "en"
            },
            "productionCode": "b08bkqg8",
            "showUrl": "",
            "lastUpdated": 1574190700,
            "dvdDiscid": "",
            "dvdSeason": 1,
            "dvdEpisodeNumber": 2,
            "dvdChapter": null,
            "absoluteNumber": 2,
            "filename": "episodes/292157/5840501.jpg",
            "seriesId": 292157,
            "lastUpdatedBy": 118266,
            "airsAfterSeason": null,
            "airsBeforeSeason": null,
            "airsBeforeEpisode": null,
            "imdbId": "tt4700604",
            "contentRating": "TV-14",
            "thumbAuthor": 1,
            "thumbAdded": "2019-11-13 11:41:13",
            "thumbWidth": "640",
            "thumbHeight": "360",
            "siteRating": 8.7,
            "siteRatingCount": 363,
            "isMovie": 0
        }
    ]
}`)),
				},
			},
		},
		{
			name:       "Normal TV Search - Multi-Season, Multi-Episodes (Multiple Pages)",
			titleInput: "simpsons",
			expected: []*types.TV{
				{
					Title:       "The Simpsons",
					SeriesCount: 3,
					Series: map[int]*types.Series{
						1: {
							Title:  "Season 1",
							Number: 1,
							Episodes: map[int]*types.Episode{
								1: {
									Title:  "Simpsons Roasting on an Open Fire",
									Number: 1,
								},
								2: {
									Title:  "Bart the Genius",
									Number: 2,
								},
							},
						},
						2: {
							Title:  "Season 2",
							Number: 2,
							Episodes: map[int]*types.Episode{
								1: {
									Title:  "Bart Gets an \"F\"",
									Number: 1,
								},
								2: {
									Title:  "Simpson and Delilah",
									Number: 2,
								},
							},
						},
						3: {
							Title:  "Season 3",
							Number: 3,
							Episodes: map[int]*types.Episode{
								1: {
									Title:  "Stark Raving Dad",
									Number: 1,
								},
								2: {
									Title:  "Mr. Lisa Goes to Washington",
									Number: 2,
								},
							},
						},
					},
				},
			},
			responses: map[string]*http.Response{
				//search response
				"https://api.thetvdb.com/search/series?name=simpsons": {
					StatusCode: http.StatusOK,
					Body: ioutil.NopCloser(bytes.NewBufferString(`{
    "data": [
        {
            "aliases": [
                "심슨"
            ],
            "banner": "/banners/graphical/71663-g13.jpg",
            "firstAired": "1987-4-19",
            "id": 71663,
            "image": "/banners/posters/71663-15.jpg",
            "network": "FOX",
            "overview": "Set in Springfield, the average American town, the show focuses on the antics and everyday adventures of the Simpson family; Homer, Marge, Bart, Lisa and Maggie, as well as a virtual cast of thousands. Since the beginning, the series has been a pop culture icon, attracting hundreds of celebrities to guest star. The show has also made name for itself in its fearless satirical take on politics, media and American life in general.",
            "poster": "/banners/posters/71663-15.jpg",
            "seriesName": "The Simpsons",
            "slug": "the-simpsons",
            "status": "Continuing"
        }
    ]
}`)),
				},
				//seriesByID response
				"https://api.thetvdb.com/series/71663": {
					StatusCode: http.StatusOK,
					Body: ioutil.NopCloser(bytes.NewBufferString(`{
    "data": {
        "id": 71663,
        "seriesId": "146",
        "seriesName": "The Simpsons",
        "aliases": [],
        "season": "32",
        "poster": "posters/71663-15.jpg",
        "banner": "graphical/71663-g13.jpg",
        "fanart": "fanart/original/71663-10.jpg",
        "status": "Continuing",
        "firstAired": "1987-04-19",
        "network": "Disney+",
        "networkId": "1306",
        "runtime": "25",
        "language": "en",
        "genre": [
            "Animation",
            "Comedy"
        ],
        "overview": "Set in Springfield, the average American town, the show focuses on the antics and everyday adventures of the Simpson family; Homer, Marge, Bart, Lisa and Maggie, as well as a virtual cast of thousands. Since the beginning, the series has been a pop culture icon, attracting hundreds of celebrities to guest star. The show has also made name for itself in its fearless satirical take on politics, media and American life in general.",
        "lastUpdated": 1603397622,
        "airsDayOfWeek": "Sunday",
        "airsTime": "8:00 PM",
        "rating": "TV-PG",
        "imdbId": "tt0096697",
        "zap2itId": "EP00018693",
        "added": "2008-02-04 00:00:00",
        "addedBy": 1,
        "siteRating": 8.9,
        "siteRatingCount": 24136,
        "slug": "the-simpsons"
    }
}`)),
				},

				//Episodes By Series ID response - page 1
				"https://api.thetvdb.com/series/71663/episodes?page=1": {
					StatusCode: http.StatusOK,
					Body: ioutil.NopCloser(bytes.NewBufferString(`{
  "links": {
    "first": 1,
    "last": 3,
    "next": 2,
    "prev": null
  },
  "data": [
    {
      "id": 55452,
      "airedSeason": 1,
      "airedSeasonID": 2727,
      "airedEpisodeNumber": 1,
      "episodeName": "Simpsons Roasting on an Open Fire",
      "firstAired": "1989-12-17",
      "guestStars": [
        "Christopher Collins"
      ],
      "directors": [
        "David Silverman"
      ],
      "writers": [
        "Mimi Pond"
      ],
      "overview": "When his Christmas bonus is cancelled, Homer becomes a department-store Santa--and then bets his meager earnings at the track. When all seems lost, Homer and Bart save Christmas by adopting the losing greyhound, Santa's Little Helper.",
      "language": {
        "episodeName": "en",
        "overview": "en"
      },
      "productionCode": "7G08",
      "showUrl": "",
      "lastUpdated": 1601811544,
      "dvdDiscid": "",
      "dvdSeason": 1,
      "dvdEpisodeNumber": 1,
      "dvdChapter": null,
      "absoluteNumber": 1,
      "filename": "episodes/71663/55452.jpg",
      "seriesId": 71663,
      "lastUpdatedBy": 2279861,
      "airsAfterSeason": null,
      "airsBeforeSeason": null,
      "airsBeforeEpisode": null,
      "imdbId": "tt0348034",
      "contentRating": "TV-PG",
      "thumbAuthor": 1,
      "thumbAdded": "2019-11-13 10:45:28",
      "thumbWidth": "640",
      "thumbHeight": "360",
      "siteRating": 7.3,
      "siteRatingCount": 2406,
      "isMovie": 0
    },
	{
		"id": 55453,
		"airedSeason": 1,
		"airedSeasonID": 2727,
		"airedEpisodeNumber": 2,
		"episodeName": "Bart the Genius",
		"firstAired": "1990-01-14",
		"guestStars": [
			"Marcia Wallace"
		],
		"directors": [
			"David Silverman"
		],
		"writers": [
			"Jon Vitti"
		],
		"overview": "After switching IQ tests with Martin, Bart is mistaken for a child genius. When he's enrolled in a school for gifted students, a series of embarrassments and mishaps makes him long for his old life.",
		"language": {
			"episodeName": "en",
			"overview": "en"
		},
		"productionCode": "7G02",
		"showUrl": "",
		"lastUpdated": 1601811833,
		"dvdDiscid": "",
		"dvdSeason": 1,
		"dvdEpisodeNumber": 2,
		"dvdChapter": null,
		"absoluteNumber": 2,
		"filename": "episodes/71663/55453.jpg",
		"seriesId": 71663,
		"lastUpdatedBy": 2279861,
		"airsAfterSeason": null,
		"airsBeforeSeason": null,
		"airsBeforeEpisode": null,
		"imdbId": "tt0756593",
		"contentRating": "TV-PG",
		"thumbAuthor": 1,
		"thumbAdded": "2019-11-13 10:45:28",
		"thumbWidth": "640",
		"thumbHeight": "360",
		"siteRating": 7.2,
		"siteRatingCount": 1854,
		"isMovie": 0
	}
  ]
}`)),
				},
				//Episodes By Series ID response - page 2
				"https://api.thetvdb.com/series/71663/episodes?page=2": {
					StatusCode: http.StatusOK,
					Body: ioutil.NopCloser(bytes.NewBufferString(`{
  "links": {
    "first": 1,
    "last": 3,
    "next": 3,
    "prev": 1
  },
  "data": [
    {
      "id": 55452,
      "airedSeason": 2,
      "airedSeasonID": 2728,
      "airedEpisodeNumber": 1,
      "episodeName": "Bart Gets an \"F\"",
      "firstAired": "1989-12-17",
      "guestStars": [
        "Christopher Collins"
      ],
      "directors": [
        "David Silverman"
      ],
      "writers": [
        "Mimi Pond"
      ],
      "overview": "When his Christmas bonus is cancelled, Homer becomes a department-store Santa--and then bets his meager earnings at the track. When all seems lost, Homer and Bart save Christmas by adopting the losing greyhound, Santa's Little Helper.",
      "language": {
        "episodeName": "en",
        "overview": "en"
      },
      "productionCode": "7G08",
      "showUrl": "",
      "lastUpdated": 1601811544,
      "dvdDiscid": "",
      "dvdSeason": 1,
      "dvdEpisodeNumber": 1,
      "dvdChapter": null,
      "absoluteNumber": 1,
      "filename": "episodes/71663/55452.jpg",
      "seriesId": 71663,
      "lastUpdatedBy": 2279861,
      "airsAfterSeason": null,
      "airsBeforeSeason": null,
      "airsBeforeEpisode": null,
      "imdbId": "tt0348034",
      "contentRating": "TV-PG",
      "thumbAuthor": 1,
      "thumbAdded": "2019-11-13 10:45:28",
      "thumbWidth": "640",
      "thumbHeight": "360",
      "siteRating": 7.3,
      "siteRatingCount": 2406,
      "isMovie": 0
    },
	{
		"id": 55453,
		"airedSeason": 2,
		"airedSeasonID": 2728,
		"airedEpisodeNumber": 2,
		"episodeName": "Simpson and Delilah",
		"firstAired": "1990-01-14",
		"guestStars": [
			"Marcia Wallace"
		],
		"directors": [
			"David Silverman"
		],
		"writers": [
			"Jon Vitti"
		],
		"overview": "After switching IQ tests with Martin, Bart is mistaken for a child genius. When he's enrolled in a school for gifted students, a series of embarrassments and mishaps makes him long for his old life.",
		"language": {
			"episodeName": "en",
			"overview": "en"
		},
		"productionCode": "7G02",
		"showUrl": "",
		"lastUpdated": 1601811833,
		"dvdDiscid": "",
		"dvdSeason": 1,
		"dvdEpisodeNumber": 2,
		"dvdChapter": null,
		"absoluteNumber": 2,
		"filename": "episodes/71663/55453.jpg",
		"seriesId": 71663,
		"lastUpdatedBy": 2279861,
		"airsAfterSeason": null,
		"airsBeforeSeason": null,
		"airsBeforeEpisode": null,
		"imdbId": "tt0756593",
		"contentRating": "TV-PG",
		"thumbAuthor": 1,
		"thumbAdded": "2019-11-13 10:45:28",
		"thumbWidth": "640",
		"thumbHeight": "360",
		"siteRating": 7.2,
		"siteRatingCount": 1854,
		"isMovie": 0
	}
  ]
}`)),
				},
				//Episodes By Series ID response - page 3
				"https://api.thetvdb.com/series/71663/episodes?page=3": {
					StatusCode: http.StatusOK,
					Body: ioutil.NopCloser(bytes.NewBufferString(`{
  "links": {
    "first": 1,
    "last": 3,
    "next": null,
    "prev": 2
  },
  "data": [
    {
      "id": 55452,
      "airedSeason": 3,
      "airedSeasonID": 2729,
      "airedEpisodeNumber": 1,
      "episodeName": "Stark Raving Dad",
      "firstAired": "1989-12-17",
      "guestStars": [
        "Christopher Collins"
      ],
      "directors": [
        "David Silverman"
      ],
      "writers": [
        "Mimi Pond"
      ],
      "overview": "When his Christmas bonus is cancelled, Homer becomes a department-store Santa--and then bets his meager earnings at the track. When all seems lost, Homer and Bart save Christmas by adopting the losing greyhound, Santa's Little Helper.",
      "language": {
        "episodeName": "en",
        "overview": "en"
      },
      "productionCode": "7G08",
      "showUrl": "",
      "lastUpdated": 1601811544,
      "dvdDiscid": "",
      "dvdSeason": 1,
      "dvdEpisodeNumber": 1,
      "dvdChapter": null,
      "absoluteNumber": 1,
      "filename": "episodes/71663/55452.jpg",
      "seriesId": 71663,
      "lastUpdatedBy": 2279861,
      "airsAfterSeason": null,
      "airsBeforeSeason": null,
      "airsBeforeEpisode": null,
      "imdbId": "tt0348034",
      "contentRating": "TV-PG",
      "thumbAuthor": 1,
      "thumbAdded": "2019-11-13 10:45:28",
      "thumbWidth": "640",
      "thumbHeight": "360",
      "siteRating": 7.3,
      "siteRatingCount": 2406,
      "isMovie": 0
    },
	{
		"id": 55453,
		"airedSeason": 3,
		"airedSeasonID": 2729,
		"airedEpisodeNumber": 2,
		"episodeName": "Mr. Lisa Goes to Washington",
		"firstAired": "1990-01-14",
		"guestStars": [
			"Marcia Wallace"
		],
		"directors": [
			"David Silverman"
		],
		"writers": [
			"Jon Vitti"
		],
		"overview": "After switching IQ tests with Martin, Bart is mistaken for a child genius. When he's enrolled in a school for gifted students, a series of embarrassments and mishaps makes him long for his old life.",
		"language": {
			"episodeName": "en",
			"overview": "en"
		},
		"productionCode": "7G02",
		"showUrl": "",
		"lastUpdated": 1601811833,
		"dvdDiscid": "",
		"dvdSeason": 1,
		"dvdEpisodeNumber": 2,
		"dvdChapter": null,
		"absoluteNumber": 2,
		"filename": "episodes/71663/55453.jpg",
		"seriesId": 71663,
		"lastUpdatedBy": 2279861,
		"airsAfterSeason": null,
		"airsBeforeSeason": null,
		"airsBeforeEpisode": null,
		"imdbId": "tt0756593",
		"contentRating": "TV-PG",
		"thumbAuthor": 1,
		"thumbAdded": "2019-11-13 10:45:28",
		"thumbWidth": "640",
		"thumbHeight": "360",
		"siteRating": 7.2,
		"siteRatingCount": 1854,
		"isMovie": 0
	}
  ]
}`)),
				},
			},
		},
	}

	for _, test := range tests {
		//initialise db with test specific mock client with test's responses
		db := TVDB{
			httpClient: database.NewHttpClient(test.responses),
		}

		//test
		result := db.SearchTV(test.titleInput)
		if diff := pretty.Compare(test.expected, result); diff != "" {
			t.Errorf("%s unexpected diff (-want +got):\n%s", test.name, diff)
		}
	}
}
