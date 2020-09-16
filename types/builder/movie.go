package builder

import (
	"log"
	"time"

	"github.com/rustedturnip/media-mapper/types"
)

type MovieBuilder struct {
	functions []func(movie *types.Movie) error
}

func NewMovieBuilder() *MovieBuilder {
	return &MovieBuilder{
		functions: make([]func(*types.Movie) error, 0),
	}
}

func (mb *MovieBuilder) Build () *types.Movie {

	movie := &types.Movie{}

	for _, f := range mb.functions {
		if err := f(movie); err != nil {
			log.Println(err.Error())
		}
	}

	return movie
}

func (mb *MovieBuilder) WithTitle(title string) *MovieBuilder {
	mb.functions = append(mb.functions, func(m *types.Movie) error {
		m.Title = title
		return nil
	})

	return mb
}

func (mb *MovieBuilder) WithReleaseDate(date time.Time) *MovieBuilder {
	mb.functions = append(mb.functions, func(m *types.Movie) error {
		m.ReleaseDate = date
		return nil
	})

	return mb
}
