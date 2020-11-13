package controller

import (
	"fmt"
	colour "github.com/fatih/color"
	ptn "github.com/middelink/go-parse-torrent-name"
	"github.com/rustedturnip/media-mapper/cli"
	"github.com/rustedturnip/media-mapper/database"
	"github.com/rustedturnip/media-mapper/filing"
	"github.com/schollz/progressbar/v3"
)

const (
	parseErr = "%s: failed to parse - %s"

	tvTitleFmt    = "%s - %dx%d - %s"
	movieTitleFmt = "%s (%d)"
)

type Worker struct {
	database   database.Database
	filer      *filing.Filer
	streamline bool
	errs       []error
}

func New(database database.Database, filer *filing.Filer, streamline bool) *Worker {
	return &Worker{
		database:   database,
		filer:      filer,
		streamline: streamline,
	}
}

func (w *Worker) Do() {

	//progress bar config
	bar := progressbar.NewOptions(w.filer.GetFileCount(),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetWidth(50),
		progressbar.OptionSetDescription("Fetching data..."),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[blue]=[reset]",
			SaucerHead:    "[blue]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))

	for _, files := range w.filer.GetFiles() {
		for _, file := range files {
			info, err := ptn.Parse(file.Name)
			if err != nil {
				fName := fmt.Sprintf("%s.%s", file.Name, file.Ext)
				w.errs = append(w.errs, fmt.Errorf(parseErr, fName, file))
			}

			file.NewName = w.getName(info)
			bar.Add(1) //progress bar
		}
	}

	//print diff
	if !w.streamline {
		fmt.Println()
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
		ui := cli.New(w.filer)
		ui.Run() //run interface for user to take control
	} else { //else silently rename all
		w.filer.RenameBatch()
	}
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
