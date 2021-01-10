package cli

import (
	"fmt"
	"sort"
	"strings"

	"github.com/rustedturnip/media-mapper/filing"
)

func New(filer *filing.Filer) *clInterface {
	return &clInterface{
		filer:           filer,
		readStringInput: readStdInput, //separated to allow for testing
	}
}

//hosts the cli interface allowing users to interactively select changes
type clInterface struct {
	actionedDirs    map[string]struct{}
	filer           *filing.Filer
	readStringInput func() string
}

//runs command line interface for user (outer/highest level interface)
func (cli *clInterface) Run() {
	for cli.runOuter() {
		cli.filer.PrintBatchDiff() //print diff at the end of each complete run
	} //run while runOuter == true
}

//runs the out-most interface
func (cli *clInterface) runOuter() bool {
	fmt.Print("How would you like to proceed [a/q/i/?]? ")

	switch input(strings.ToLower(cli.readStringInput())) {

	case ALL: //a, rename all files
		fmt.Println("Proceeding with all changes...")
		cli.filer.RenameBatch()
		return false

	case QUIT: //q - exit program
		fmt.Println("Exiting without making changes...")
		return false

	case INTERACTIVE: //i - interactive, choose changes

		cli.actionedDirs = make(map[string]struct{})                                            //initialise/reset actioned dirs
		for ordered := cli.getOrderedDirs(); len(ordered) > 0; ordered = cli.getOrderedDirs() { //while not all dirs actioned, get next alphabetically
			cli.runInteractive(ordered[0])
		}
		return true

	case OPTIONS: //? - list options
		fmt.Println("a - continue with ALL changes listed\n" +
			"q - exit without making changes\n" +
			"i - interactive, choose which changes to keep")
		return true

	default: //unrecognised input, retry
		fmt.Println("type '?' for options")
		return true
	}
}

//offers user interactive interface
func (cli *clInterface) runInteractive(dir string) {

	fmt.Print(fmt.Sprintf("\n%s\nStep into directory [y/n/f/?]? ", dir))
	switch input(strings.ToLower(cli.readStringInput())) {

	case YES: //y - step into
		cli.interactiveStepIn(dir)
		return

	case NO: //n - skip dir and sub dirs - (leaves changes to be implemented)
		cli.interactiveSkipDir(dir)
		return

	case FORGET: //f - remove (forget) nested changes
		cli.interactiveForget(dir)
		return

	case OPTIONS: //? - list options
		fmt.Println("y - yes (step into directory to review contained changes)\n" +
			"n - no (skip directory, all changes contained will be made)\n" +
			"f - forget (discard all changes contained within directory)")

		return

	default: //unrecognised input, retry
		fmt.Println("type '?' for options")

		return
	}

}

//steps into dir and offers choice to user for each file contained
func (cli *clInterface) interactiveStepIn(dir string) {

	files := cli.filer.GetFiles() //reference to files

	//make copy of files
	var filesCopy []filing.File
	for _, file := range files[dir] {
		filesCopy = append(filesCopy, *file)
	}

	for _, file := range filesCopy {
		cli.runChange(dir, file)
	}

	cli.actionedDirs[dir] = struct{}{} //action dir
}

//keeps changes of all files (recursively) within specified directory
func (cli *clInterface) interactiveSkipDir(dir string) {

	for _, iDir := range cli.getOrderedDirs() {

		if strings.HasPrefix(iDir, dir) {
			cli.actionedDirs[iDir] = struct{}{}
		}
	}
}

//removes changes of all files (recursively) within specified directory
func (cli *clInterface) interactiveForget(dir string) {

	for _, iDir := range cli.getOrderedDirs() {

		if strings.HasPrefix(iDir, dir) {
			cli.filer.GetFiles().ForgetDir(iDir)
			cli.actionedDirs[iDir] = struct{}{}
		}
	}
}

//returns list of remaining dirs in alphabetical order
func (cli *clInterface) getOrderedDirs() []string {

	list := make([]string, 0, len(cli.filer.GetFiles()))
	for k := range cli.filer.GetFiles() {

		if _, ok := cli.actionedDirs[k]; !ok {
			list = append(list, k) //if directory hasn't been actioned, add to list
		}
	}

	sort.Strings(list)
	return list
}

//offers user interface for individual change
func (cli *clInterface) runChange(dir string, file filing.File) {

	file.PrintDiff()
	fmt.Print("\nKeep change [y/n]? ")

	switch input(strings.ToLower(cli.readStringInput())) {

	case YES: //y - keep change (i.e. do nothing)
		return

	case NO: //n - remove change (see default)
		fallthrough
	default:
		cli.filer.GetFiles().ForgetFile(dir, file.Name)
	}
}
