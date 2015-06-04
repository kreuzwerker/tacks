package tacks

import (
	"encoding/json"
	"fmt"
	"sort"
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

func (s Stack) IsIamCapabilitiesRequired() bool {

	capabilityIamRequired := []string{
		"AWS::CloudFormation::Stack",
		"AWS::IAM::AccessKey",
		"AWS::IAM::Group",
		"AWS::IAM::InstanceProfile",
		"AWS::IAM::Policy",
		"AWS::IAM::Role",
		"AWS::IAM::User",
		"AWS::IAM::UserToGroupAddition",
	}

	var types = s.Types()

	for _, req := range capabilityIamRequired {

		for _, got := range types {

			if req == got {
				return true
			}

		}

	}

	return false

}

func (s Stack) Types() []string {

	const _type = "Type"

	var (
		nothing = struct{}{}
		types   = make(map[string]struct{})
	)

	for _, resource := range s["Resources"].(Stack) {

		for key, value := range resource.(Stack) {

			if key == _type {
				types[value.(string)] = nothing
			}

		}

	}

	var sortedTypes sort.StringSlice

	for key, _ := range types {
		sortedTypes = append(sortedTypes, key)
	}

	sortedTypes.Sort()

	return sortedTypes

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
