package filing

import "fmt"

var (
	//supported supportedVideo file types
	supportedVideo = map[string]struct{}{
		".asf": {},
		".avi": {},
		".mov": {},
		".mp4": {},
		".ts":  {},
		".mkv": {},
		".wmv": {},
	}

	//supported supportedSubtitle file types
	supportedSubtitle = map[string]struct{}{
		".srt": {},
		".smi": {},
		".ssa": {},
		".ass": {},
		".vtt": {},
	}
)

type File struct {
	Name    string //file name without extension
	NewName string //
	Ext     string //file extension
}

func (f *File) GetName() string {
	return fmt.Sprintf("%s%s", f.Name, f.Ext)
}

func (f *File) GetNewName() string {
	return fmt.Sprintf("%s%s", f.NewName, f.Ext)
}
