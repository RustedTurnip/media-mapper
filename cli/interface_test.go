package cli

import (
	"testing"

	"github.com/kylelemons/godebug/pretty"
	"github.com/rustedturnip/media-mapper/filing"
)

func TestClInterface_getOrderedDirs(t *testing.T) {

	var tests = []struct {
		name       string
		files      filing.FileStruct
		transforms func(cli *clInterface)
		expected   []string
	}{
		{
			name: "Test Ordered Dirs - No Transforms",
			files: filing.FileStruct{
				"abc/123": {},
				"abc":     {},
				"xyz":     {},
				"abc/234": {},
			},
			transforms: func(cli *clInterface) {
				//no tansforms
			},
			expected: []string{
				"abc",
				"abc/123",
				"abc/234",
				"xyz",
			},
		},
		{
			name: "Test Ordered Dirs - Forget High Level Dir",
			files: filing.FileStruct{
				"abc/123": {},
				"abc":     {},
				"xyz":     {},
				"abc/234": {},
			},
			transforms: func(cli *clInterface) {
				cli.interactiveForget("abc") //forget abc (and all subdirs)
			},
			expected: []string{
				"xyz",
			},
		},
		{
			name: "Test Ordered Dirs - Forget All",
			files: filing.FileStruct{
				"abc/123": {},
				"abc":     {},
				"xyz":     {},
				"abc/234": {},
			},
			transforms: func(cli *clInterface) {
				cli.interactiveForget("abc/123")
				cli.interactiveForget("abc")
				cli.interactiveForget("xyz")
				cli.interactiveForget("abc/234")
			},
			expected: []string{},
		},
	}

	for _, test := range tests {
		cli := &clInterface{
			actionedDirs: map[string]struct{}{},
			filer: &filing.Filer{
				Files: test.files,
			},
		}

		test.transforms(cli) //run tests transforms

		result := cli.getOrderedDirs()

		if diff := pretty.Compare(test.expected, result); diff != "" {
			t.Errorf("%s unexpected diff (-want +got):\n%s", test.name, diff)
		}
	}
}

func TestClInterface_interactiveSkipDir(t *testing.T) {

	var tests = []struct {
		name     string
		files    filing.FileStruct
		inputs   []string
		expected map[string]struct{}
	}{
		{
			name: "SkipDir - One Nested Dir",
			files: filing.FileStruct{
				"test":      {},
				"test/dir":  {},
				"other_dir": {},
			},
			inputs: []string{
				"test/dir",
			},
			expected: map[string]struct{}{
				"test/dir": {},
			},
		},
		{
			name: "SkipDir - One Outer Dir",
			files: filing.FileStruct{
				"test":      {},
				"test/dir":  {},
				"other_dir": {},
			},
			inputs: []string{
				"test",
			},
			expected: map[string]struct{}{
				"test":     {},
				"test/dir": {},
			},
		},
	}

	for _, test := range tests {
		cli := &clInterface{
			actionedDirs: map[string]struct{}{},
			filer: &filing.Filer{
				Files: test.files,
			},
		}

		for _, dir := range test.inputs {
			cli.interactiveSkipDir(dir)
		}

		if diff := pretty.Compare(test.expected, cli.actionedDirs); diff != "" {
			t.Errorf("%s unexpected diff (-want +got):\n%s", test.name, diff)
		}

	}
}

func TestClInterface_interactiveStepIn(t *testing.T) {

	var tests = []struct {
		name     string
		files    filing.FileStruct
		inputs   map[string][]string //map[dir]inputStrings
		expected filing.FileStruct
	}{
		{
			name: "Normal Step In - Mixed Inputs, 1 Directory",
			files: filing.FileStruct{
				"movies": {
					{
						Name:    "jaws.1",
						NewName: "Jaws",
					},
					{
						Name:    "jaws.2",
						NewName: "Jaws 2",
					},
					{
						Name:    "batman_the_dark_knight",
						NewName: "Batman: The Dark Knight",
					},
					{
						Name:    "batman.begins",
						NewName: "Batman Begins",
					},
				},
			},
			inputs: map[string][]string{
				"movies": {"y", "n", "n", "y"},
			},
			expected: filing.FileStruct{
				"movies": {
					{
						Name:    "jaws.1",
						NewName: "Jaws",
					},
					{
						Name:    "batman.begins",
						NewName: "Batman Begins",
					},
				},
			},
		},
		{
			name: "Normal Step In - Mixed Inputs (Non-No), Multi-Directories",
			files: filing.FileStruct{
				"movies": {
					{
						Name:    "jaws.1",
						NewName: "Jaws",
					},
					{
						Name:    "jaws.2",
						NewName: "Jaws 2",
					},
				},
				"tv": {
					{
						Name:    "taboo.s01eo1.episode.1",
						NewName: "Taboo 1x01 Episode 1",
					},
					{
						Name:    "taboo.s01eo2.episode.2",
						NewName: "Taboo 1x02 Episode 2",
					},
				},
			},
			inputs: map[string][]string{
				"movies": {"xvc", "qq"},
				"tv":     {"n", "y"},
			},
			expected: filing.FileStruct{
				"tv": {
					{
						Name:    "taboo.s01eo2.episode.2",
						NewName: "Taboo 1x02 Episode 2",
					},
				},
			},
		},
	}

	for _, test := range tests {
		cli := &clInterface{
			actionedDirs: map[string]struct{}{},
			filer: &filing.Filer{
				Files: test.files,
			},
		}

		for dir, _ := range test.files {
			cli.readStringInput = func() string { //consumes inputs, passing next in queue on each call
				next := test.inputs[dir][0]
				test.inputs[dir] = test.inputs[dir][1:len(test.inputs[dir])]
				return next
			}

			cli.interactiveStepIn(dir)
		}

		if diff := pretty.Compare(test.expected, cli.filer.Files); diff != "" {
			t.Errorf("%s unexpected diff (-want +got):\n%s", test.name, diff)
		}

	}
}

func TestClInterface_runChange(t *testing.T) {

	type input struct {
		file     filing.File
		inputStr string
	}

	var tests = []struct {
		name     string
		files    filing.FileStruct
		inputs   map[string][]input //map[dir]input
		expected filing.FileStruct
	}{
		{
			name: "Forget File - Multiple Files & Dirs",
			files: filing.FileStruct{
				"personal/movies": {
					{
						Name:    "american-beauty",
						NewName: "American Beauty (1999)",
						Ext:     ".mkv",
					},
					{
						Name:    "gladiator",
						NewName: "Gladiator (2000)",
						Ext:     ".mp4",
					},
				},
				"personal/other-media": {
					{
						Name:    "simpsons.s01e01",
						NewName: "Simpsons - 1x01 - Episode 1",
						Ext:     ".mkv",
					},
					{
						Name:    "futurama.s04.e05",
						NewName: "Futurama 4x05 - Episode 5",
						Ext:     ".mp4",
					},
				},
			},
			inputs: map[string][]input{
				"personal/other-media": {
					{
						file: filing.File{
							Name:    "futurama.s04.e05",
							NewName: "Futurama 4x05 - Episode 5",
							Ext:     ".mp4",
						},
						inputStr: "n",
					},
				},
			},
			expected: filing.FileStruct{
				"personal/movies": {
					{
						Name:    "american-beauty",
						NewName: "American Beauty (1999)",
						Ext:     ".mkv",
					},
					{
						Name:    "gladiator",
						NewName: "Gladiator (2000)",
						Ext:     ".mp4",
					},
				},
				"personal/other-media": {
					{
						Name:    "simpsons.s01e01",
						NewName: "Simpsons - 1x01 - Episode 1",
						Ext:     ".mkv",
					},
				},
			},
		},
		{
			name: "Forget All Files In Dir - Single Dir",
			files: filing.FileStruct{
				"personal/movies": {
					{
						Name:    "american-beauty",
						NewName: "American Beauty (1999)",
						Ext:     ".mkv",
					},
					{
						Name:    "gladiator",
						NewName: "Gladiator (2000)",
						Ext:     ".mp4",
					},
				},
			},
			inputs: map[string][]input{
				"personal/movies": {
					{
						file: filing.File{
							Name:    "american-beauty",
							NewName: "American Beauty (1999)",
							Ext:     ".mkv",
						},
						inputStr: "n",
					},
					{
						file: filing.File{
							Name:    "gladiator",
							NewName: "Gladiator (2000)",
							Ext:     ".mp4",
						},
						inputStr: "n",
					},
				},
			},
			expected: filing.FileStruct{},
		},
		{
			name: "Forget File - Unknown Input", //should remove change as if selecting no (n)
			files: filing.FileStruct{
				"personal/movies": {
					{
						Name:    "american-beauty",
						NewName: "American Beauty (1999)",
						Ext:     ".mkv",
					},
				},
			},
			inputs: map[string][]input{
				"personal/movies": {
					{
						file: filing.File{
							Name:    "american-beauty",
							NewName: "American Beauty (1999)",
							Ext:     ".mkv",
						},
						inputStr: "xyz",
					},
				},
			},
			expected: filing.FileStruct{},
		},
	}

	for _, test := range tests {

		cli := &clInterface{
			filer: &filing.Filer{
				Files: test.files,
			},
		}

		for dir, fileInput := range test.inputs {
			//set input handler
			for _, fi := range fileInput {
				cli.readStringInput = func() string {
					return fi.inputStr
				}
				cli.runChange(dir, fi.file)
			}
		}

		if diff := pretty.Compare(test.expected, cli.filer.Files); diff != "" {
			t.Errorf("%s unexpected diff (-want +got):\n%s", test.name, diff)
		}
	}
}
