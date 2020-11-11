package filing

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

type FileStruct map[string][]*File

//forgets a dir (preventing any changes happening to it)
func (fs FileStruct) ForgetDir(dir string) {

	//forget dir and all subdirectories
	for d, _ := range fs {
		if strings.HasPrefix(d, dir) {
			delete(fs, dir)
		}
	}
}

//forgets a file (to prevent it being changed)
func (fs FileStruct) ForgetFile(dir string, fileName string) {
	if _, ok := fs[dir]; !ok {
		return
	}

	for i, file := range fs[dir] {
		if file.Name == fileName {
			//if last in array, use all elements up until point
			if i == len(fs[dir])-1 {
				fs[dir] = fs[dir][:i]
				continue
			}
			fs[dir] = append(fs[dir][:i], fs[dir][i+1:]...)
		}
	}

	if len(fs[dir]) == 0 {
		fs.ForgetDir(dir)
	}
}

type Filer struct {
	Root  string
	Files FileStruct //key: file path (string), value: name, extension (File)
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
		Root: root,
	}

	if err := filer.findFiles(); err != nil {
		return nil, err
	}

	return filer, nil
}

func (f *Filer) GetFiles() FileStruct {
	return f.Files
}

//returns map, with folder location as key, relevant contained files as values (array)
func (f *Filer) findFiles() error {

	files, err := f.listAllFiles()
	if err != nil {
		return err
	}

	mFiles := f.extractMediaFiles(files)

	if len(mFiles) == 0 {
		return fmt.Errorf(fmt.Sprintf("no supported media files located under specified root: %s", f.Root))
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

	f.Files = fileMap
	return nil
}

//recursively retrieve all files/directories under specified root
func (f *Filer) listAllFiles() ([]string, error) {

	var files []string

	if err := filepath.Walk(f.Root, func(location string, info os.FileInfo, err error) error {

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

	for loc, dir := range f.Files {
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
	for loc, dir := range f.Files {
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

			file.PrintDiff()
			fmt.Println()
		}
	}
}
