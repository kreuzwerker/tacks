package tacks

import (
	"fmt"
	"os"
)

// Env represents an environment variable key
type Env string

// Value returns the value of the environment key defined
func (e Env) Value() (interface{}, error) {

	const nothing = ""

	if value := os.Getenv(string(e)); value == nothing {
		return nil, fmt.Errorf("no such environment key %q", string(e))
	} else {
		return value, nil
	}

}
