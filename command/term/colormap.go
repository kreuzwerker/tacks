package term

import (
	"regexp"

	"github.com/fatih/color"
)

var CloudFormation = NewColorMap(map[string]*color.Color{
	`_FAILED$`:      color.New(color.FgRed),
	`_IN_PROGRESS$`: color.New(color.FgBlue),
	`_COMPLETE$`:    color.New(color.FgGreen),
})

type ColorMap struct {
	c []func(a ...interface{}) string
	r []*regexp.Regexp
}

func NewColorMap(m map[string]*color.Color) *ColorMap {

	var cm ColorMap

	for k, v := range m {

		cm.c = append(cm.c, v.SprintFunc())
		cm.r = append(cm.r, regexp.MustCompile(k))
	}

	return &cm

}

func (c ColorMap) Colorize(text string) string {

	for i, e := range c.r {

		if e.MatchString(text) {
			return c.c[i](text)
		}

	}

	return text

}
