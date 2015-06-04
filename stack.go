package tacks

import (
	"encoding/json"
	"fmt"
)

type Stack map[interface{}]interface{}

func (s Stack) MarshalText() ([]byte, error) {

	if cast, err := jsonify(s); err != nil {
		return nil, err
	} else if cast, ok := cast.(map[string]interface{}); !ok {
		return nil, fmt.Errorf("unexpected value %v (%T)", cast, cast)
	} else if text, err := json.MarshalIndent(cast, "", "  "); err != nil {
		return nil, err
	} else {
		return text, nil
	}

}

func jsonify(in interface{}) (interface{}, error) {

	switch cast := in.(type) {
	case Stack:
		return jsonify(map[interface{}]interface{}(cast))

	case map[interface{}]interface{}:

		var out = make(map[string]interface{}, len(cast))

		for key, value := range cast {

			if value, err := jsonify(value); err != nil {
				return nil, err
			} else {
				out[key.(string)] = value
			}

		}

		return out, nil

	case []interface{}:

		var out = make([]interface{}, len(cast))

		for idx, value := range cast {

			if value, err := jsonify(value); err != nil {
				return nil, err
			} else {
				out[idx] = value
			}

		}

		return out, nil

	case float64, float32, int64, int, string:
		return cast, nil

	default:
		return nil, fmt.Errorf("unhandled value %v (%T)", cast, cast)

	}

}
