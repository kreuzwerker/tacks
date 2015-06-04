package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/kreuzwerker/tacks"
)

func main() {

	var filename string

	if len(os.Args) > 1 {
		filename = os.Args[1]
	}

	c := &tacks.Context{
		Args0:    os.Args[0],
		Filename: filename,
	}

	if err := c.DetectHashbang(); err != nil {
		log.Fatal(err)
	} else if data, err := c.Data(); err != nil {
		log.Fatal(err)
	} else if out, err := ioutil.ReadAll(data); err != nil {
		log.Fatal(err)
	} else {
		out := string(out)
		fmt.Println(strings.ToUpper(out))
	}

}
