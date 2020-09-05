package tvdb

import "github.com/rustedturnip/media-mapper/dbs"

var (
	login = "https://api.thetvdb.com/login"
)

type TVDB struct {
}

func New() dbs.Database {
	//TODO
	return nil
}

func (db *TVDB) GetMovie(title string) string {

	return title
}
