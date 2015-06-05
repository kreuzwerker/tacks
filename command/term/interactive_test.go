package term

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var buf = new(bytes.Buffer)

func answer(answer string) bool {

	buf.WriteString(answer)
	buf.WriteString("\n")
	defer buf.Reset()

	return Confirm("Are you sure? ")

}

func TestAsk(t *testing.T) {

	assert := assert.New(t)

	confirmIn = buf
	confirmOut = ioutil.Discard

	var tt = []struct {
		in  string
		out bool
	}{
		{"YES", true},
		{"yes", true},
		{"ye", false},
		{"y", false},
		{"no", false},
	}

	for _, t := range tt {
		assert.Equal(t.out, answer(t.in))
	}

}

func TestWait(t *testing.T) {

	assert := assert.New(t)

	now := time.Now()

	go func() {
		time.Sleep(100 * time.Millisecond)
		os.Stdin.WriteString("\n")
	}()

	Wait()

	elapsed := time.Now().Sub(now).Seconds()

	assert.InDelta(0.1, elapsed, 0.1)

}
