package types

import "time"

type Movie struct {
	Title       string
	ReleaseDate time.Time
}

type TV struct {
	Title       string
	SeriesCount int
	ReleaseDate time.Time
	Series      map[int]*Series
}

type Series struct {
	Title    string
	Number   int
	Episodes map[int]*Episode
}

type Episode struct {
	Title  string
	Number int
}

//Constructors
func NewTV() *TV {
	return &TV{
		Series: make(map[int]*Series),
	}
}

func NewSeries() *Series{
	return &Series{
		Episodes: make(map[int]*Episode),
	}
}

func NewEpisode() *Episode{
	return &Episode{}
}
