package tacks

import "sync"

type Maybe struct {
	Error error
	Value interface{}
}

func NewError(err error) Maybe {
	return Maybe{Error: err}
}

func NewValue(value interface{}) Maybe {
	return Maybe{Value: value}
}

type Value interface {
	Value() (interface{}, error)
}

type Values map[string]Value

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
