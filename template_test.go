package tacks

import (
	"io/ioutil"
	"os"
	"os/user"
	"testing"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestConfigEvaluate(t *testing.T) {

	logger.Level = logrus.DebugLevel

	assert := assert.New(t)

	file, err := os.Open("test/stack.yml")
	defer file.Close()

	assert.NoError(err)

	conf, err := NewConfigFromReader(file)
	assert.NoError(err)

	assert.Empty(conf.values)
	assert.NotEmpty(conf.Stack)

	user, _ := user.Current()

	for _, name := range []string{"production", "staging"} {

		env, err := conf.Evaluate(name)
		assert.NoError(err, "failure in env %q")

		assert.Equal(name == "production", env.Ask)

		if name == "staging" {
			assert.Equal("upsert", env.Mode)
		} else {
			assert.Equal("create", env.Mode)
		}

		assert.Equal(map[string]string{
			"foo": "bar",
		}, env.Tags)

		assert.EqualValues(15, env.Timeout)

		assert.EqualValues("Makefile", conf.values["a"])

		if name == "staging" {
			assert.EqualValues(user.Username, conf.values["b"])
		} else {
			assert.EqualValues("dennis", conf.values["b"])
		}

		assert.Regexp("tacks$", conf.values["c"])
		assert.EqualValues(time.Now().Unix(), conf.values["d"])

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
