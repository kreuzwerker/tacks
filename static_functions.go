package tacks

import (
	"fmt"
	"strings"
	"text/template"
	"time"
)

// TODO: find more functions
// TODO: move functions into document

/*

  - SHA2 hash over stack
  - hostname of computer
  - ...

*/

var StaticFunctions = template.FuncMap{
	"now": func() string {
		return time.Now().UTC().Format("20060102-1504 UTC")
	},
	"partial": func(args ...string) (string, error) {
		res := fmt.Sprintf("Including %v", strings.Join(args, " + "))
		return res, nil
	},
}
