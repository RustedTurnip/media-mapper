package filing

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	colour "github.com/fatih/color"
)

type Filer struct {
	root  string
	files map[string][]*File //key: file path (string), value: name, extension (File)
}

func New(root string) (*Filer, error) {

	//try and use working directory if none specified ($ pwd)
	if root == "" {
		if wd, err := os.Getwd(); err != nil {
			return nil, fmt.Errorf("unable to establish a working directory")
		} else {
			root = wd
		}
	}

	filer := &Filer{
		root: root,
	}

	if err := filer.findFiles(); err != nil {
		return nil, err
	}

	return filer, nil
}

func (f *Filer) GetFiles() map[string][]*File {
	return f.files
}

//returns map, with folder location as key, relevant contained files as values (array)
func (f *Filer) findFiles() error {

	files, err := f.listAllFiles()
	if err != nil {
		return err
	}

	mFiles := f.extractMediaFiles(files)

	if len(mFiles) == 0 {
		return fmt.Errorf(fmt.Sprintf("no supported media files located under specified root: %s", f.root))
	}

	fileMap := make(map[string][]*File)

	for _, file := range mFiles {
		path := filepath.Dir(file)
		ext := filepath.Ext(file)
		name := strings.TrimSuffix(filepath.Base(file), ext)

		fileMap[path] = append(fileMap[path], &File{
			Name: name,
			Ext:  ext,
		})
	}

	f.files = fileMap
	return nil
}

//recursively retrieve all files/directories under specified root
func (f *Filer) listAllFiles() ([]string, error) {

	var files []string

	if err := filepath.Walk(f.root, func(location string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		files = append(files, location)
		return nil

	}); err != nil { //.Walk() error handle
		return nil, err
	}

	return files, nil
}

func (f *Filer) extractMediaFiles(files []string) []string {

	var mediaFiles []string

	for _, file := range files {
		extension := filepath.Ext(file)

		if _, ok := supportedVideo[strings.ToLower(extension)]; ok {
			mediaFiles = append(mediaFiles, file)
		}
	}

	//TODO same for supportedSubtitle files

	return mediaFiles
}

func (f *Filer) RenameBatch() {

	for loc, dir := range f.files {
		for _, file := range dir {

			if file.NewName == "" {
				continue
			}

			old := path.Join(loc, file.GetName())
			new := path.Join(loc, file.GetNewName())

			err := os.Rename(old, new)

			if err != nil {
				log.Println(fmt.Sprintf("Failed to rename file: %s", old))
			}
		}
	}
}

func (f *Filer) PrintBatchDiff() {

	//print location of diffs
	for loc, dir := range f.files {
		var once sync.Once //only want to print loc once per directory

		//print diff
		for _, file := range dir {
			if file.NewName == "" {
				continue
			}

			//only want to print loc once per directory
			once.Do(func() {
				fmt.Println(fmt.Sprintf("\ndiff %s:", loc))
			})

			colour.Red("- %s", file.GetName())
			colour.Green("+ %s", file.GetNewName())
			fmt.Println()
		}
	}
}
