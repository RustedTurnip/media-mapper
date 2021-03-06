package controller

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	colour "github.com/fatih/color"
	ptn "github.com/middelink/go-parse-torrent-name"
	"github.com/rustedturnip/media-mapper/dbs"
	"github.com/rustedturnip/media-mapper/filing"
)

const (
	parseErr = "%s: failed to parse - %s"

	tvTitleFmt    = "%s - %dx%d - %s"
	movieTitleFmt = "%s (%d)"
)

type Worker struct {
	database   dbs.Database
	filer      *filing.Filer
	streamline bool
	errs       []error
}

func New(database dbs.Database, filer *filing.Filer, streamline bool) *Worker {
	return &Worker{
		database:   database,
		filer:      filer,
		streamline: streamline,
	}
}

func (w *Worker) Do() {

	for _, files := range w.filer.GetFiles() {
		for _, file := range files {
			info, err := ptn.Parse(file.Name)
			if err != nil {
				fName := fmt.Sprintf("%s.%s", file.Name, file.Ext)
				w.errs = append(w.errs, fmt.Errorf(parseErr, fName, file))
			}

			file.NewName = w.getName(info)
		}
	}

	//print diff
	if !w.streamline {
		w.filer.PrintBatchDiff()
	}

	//display failed files
	if len(w.errs) > 1 {
		fmt.Println("\nMatch errors:")
		for _, err := range w.errs {
			colour.Yellow("! %s", err.Error())
		}
	}

	//user input, proceed?
	if !w.streamline {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Proceed with changes? (y/n): ")
		text, _ := reader.ReadString('\n')

		if strings.ToLower(strings.Trim(text, "\r\n")) != "y" {
			fmt.Println("Cancelling...")
			return
		}
	}

	//continue with file rename
	if !w.streamline {
		fmt.Println("Renaming files...")
	}
	w.filer.RenameBatch()
}

func (w *Worker) getName(info *ptn.TorrentInfo) string {

	switch info.Episode {
	case 0: //Movie
		//TODO - better movie selection
		results := w.database.SearchMovies(info.Title)
		if len(results) == 0 {
			return ""
		}
		movie := results[0]
		return fmt.Sprintf(movieTitleFmt, movie.Title, movie.ReleaseDate.Year())

	default: //Episode of TV Series
		//TODO - better movie selection
		results := w.database.SearchTV(info.Title)
		if len(results) == 0 {
			return ""
		}
		show := results[0]

		if _, ok := show.Series[info.Season]; !ok {
			return "" //can't find series
		}
		series := show.Series[info.Season]

		if _, ok := series.Episodes[info.Episode]; !ok {
			return "" //can't find episode in series
		}
		episode := series.Episodes[info.Episode]

		return fmt.Sprintf(tvTitleFmt, show.Title, series.Number, episode.Number, episode.Title)
	}
}
