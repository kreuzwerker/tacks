package tacks

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"regexp"
)

type Execution struct {
	Command  string
	External bool
	File     io.Reader `json:"-"`
	Filename string
	PipeOpen bool
}

func NewExecution() (e *Execution, err error) {

	e = new(Execution)
	e.PipeOpen = isPipeOpen()

	if filename, external := isExternal(); external {

		e.Command, e.File, err = readEmbedded(filename)
		e.External = true
		e.Filename = filename

		if err != nil {
			return nil, err
		}

	}

	return e, nil

}

func isExternal() (string, bool) {

	arg0 := path.Base(os.Args[0])
	external := os.Getenv("_")

	return external, arg0 != path.Base(external)

}

func isPipeOpen() bool {

	stat, err := os.Stdin.Stat()

	return err == nil && (stat.Mode()&os.ModeCharDevice) == 0

}

func readEmbedded(file string) (string, io.Reader, error) {

	fd, err := os.Open(file)

	if err != nil {
		return Nothing, nil, err
	}

	defer fd.Close()

	var (
		buf  = new(bytes.Buffer)
		cmd  string
		line = 1
	)

	scanner := bufio.NewScanner(fd)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {

		text := scanner.Text()

		if line == 1 {

			if matches := regexp.MustCompile(`^#!\s*(.+)$`).FindAllStringSubmatch(text, 1); len(matches) < 1 {
				return Nothing, nil, fmt.Errorf("No hashbang found in line 1: %q", text)
			} else {
				cmd = matches[0][1]
			}

		} else if _, err := buf.WriteString(fmt.Sprintf("%s\n", text)); err != nil {
			return Nothing, nil, err
		}

		line = line + 1

	}

	return cmd, buf, nil

}
