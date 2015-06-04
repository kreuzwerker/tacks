package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/kreuzwerker/tacks"
)

func main() {

	if e, err := tacks.NewExecution(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else if out, err := json.Marshal(e); err != nil {
		fmt.Println(err)
		os.Exit(2)
	} else {
		fmt.Println(string(out))
	}

}
