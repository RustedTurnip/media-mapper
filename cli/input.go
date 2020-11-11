package cli

import (
	"bufio"
	"os"
	"strings"
)

type input string

const (
	YES         input = "y"
	NO          input = "n"
	ALL         input = "a"
	QUIT        input = "q"
	INTERACTIVE input = "i"
	FORGET      input = "f"
	OPTIONS     input = "?"
)

func readStdInput() string {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')

	return strings.Trim(text, "\r\n")
}
