package tacks

import (
	"text/template"
	"time"
)

// TODO: find more functions

/*

  - SHA2 hash over stack
  - hostname of computer
  - ...

*/

var Functions = template.FuncMap{
	"now": func() string {
		return time.Now().UTC().Format("20060102-1504 UTC")
	},
}
