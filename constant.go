package tacks

// Constant represents a fixed value
type Constant string

// Value returns the contant as value
func (c Constant) Value() (interface{}, error) {
	return c, nil
}
