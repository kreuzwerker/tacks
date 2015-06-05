package term

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

var (
	confirmIn  io.Reader = os.Stdin
	confirmOut io.Writer = os.Stdout
)

func Confirm(question string, args ...interface{}) bool {

	const YES = "YES"

	var (
		err    error
		answer string
	)

	fmt.Fprintf(confirmOut, question, args...)

	if _, err = fmt.Fscanln(confirmIn, &answer); err != nil {
		return false
	} else {
		return strings.ToUpper(answer) == YES
	}

}

func Size() (width, height int, err error) {
	return terminal.GetSize(int(os.Stdout.Fd()))
}

func Wait() {

	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()

	os.Stdin.Read(make([]byte, 1))

}
