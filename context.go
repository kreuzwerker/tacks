package tacks

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"regexp"
)

// Context encapsulates access to some data via hashbang or some given filename
type Context struct {
	Args0    string
	Filename string
	Hashbang bool
	data     io.ReadCloser
}

// Data returns the data associated with this context: either the file initially set or
// the script itself (minus the first line) when the mode is hashbang
func (c *Context) Data() (io.ReadCloser, error) {

	if c.data != nil {
		return c.data, nil
	} else if data, err := os.Open(c.Filename); err != nil {
		return nil, err
	} else {
		return data, nil
	}

}

// DetectHashbang changes context fields if the execution mode is hashbang
func (c *Context) DetectHashbang() error {

	var (
		args0    = path.Base(os.Args[0])
		file     = os.Getenv("_")
		hashbang = regexp.MustCompile(`^#!\s*(.+)$`)
	)

	if args0 != path.Base(file) {

		c.Hashbang = true
		c.Filename = file

		fd, err := os.Open(file)

		if err != nil {
			return err
		}

		defer fd.Close()

		buf := new(bytes.Buffer)
		line := 1

		scanner := bufio.NewScanner(fd)
		scanner.Split(bufio.ScanLines)

		for scanner.Scan() {

			text := scanner.Text()

			if line == 1 {

				if matches := hashbang.FindAllStringSubmatch(text, 1); len(matches) < 1 {
					return fmt.Errorf("no hashbang found in line 1: %q", text)
				} else {
					c.Args0 = matches[0][1]
				}

			} else if _, err := buf.WriteString(fmt.Sprintf("%s\n", text)); err != nil {
				return err
			}

			line = line + 1

		}

		c.data = ioutil.NopCloser(buf)

	}

	return nil

}
