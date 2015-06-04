package tacks

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func run(cmd string) (string, error) {

	const null = ""

	if out, err := exec.Command("/bin/bash", "-c", cmd).CombinedOutput(); err != nil {
		return null, err
	} else {
		return string(out), nil
	}

}

func init() {

	if out, err := exec.Command("/usr/bin/env", "go", "build", "-o", "/tmp/context", "test/context.go").CombinedOutput(); err != nil {
		panic(string(out))
	}

}

func TestContextHashbang(t *testing.T) {

	t.Parallel()

	assert := assert.New(t)

	out, err := run("test/context.sh")
	assert.NoError(err)
	assert.Equal(`A
B
C

`, out)

}

func TestContextHashbangFail(t *testing.T) {

	t.Parallel()

	assert := assert.New(t)

	_, err := run("test/context-err.sh")

	assert.Error(err)

}

func TestContextStandalone(t *testing.T) {

	t.Parallel()

	assert := assert.New(t)

	out, err := run("/tmp/context test/context.sh")

	assert.NoError(err)
	assert.Equal(`#! /TMP/CONTEXT
A
B
C

`, out)

}
