package filing

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Filer struct {
	root  string
	files map[string][]*File //key: file path (string), value: name, extension (File)
}

func New(root string) (*Filer, error) {

	if root == "" {
		return nil, fmt.Errorf("no directory specified to work on")
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
			old := path.Join(loc, file.GetName())
			new := path.Join(loc, file.GetNewName())

			err := os.Rename(old, new)

			if err != nil {
				log.Println(fmt.Sprintf("Failed to rename file: %s", old))
			}
		}
	}
}
