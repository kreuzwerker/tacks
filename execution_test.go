package tacks

import (
	"encoding/json"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func run(cmd string) (*Execution, error) {

	var e Execution

	if out, err := exec.Command("/bin/bash", "-c", cmd).Output(); err != nil {
		return nil, err
	} else if err := json.Unmarshal(out, &e); err != nil {
		return nil, err
	} else {
		return &e, nil
	}

}

func init() {

	if err := exec.Command("/usr/bin/env", "go", "build", "-o", "/tmp/execution", "test/execution.go").Run(); err != nil {
		panic(err)
	}

}

func TestPipe(t *testing.T) {

	assert := assert.New(t)

	e, err := run("/bin/cat test/ugly.json | /tmp/execution")
	assert.NoError(err)

	assert.Equal(Nothing, e.Command)
	assert.Equal(false, e.External)
	assert.Equal(Nothing, e.Filename)
	assert.Equal(true, e.PipeOpen)

}

func TestShebang(t *testing.T) {

	assert := assert.New(t)

	e, err := run("test/execution.sh")
	assert.NoError(err)

	assert.Equal("/tmp/execution", e.Command)
	assert.Equal(true, e.External)
	assert.Equal("test/execution.sh", e.Filename)
	assert.Equal(false, e.PipeOpen)

}

func TestShebangFail(t *testing.T) {

	assert := assert.New(t)

	e, err := run("test/execution-err.sh")
	assert.Error(err)

	assert.Nil(e)

}

func TestStandalone(t *testing.T) {

	assert := assert.New(t)

	e, err := run("/tmp/execution")
	assert.NoError(err)

	assert.Equal(Nothing, e.Command)
	assert.Equal(false, e.External)
	assert.Equal(Nothing, e.Filename)
	assert.Equal(false, e.PipeOpen)

}
