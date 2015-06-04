package tacks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConstant(t *testing.T) {

	assert := assert.New(t)

	resp, err := Constant("Hello").Value()
	assert.NoError(err)
	assert.EqualValues("Hello", resp)

}
