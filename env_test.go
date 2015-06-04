package tacks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnv(t *testing.T) {

	assert := assert.New(t)

	resp, err := Env("NO_SUCH_KEY_DEFINED").Value()
	assert.Error(err)
	assert.Nil(resp)

	resp, err = Env("PWD").Value()
	assert.NoError(err)
	assert.Regexp("tacks$", resp)

}
