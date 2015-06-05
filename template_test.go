package tacks

import (
	"io/ioutil"
	"os"
	"os/user"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTemplateEvaluate(t *testing.T) {

	assert := assert.New(t)

	file, err := os.Open("test/stack.yml")
	defer file.Close()

	assert.NoError(err)

	conf, err := NewTemplateFromReader(file)
	assert.NoError(err)

	user, _ := user.Current()

	for _, name := range []string{"production", "staging"} {

		err := conf.Evaluate(name, func(d Document) error {

			assert.Equal(name, d.Environment.Name)

			if name == "staging" {
				assert.Equal("upsert", d.Environment.Mode)
			} else {
				assert.Equal("create", d.Environment.Mode)
			}

			assert.Equal(map[string]string{
				"foo": "bar",
			}, d.Environment.Tags)

			assert.EqualValues(15, d.Environment.Timeout)

			assert.EqualValues("Makefile", d.Variables["a"])

			if name == "staging" {
				assert.EqualValues(user.Username, d.Variables["b"])
			} else {
				assert.EqualValues("dennis", d.Variables["b"])
			}

			vd, _ := strconv.Atoi(d.Variables["d"].(string))

			assert.Regexp("tacks$", d.Variables["c"])
			assert.InDelta(time.Now().Unix(), vd, 1.0)

			return nil

		})

		assert.NoError(err, "failure in env %q")

	}

	var (
		stats []os.FileInfo
		datas []string
	)

	for _, e := range []string{"/tmp/pre.log", "/tmp/post.log"} {

		defer os.Remove(e)

		stat, err := os.Stat(e)
		assert.NoError(err)

		stats = append(stats, stat)

		data, err := ioutil.ReadFile(e)
		assert.NoError(err)

		datas = append(datas, string(data))

	}

	assert.Equal("pre\ndone\n", datas[0])
	assert.Equal("post\ndone\n", datas[1])
	assert.Len(stats, 2)
	assert.True(stats[0].ModTime().Before(stats[1].ModTime()))

}
