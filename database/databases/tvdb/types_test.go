package tvdb

import (
	"testing"

	"github.com/kylelemons/godebug/pretty"
	"github.com/rustedturnip/media-mapper/types"
)

func TestTVDB_buildMovie(t *testing.T) {

	var tests = []struct {
		name     string
		input    *tvShow
		expected *types.TV
	}{
		{
			name: "Normal TV Multi-Series",
			input: &tvShow{
				SeriesName: "Line of Duty",
				Series: &tvSeriesEpisodes{
					Episodes: []*episode{
						{
							EpisodeName:        "Episode 1",
							AiredEpisodeNumber: 1,
							AiredSeason:        1,
						},
						{
							EpisodeName:        "Episode 2",
							AiredEpisodeNumber: 2,
							AiredSeason:        1,
						},
						{
							EpisodeName:        "Episode 3",
							AiredEpisodeNumber: 3,
							AiredSeason:        1,
						},
						{
							EpisodeName:        "Episode 1",
							AiredEpisodeNumber: 1,
							AiredSeason:        2,
						},
						{
							EpisodeName:        "Episode 2",
							AiredEpisodeNumber: 2,
							AiredSeason:        2,
						},
						{
							EpisodeName:        "Episode 3",
							AiredEpisodeNumber: 3,
							AiredSeason:        2,
						},
					},
				},
			},
			expected: &types.TV{
				Title:       "Line of Duty",
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
							3: {
								Title:  "Episode 3",
								Number: 3,
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
							3: {
								Title:  "Episode 3",
								Number: 3,
							},
						},
					},
				},
			},
		},
		{
			name: "Series with Specials",
			input: &tvShow{
				SeriesName: "Dexter",
				Series: &tvSeriesEpisodes{
					Episodes: []*episode{
						{
							EpisodeName:        "Bloopers",
							AiredEpisodeNumber: 1,
							AiredSeason:        0,
						},
						{
							EpisodeName:        "Cast Interviews",
							AiredEpisodeNumber: 2,
							AiredSeason:        0,
						},
						{
							EpisodeName:        "Episode 1",
							AiredEpisodeNumber: 1,
							AiredSeason:        1,
						},
						{
							EpisodeName:        "Episode 1",
							AiredEpisodeNumber: 1,
							AiredSeason:        2,
						},
					},
				},
			},
			expected: &types.TV{
				Title:       "Dexter",
				SeriesCount: 2, //This should ignore series 0 as special episodes extras etc.
				Series: map[int]*types.Series{
					0: {
						Title:  "Specials",
						Number: 0,
						Episodes: map[int]*types.Episode{
							1: {
								Title:  "Bloopers",
								Number: 1,
							},
							2: {
								Title:  "Cast Interviews",
								Number: 2,
							},
						},
					},
					1: {
						Title:  "Season 1",
						Number: 1,
						Episodes: map[int]*types.Episode{
							1: {
								Title:  "Episode 1",
								Number: 1,
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
						},
					},
				},
			},
		},
	}

	//run tests
	for _, test := range tests {
		result := buildTV(test.input)
		if diff := pretty.Compare(test.expected, result); diff != "" {
			t.Errorf("%s unexpected diff (-want +got):\n%s", test.name, diff)
		}
	}
}
