package tacks

import (
	"fmt"
	"os"
)

type Env string

func (e Env) Value() (interface{}, error) {

	const nothing = ""

	if value := os.Getenv(string(e)); value == nothing {
		return nil, fmt.Errorf("no such environment key %q", string(e))
	} else {
		return value, nil
	}

}
