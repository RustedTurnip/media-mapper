package builder

import (
	"github.com/rustedturnip/media-mapper/types"
)


//Builder definitions
type TVBuilder struct {
	functions []func(tv *types.TV) error
}

type SeriesBuilder struct {
	functions []func(series *types.Series) error
}

type EpisodeBuilder struct {
	functions []func(series *types.Episode) error
}

//Constructors
func NewTVBuilder() *TVBuilder{
	return &TVBuilder{
		functions: make([]func(tv *types.TV) error, 0),
	}
}

func NewSeriesBuilder() *SeriesBuilder{
	return &SeriesBuilder{
		functions: make([]func(tv *types.Series) error, 0),
	}
}

func NewEpisodeBuilder() *EpisodeBuilder{
	return &EpisodeBuilder{
		functions: make([]func(tv *types.Episode) error, 0),
	}
}


//Build functions
func (tvb *TVBuilder) Build () *types.TV {

	tv := types.NewTV()

	for _, f := range tvb.functions {
		if err := f(tv); err != nil {
			//todo
			panic(err)
		}
	}

	return tv
}

func (sb *SeriesBuilder) Build () *types.Series {

	series := types.NewSeries()

	for _, f := range sb.functions {
		if err := f(series); err != nil {
			//todo
			panic(err)
		}
	}

	return series
}

func (eb *EpisodeBuilder) Build () *types.Episode {

	episode := types.NewEpisode()

	for _, f := range eb.functions {
		if err := f(episode); err != nil {
			//todo
			panic(err)
		}
	}

	return episode
}


//TV Builder functions
func (tvb *TVBuilder) WithTitle(title string) *TVBuilder {
	tvb.functions = append(tvb.functions, func(tv *types.TV) error {
		tv.Title = title
		return nil
	})

	return tvb
}

func (tvb *TVBuilder) WithSeriesCount(count int) *TVBuilder {
	tvb.functions = append(tvb.functions, func(tv *types.TV) error {
		tv.SeriesCount = count
		return nil
	})

	return tvb
}

func (tvb *TVBuilder) WithSeries(seriesBuilder *SeriesBuilder) *TVBuilder {
	tvb.functions = append(tvb.functions, func(tv *types.TV) error {

		series := seriesBuilder.Build()
		tv.Series[series.Number] = series
		return nil
	})

	return tvb
}


//Series Builder functions
func (sb *SeriesBuilder) WithTitle(title string) *SeriesBuilder {
	sb.functions = append(sb.functions, func(s *types.Series) error {
		s.Title = title
		return nil
	})

	return sb
}

func (sb *SeriesBuilder) WithNumber(number int) *SeriesBuilder {
	sb.functions = append(sb.functions, func(s *types.Series) error {
		s.Number = number
		return nil
	})

	return sb
}

func (sb *SeriesBuilder) WithEpisode(episodeBuilder *EpisodeBuilder) *SeriesBuilder {
	sb.functions = append(sb.functions, func(s *types.Series) error {

		episode := episodeBuilder.Build()
		s.Episodes[episode.Number] = episode
		return nil
	})

	return sb
}


//Episode Builder functions
func (eb *EpisodeBuilder) WithTitle(title string) *EpisodeBuilder {
	eb.functions = append(eb.functions, func(e *types.Episode) error{
		e.Title = title
		return nil
	})

	return eb
}

func (eb *EpisodeBuilder) WithNumber(number int) *EpisodeBuilder {
	eb.functions = append(eb.functions, func(e *types.Episode) error{
		e.Number = number
		return nil
	})

	return eb
}
