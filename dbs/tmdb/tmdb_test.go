package tmdb

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/kylelemons/godebug/pretty"
	"github.com/rustedturnip/media-mapper/dbs"
	"github.com/rustedturnip/media-mapper/types"
)

const (
	testAPIToken = "TEST_TOKEN"
)

func TestTMDB_SearchMovies(t *testing.T) {

	var tests = []struct {
		name       string
		titleInput string
		responses  map[string]*http.Response
		expected   []*types.Movie
	}{
		{
			name:       "Normal Movie Search - 1 Result",
			titleInput: "Requiem for a Dream",
			responses: map[string]*http.Response{
				"https://api.themoviedb.org/3/search/movie?api_key=TEST_TOKEN&language=en-GB&query=Requiem+for+a+Dream&page=1&include_adult=true": {
					StatusCode: http.StatusOK,
					Body: ioutil.NopCloser(bytes.NewBufferString(`{
    "page": 1,
    "total_results": 1,
    "total_pages": 1,
    "results": [
        {
            "popularity": 8.806,
            "vote_count": 6623,
            "video": false,
            "poster_path": "/nOd6vjEmzCT0k4VYqsA2hwyi87C.jpg",
            "id": 641,
            "adult": false,
            "backdrop_path": "/c5g1Dn1tF22CS2oOvHDNKr1Ve2U.jpg",
            "original_language": "en",
            "original_title": "Requiem for a Dream",
            "genre_ids": [
                80,
                18,
                53
            ],
            "title": "Requiem for a Dream",
            "vote_average": 8,
            "overview": "The hopes and dreams of four ambitious people are shattered when their drug addictions begin spiraling out of control. A look into addiction and how it overcomes the mind and body.",
            "release_date": "2000-10-06"
        }
    ]
}`)),
				},
			},
			expected: []*types.Movie{
				{
					Title:       "Requiem for a Dream",
					ReleaseDate: time.Unix(970790400, 0).UTC(), //2000-10-06
				},
			},
		},
		{
			name:       "Normal Movie Search - Many Results",
			titleInput: "Lord of the Rings",
			responses: map[string]*http.Response{
				"https://api.themoviedb.org/3/search/movie?api_key=TEST_TOKEN&language=en-GB&query=Lord+of+the+Rings&page=1&include_adult=true": {
					StatusCode: http.StatusOK,
					Body: ioutil.NopCloser(bytes.NewBufferString(`{
    "page": 1,
    "total_results": 3,
    "total_pages": 1,
    "results": [
        {
            "popularity": 51.45,
            "vote_count": 15426,
            "video": false,
            "poster_path": "/5VTN0pR8gcqV3EPUHHfMGnJYN9L.jpg",
            "id": 121,
            "adult": false,
            "backdrop_path": "/9BUvLUz1GhbNpcyQRyZm1HNqMq4.jpg",
            "original_language": "en",
            "original_title": "The Lord of the Rings: The Two Towers",
            "genre_ids": [
                28,
                12,
                14
            ],
            "title": "The Lord of the Rings: The Two Towers",
            "vote_average": 8.3,
            "overview": "Frodo and Sam are trekking to Mordor to destroy the One Ring of Power while Gimli, Legolas and Aragorn search for the orc-captured Merry and Pippin. All along, nefarious wizard Saruman awaits the Fellowship members at the Orthanc Tower in Isengard.",
            "release_date": "2002-12-18"
        },
        {
            "popularity": 52.865,
            "id": 122,
            "video": false,
            "vote_count": 16391,
            "vote_average": 8.5,
            "title": "The Lord of the Rings: The Return of the King",
            "release_date": "2003-12-01",
            "original_language": "en",
            "original_title": "The Lord of the Rings: The Return of the King",
            "genre_ids": [
                12,
                14,
                28
            ],
            "backdrop_path": "/9DeGfFIqjph5CBFVQrD6wv9S7rR.jpg",
            "adult": false,
            "overview": "Aragorn is revealed as the heir to the ancient kings as he, Gandalf and the other members of the broken fellowship struggle to save Gondor from Sauron's forces. Meanwhile, Frodo and Sam take the ring closer to the heart of Mordor, the dark lord's realm.",
            "poster_path": "/rCzpDGLbOoPwLjy3OAm5NUPOTrC.jpg"
        },
        {
            "popularity": 53.291,
            "vote_count": 17860,
            "video": false,
            "poster_path": "/6oom5QYQ2yQTMJIbnvbkBL9cHo6.jpg",
            "id": 120,
            "adult": false,
            "backdrop_path": "/vRQnzOn4HjIMX4LBq9nHhFXbsSu.jpg",
            "original_language": "en",
            "original_title": "The Lord of the Rings: The Fellowship of the Ring",
            "genre_ids": [
                28,
                12,
                14
            ],
            "title": "The Lord of the Rings: The Fellowship of the Ring",
            "vote_average": 8.3,
            "overview": "Young hobbit Frodo Baggins, after inheriting a mysterious ring from his uncle Bilbo, must leave his home in order to keep it from falling into the hands of its evil creator. Along the way, a fellowship is formed to protect the ringbearer and make sure that the ring arrives at its final destination: Mt. Doom, the only place where it can be destroyed.",
            "release_date": "2001-12-18"
        }
    ]
}`)),
				},
			},
			expected: []*types.Movie{
				{
					Title:       "The Lord of the Rings: The Two Towers",
					ReleaseDate: time.Unix(1040169600, 0).UTC(),
				},
				{
					Title:       "The Lord of the Rings: The Return of the King",
					ReleaseDate: time.Unix(1070236800, 0).UTC(),
				},
				{
					Title:       "The Lord of the Rings: The Fellowship of the Ring",
					ReleaseDate: time.Unix(1008633600, 0).UTC(),
				},
			},
		},
	}

	for _, test := range tests {

		//create db instance with mocked http client
		db := TMDB{
			apiKey:     testAPIToken,
			httpClient: dbs.NewHttpClient(test.responses),
		}

		results := db.SearchMovies(test.titleInput)

		if diff := pretty.Compare(test.expected, results); diff != "" {
			t.Errorf("%s unexpected diff (-want +got):\n%s", test.name, diff)
		}
	}
}

func TestTMDB_SearchTV(t *testing.T) {

	var tests = []struct {
		name       string
		titleInput string
		responses  map[string]*http.Response
		expected   []*types.TV
	}{
		{
			name:       "Normal TV Search - Multiple Series, Multiple Episodes",
			titleInput: "Paradise PD",
			responses: map[string]*http.Response{
				"https://api.themoviedb.org/3/search/tv?api_key=TEST_TOKEN&language=en-GB&query=Paradise+PD&page=1&include_adult=true": {
					StatusCode: http.StatusOK,
					Body: ioutil.NopCloser(bytes.NewBufferString(`{
    "page": 1,
    "total_results": 1,
    "total_pages": 1,
    "results": [
        {
            "original_name": "Paradise PD",
            "genre_ids": [
                16,
                35
            ],
            "name": "Paradise PD",
            "popularity": 15.595,
            "origin_country": [
                "US"
            ],
            "vote_count": 95,
            "first_air_date": "2018-08-31",
            "backdrop_path": "/vVlhy5xJPHTJ0pMprsI0zxbrrpM.jpg",
            "original_language": "en",
            "id": 81983,
            "vote_average": 8,
            "overview": "An eager young rookie joins the ragtag small-town police force led by his dad as they bumble, squabble and snort their way through a big drug case.",
            "poster_path": "/desSj4kx0y9p61vm9QBE3Wm8GxK.jpg"
        }
    ]
}`)),
				},
				"https://api.themoviedb.org/3/tv/81983?api_key=TEST_TOKEN&language=en-GB": {
					StatusCode: http.StatusOK,
					Body: ioutil.NopCloser(bytes.NewBufferString(`{
    "backdrop_path": "/vVlhy5xJPHTJ0pMprsI0zxbrrpM.jpg",
    "created_by": [
        {
            "id": 948351,
            "credit_id": "5b875db60e0a266f34004f48",
            "name": "Roger Black",
            "gender": 2,
            "profile_path": null
        },
        {
            "id": 1094963,
            "credit_id": "5b875dbe0e0a266f2b005963",
            "name": "Waco O'Guin",
            "gender": 2,
            "profile_path": null
        }
    ],
    "episode_run_time": [
        28
    ],
    "first_air_date": "2018-08-31",
    "genres": [
        {
            "id": 16,
            "name": "Animation"
        },
        {
            "id": 35,
            "name": "Comedy"
        }
    ],
    "homepage": "https://www.netflix.com/title/80191522",
    "id": 81983,
    "in_production": true,
    "languages": [
        "en"
    ],
    "last_air_date": "2020-03-06",
    "last_episode_to_air": {
        "air_date": "2020-03-06",
        "episode_number": 8,
        "id": 2182163,
        "name": "Operation DD",
        "overview": "After learning that Fitz isn't quite what he seems, the squad races to stop their real enemy — and save Paradise from a nuclear disaster.",
        "production_code": "",
        "season_number": 2,
        "show_id": 81983,
        "still_path": "/5sBIxsBbpbtJy29EnK7mhAEfBml.jpg",
        "vote_average": 0.0,
        "vote_count": 0
    },
    "name": "Paradise PD",
    "next_episode_to_air": null,
    "networks": [
        {
            "name": "Netflix",
            "id": 213,
            "logo_path": "/wwemzKWzjKYJFfCeiB57q3r4Bcm.png",
            "origin_country": ""
        }
    ],
    "number_of_episodes": 18,
    "number_of_seasons": 2,
    "origin_country": [
        "US"
    ],
    "original_language": "en",
    "original_name": "Paradise PD",
    "overview": "An eager young rookie joins the ragtag small-town police force led by his dad as they bumble, squabble and snort their way through a big drug case.",
    "popularity": 15.595,
    "poster_path": "/desSj4kx0y9p61vm9QBE3Wm8GxK.jpg",
    "production_companies": [
        {
            "id": 86647,
            "logo_path": null,
            "name": "Odenkirk Provissiero Entertainment",
            "origin_country": ""
        },
        {
            "id": 30452,
            "logo_path": "/zmU1ElCS02iL5N7E5MuY4fV7bCX.png",
            "name": "Bento Box Entertainment",
            "origin_country": "US"
        }
    ],
    "seasons": [
        {
            "air_date": "2018-08-31",
            "episode_count": 10,
            "id": 108605,
            "name": "Season 1",
            "overview": "",
            "poster_path": "/desSj4kx0y9p61vm9QBE3Wm8GxK.jpg",
            "season_number": 1
        },
        {
            "air_date": "2020-03-06",
            "episode_count": 8,
            "id": 143701,
            "name": "Season 2",
            "overview": "",
            "poster_path": "/ij4M1eGTJHU4UOhqGKQfAXNWxDC.jpg",
            "season_number": 2
        }
    ],
    "status": "Returning Series",
    "type": "Scripted",
    "vote_average": 8.0,
    "vote_count": 95
}`)),
				},
				"https://api.themoviedb.org/3/tv/81983/season/1?api_key=TEST_TOKEN&language=en-GB": {
					StatusCode: http.StatusOK,
					Body: ioutil.NopCloser(bytes.NewBufferString(`{
    "_id": "5b875d8b9251412d49004f82",
    "air_date": "2018-08-31",
    "episodes": [
        {
            "air_date": "2018-08-31",
            "episode_number": 1,
            "id": 1560627,
            "name": "Welcome to Paradise",
            "overview": "At 18, Kevin Crawford finally gets a shot at joining the police force run by his dad, just as a new drug dubbed \"argyle meth\" hits the streets.",
            "production_code": "",
            "season_number": 1,
            "show_id": 81983,
            "still_path": "/nUJwHB4aXvRe2GNRhSbkGRTXfz0.jpg",
            "vote_average": 0.0,
            "vote_count": 0,
            "crew": [],
            "guest_stars": []
        },
        {
            "air_date": "2018-08-31",
            "episode_number": 2,
            "id": 1561112,
            "name": "Ass on the Line",
            "overview": "Bullet finds fame and glory in an underground dogfighting ring, and Chief Crawford butts heads with his biggest rival on a maddening homicide case.",
            "production_code": "",
            "season_number": 1,
            "show_id": 81983,
            "still_path": "/cayB67rEdZeqrQwiaI0uhzHuGGk.jpg",
            "vote_average": 0.0,
            "vote_count": 0,
            "crew": [],
            "guest_stars": []
        },
        {
            "air_date": "2018-08-31",
            "episode_number": 3,
            "id": 1561113,
            "name": "Black & Blue",
            "overview": "At Gina's insistence, Fitz starts carrying a gun -- and ignites a national media scandal. Shipped off to a nursing home, Hopson uncovers a conspiracy.",
            "production_code": "",
            "season_number": 1,
            "show_id": 81983,
            "still_path": "/dvH7qAu1dHOedJTeZvwyN22diNj.jpg",
            "vote_average": 0.0,
            "vote_count": 0,
            "crew": [],
            "guest_stars": []
        },
        {
            "air_date": "2018-08-31",
            "episode_number": 4,
            "id": 1561114,
            "name": "Karla",
            "overview": "When Kevin's mom buys him a sleek new talking police car, it's love at first sight. Bullet turns Dusty into a fried-chicken kingpin.",
            "production_code": "",
            "season_number": 1,
            "show_id": 81983,
            "still_path": "/1T3OQHcoaNmWjcngfViKezV6v1F.jpg",
            "vote_average": 0.0,
            "vote_count": 0,
            "crew": [],
            "guest_stars": []
        }
    ],
    "name": "Season 1",
    "overview": "",
    "id": 108605,
    "poster_path": "/desSj4kx0y9p61vm9QBE3Wm8GxK.jpg",
    "season_number": 1
}`)),
				},
				"https://api.themoviedb.org/3/tv/81983/season/2?api_key=TEST_TOKEN&language=en-GB": {
					StatusCode: http.StatusOK,
					Body: ioutil.NopCloser(bytes.NewBufferString(`{
    "_id": "5e4e8c5835811d0015509819",
    "air_date": "2020-03-06",
    "episodes": [
        {
            "air_date": "2020-03-06",
            "episode_number": 1,
            "id": 2170009,
            "name": "Paradise Found",
            "overview": "As tourists flock to the new, peaceful Paradise, Gina plots to bust Dusty out of prison, Kevin savors his hero status, and Karen plans an execution.",
            "production_code": "",
            "season_number": 2,
            "show_id": 81983,
            "still_path": "/1yHtLtLBRxv8WcVwKv6Vgx5IOBj.jpg",
            "vote_average": 0.0,
            "vote_count": 0,
            "crew": [],
            "guest_stars": []
        },
        {
            "air_date": "2020-03-06",
            "episode_number": 2,
            "id": 2182152,
            "name": "Big Ball Energy",
            "overview": "On Kevin Sucks Day, Fitz hunts down a new meth supplier, the chief discovers Karen's secret fetish, and Kevin vows to defy an embarrassing prediction.",
            "production_code": "",
            "season_number": 2,
            "show_id": 81983,
            "still_path": "/oGTNDaczIMDKzd9USo3zjlfWo8c.jpg",
            "vote_average": 0.0,
            "vote_count": 0,
            "crew": [],
            "guest_stars": []
        },
        {
            "air_date": "2020-03-06",
            "episode_number": 3,
            "id": 2182153,
            "name": "Tucker Carlson Is a Huge D**k",
            "overview": "A rant by Tucker Carlson sparks a war of the sexes, leaving Paradise with two police forces. Fitz's new evil plan is thwarted by Gal-Qaeda.",
            "production_code": "",
            "season_number": 2,
            "show_id": 81983,
            "still_path": "/d6QscY64YIm55bvjf8YmjOQtHIB.jpg",
            "vote_average": 0.0,
            "vote_count": 0,
            "crew": [],
            "guest_stars": []
        },
        {
            "air_date": "2020-03-06",
            "episode_number": 4,
            "id": 2182154,
            "name": "Who Ate Wally's Waffles",
            "overview": "Dusty finds a long-lost sitcom star living in Paradise and sets out to reboot his career. The squad obsesses over Kevin's bathroom habits.",
            "production_code": "",
            "season_number": 2,
            "show_id": 81983,
            "still_path": "/jYAjNlyjmZrLlUPewhwQ7ZiN6bh.jpg",
            "vote_average": 0.0,
            "vote_count": 0,
            "crew": [],
            "guest_stars": []
        }
    ],
    "name": "Season 2",
    "overview": "",
    "id": 143701,
    "poster_path": "/ij4M1eGTJHU4UOhqGKQfAXNWxDC.jpg",
    "season_number": 2
}`)),
				},
			},
			expected: []*types.TV{
				{
					Title:       "Paradise PD",
					SeriesCount: 2,
					Series: map[int]*types.Series{
						1: {
							Title:  "Season 1",
							Number: 1,
							Episodes: map[int]*types.Episode{
								1: {
									Title:  "Welcome to Paradise",
									Number: 1,
								},
								2: {
									Title:  "Ass on the Line",
									Number: 2,
								},
								3: {
									Title:  "Black & Blue",
									Number: 3,
								},
								4: {
									Title:  "Karla",
									Number: 4,
								},
							},
						},
						2: {
							Title:  "Season 2",
							Number: 2,
							Episodes: map[int]*types.Episode{
								1: {
									Title:  "Paradise Found",
									Number: 1,
								},
								2: {
									Title:  "Big Ball Energy",
									Number: 2,
								},
								3: {
									Title:  "Tucker Carlson Is a Huge D**k",
									Number: 3,
								},
								4: {
									Title:  "Who Ate Wally's Waffles",
									Number: 4,
								},
							},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		//test specific db instance
		db := &TMDB{
			apiKey:     testAPIToken,
			httpClient: dbs.NewHttpClient(test.responses), //mocked http client with test's responses to queries
		}

		//run test
		results := db.SearchTV(test.titleInput)
		if diff := pretty.Compare(test.expected, results); diff != "" {
			t.Errorf("%s unexpected diff (-want +got):\n%s", test.name, diff)
		}
	}
}
