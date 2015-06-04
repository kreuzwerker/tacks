package tacks

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type job func() (interface{}, error)

func (j job) Value() (interface{}, error) {
	return j()
}

func TestValues(t *testing.T) {

	assert := assert.New(t)

	var j1 job = func() (interface{}, error) {
		time.Sleep(100 * time.Millisecond)
		return 1, nil
	}

	var j2 job = func() (interface{}, error) {
		time.Sleep(200 * time.Millisecond)
		return 2, nil
	}

	var j3 job = func() (interface{}, error) {
		time.Sleep(5 * time.Millisecond)
		return 3, fmt.Errorf("not enough sleep!")
	}

	values := Values{
		"j1": j1,
		"j2": j2,
	}

	result := values.Evaluate()

	assert.Len(result, 2)

	assert.Equal(1, result["j1"].Value)
	assert.Equal(nil, result["j1"].Error)

	assert.Equal(2, result["j2"].Value)
	assert.Equal(nil, result["j2"].Error)

	values["j3"] = j3

	result = values.Evaluate()

	assert.Len(result, 3)

	assert.Equal(1, result["j1"].Value)
	assert.Nil(result["j1"].Error)

	assert.Equal(2, result["j2"].Value)
	assert.Nil(result["j2"].Error)

	assert.Equal(nil, result["j3"].Value)
	assert.NotNil(result["j3"].Error)

}
