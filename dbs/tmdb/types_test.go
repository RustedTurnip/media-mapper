package tmdb

import (
	"testing"
	"time"

	"github.com/kylelemons/godebug/pretty"
	"github.com/rustedturnip/media-mapper/types"
)

func TestTMDB_buildMovie(t *testing.T) {

	var tests = []struct {
		name     string
		input    movieSearchResult
		expected types.Movie
	}{
		{
			name: "Normal Movie",
			input: movieSearchResult{
				Title:       "Requiem for a Dream",
				ReleaseDate: "2000-10-06",
			},
			expected: types.Movie{
				Title:       "Requiem for a Dream",
				ReleaseDate: time.Unix(970790400, 0).UTC(), //2000-10-06
			},
		},
	}

	for _, test := range tests {
		result := buildMovie(test.input)

		if diff := pretty.Compare(test.expected, *result); diff != "" {
			t.Errorf("%s unexpected diff (-want +got):\n%s", test.name, diff)
		}
	}
}

func TestTMDB_buildTV(t *testing.T) {

	var tests = []struct {
		name     string
		input    tvShow
		expected types.TV
	}{
		{
			name: "Normal, multi-series TV Show",
			input: tvShow{
				ID:               47665,
				Name:             "Black Sails",
				NumberOfSeasons:  2,
				NumberOfEpisodes: 4,
				Seasons: []*tvShowSeriesInfo{
					{
						Name:         "Season 1",
						SeasonNumber: 1,
						EpisodeCount: 2,
						SeasonData: &tvShowSeriesData{
							Name:         "Season 1",
							SeasonNumber: 1,
							Episodes: []*tvShowSeriesEpisode{
								{
									Name:          "Episode 1",
									EpisodeNumber: 1,
									SeasonNumber:  1,
								},
								{
									Name:          "Episode 2",
									EpisodeNumber: 2,
									SeasonNumber:  1,
								},
							},
						},
					},
					{
						Name:         "Season 2",
						SeasonNumber: 2,
						EpisodeCount: 2,
						SeasonData: &tvShowSeriesData{
							Name:         "Season 2",
							SeasonNumber: 2,
							Episodes: []*tvShowSeriesEpisode{
								{
									Name:          "Episode 1",
									EpisodeNumber: 1,
									SeasonNumber:  2,
								},
								{
									Name:          "Episode 2",
									EpisodeNumber: 2,
									SeasonNumber:  1,
								},
							},
						},
					},
				},
			},
			expected: types.TV{
				Title:       "Black Sails",
				SeriesCount: 2,
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
					2: {
						Title:  "Season 2",
						Number: 2,
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
	}

	for _, test := range tests {

		result := buildTV(&test.input)
		if diff := pretty.Compare(test.expected, *result); diff != "" {
			t.Errorf("%s unexpected diff (-want +got):\n%s", test.name, diff)
		}
	}
}
