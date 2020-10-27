package cache

import (
	"fmt"
	"testing"

	"github.com/rustedturnip/media-mapper/types"
	"github.com/stretchr/testify/assert"
)

const (
	maxCalls = 1
)

type test struct {
	name   string
	inputs map[string]int //map[input]frequency
}

//implements database interface with internal calls tallied for testing
type testDBClient struct {
	movieTally map[string]int
	tvTally    map[string]int
}

func newTestDB() *testDBClient {
	return &testDBClient{
		movieTally: map[string]int{},
		tvTally:    map[string]int{},
	}
}

func (tc *testDBClient) SearchMovies(title string) []*types.Movie {

	if _, ok := tc.movieTally[title]; !ok {
		tc.movieTally[title] = 0
	}

	tc.movieTally[title]++
	return nil //no response necessary, only need to count number of calls made to function for each title
}

func (tc *testDBClient) SearchTV(title string) []*types.TV {
	if _, ok := tc.tvTally[title]; !ok {
		tc.tvTally[title] = 0
	}

	tc.tvTally[title]++
	return nil //no response necessary, only need to count number of calls made to function for each title
}

func TestDatabaseCache_SearchTV(t *testing.T) {

	var tests = []test{
		{
			name: "Cached DB - SearchTV()",
			inputs: map[string]int{
				"broadchurch":   3,
				"fawlty towers": 2,
				"Happy Vally":   6,
			},
		},
	}

	for _, test := range tests {
		//create cached-db instance with testDBClient
		client := newTestDB()
		db := New(client)

		for title, count := range test.inputs {
			for i := 0; i < count; i++ {
				db.SearchTV(title)
			}
		}

		for title, count := range client.tvTally {
			assert.Equal(t, maxCalls, count, fmt.Sprintf("caching failed, SearchTV() was called %d times for: %s", count, title))
		}
	}
}

func TestDatabaseCache_SearchMovies(t *testing.T) {

	var tests = []test{
		{
			name: "Cached DB - SearchMovies()",
			inputs: map[string]int{
				"Requiem for a Dream": 8,
				"TPB AFK":             2,
				"gladiator":           1,
			},
		},
	}

	for _, test := range tests {
		//create cached-db instance with testDBClient
		client := newTestDB()
		db := New(client)

		for title, count := range test.inputs {
			for i := 0; i < count; i++ {
				db.SearchMovies(title)
			}
		}

		for title, count := range client.tvTally {
			assert.Equal(t, maxCalls, count, fmt.Sprintf("caching failed, SearchMovie() was called %d times for: %s", count, title))
		}
	}
}
