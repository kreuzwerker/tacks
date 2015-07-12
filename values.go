package tacks

import "sync"

// Maybe encapsulates a successful or failed call to the Value
// method a given Value interface
type Maybe struct {
	Error error
	Value interface{}
}

// NewError represents a failed call to a Value method
func NewError(err error) Maybe {
	return Maybe{Error: err}
}

// NewValue represents a successful call to a Value method
func NewValue(value interface{}) Maybe {
	return Maybe{Value: value}
}

// Value implementers are lazy-evaluated configuration sources
// such as constants, environment values and bash executions
type Value interface {
	Value() (interface{}, error)
}

// Values are the sum of all values defined within the environment
// of a Template
type Values map[string]Value

// Evaluate concurrently call Value on the defined values and
// returns appropriate Maybe instances
func (v Values) Evaluate() map[string]Maybe {

	result := make(map[string]Maybe, len(v))

	wg := &sync.WaitGroup{}
	wg.Add(len(v))

	lock := &sync.Mutex{}

	for k, v := range v {

		go func(k string, v Value) {

			resp, err := v.Value()
			defer wg.Done()

			lock.Lock()
			defer lock.Unlock()

			if err != nil {
				result[k] = NewError(err)
			} else {
				result[k] = NewValue(resp)
			}

		}(k, v)

	}

	wg.Wait()

	return result

}
