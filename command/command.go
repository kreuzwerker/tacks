package command

type Command interface {
	Run() error
}

func Background(c Command, f func(...interface{})) {
	go Foreground(c, f)
}

func Foreground(c Command, f func(...interface{})) {

	if err := c.Run(); err != nil {
		f(err)
	}

}
