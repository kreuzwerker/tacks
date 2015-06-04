package tacks

type Constant string

func (c Constant) Value() (interface{}, error) {
	return c, nil
}
