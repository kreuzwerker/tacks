package tacks

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	bad     = `false`
	missing = `no_such_binary_exists`
	good1   = `ls -1 | sort -r | head -n 1`
	good2   = `ls -1 | sort | head -n 1`
)

func TestCommand(t *testing.T) {

	t.Parallel()

	assert := assert.New(t)

	c := Command(good1)
	res, err := c.Value()

	assert.Equal("vars_test.go", res)
	assert.NoError(err)

	c = Command(bad)
	res, err = c.Value()

	assert.Empty(res)
	assert.Equal(`failed to execute "false": "" (exit status 1)`, err.Error())

	c = Command(missing)
	res, err = c.Value()

	assert.Empty(res)
	assert.Equal(`failed to execute "no_such_binary_exists": /bin/bash: no_such_binary_exists: command not found (exit status 127)`, err.Error())

}

func TestCommandEnv(t *testing.T) {

	t.Parallel()

	assert := assert.New(t)

	c := Command(`env`)
	res, err := c.Value()
	assert.NoError(err)

	var (
		golang = os.Environ()
		shell  = strings.Split(res.(string), "\n")
	)

	for _, e := range golang {

		if strings.HasPrefix(e, "PATH=") {

			for _, ee := range shell {

				if strings.HasPrefix(ee, "PATH=") {

					assert.Equal(e, ee)
					return

				}

			}

		}

	}

	assert.Fail("PATH not inherited")

}

func TestCommandPipefail(t *testing.T) {

	assert := assert.New(t)

	c := Command(`false | echo`)
	_, err := c.Value()

	assert.Equal(`failed to execute "false | echo": "" (exit status 1)`, err.Error())

}
