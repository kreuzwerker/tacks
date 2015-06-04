package tacks

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestStackMarshalText(t *testing.T) {

	assert := assert.New(t)

	in := []byte(`{
  "a": "b",
  "c": [
    0,
    1,
    {
      "d": {
        "e": "f"
      }
    },
    3
  ]
}`)

	var stack Stack

	if err := yaml.Unmarshal(in, &stack); err != nil {
		assert.NoError(err)
	}

	out, err := stack.MarshalText()
	assert.NoError(err)

	assert.Equal(in, out)

}
