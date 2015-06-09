package tacks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"text/template"
)

type Document struct {
	stack       map[interface{}]interface{}
	Environment *Environment
	Variables   map[string]interface{}
}

func (d Document) IsIamCapabilitiesRequired() bool {

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

	var types = d.Types()

	for _, req := range capabilityIamRequired {

		for _, got := range types {

			if req == got {
				return true
			}

		}

	}

	return false

}

func (d Document) Parse() (string, error) {

	const null = ""

	tree, err := jsonify(d.stack)

	if err != nil {
		return null, err
	}

	text, err := json.MarshalIndent(tree.(map[string]interface{}), "", "  ")

	if err != nil {
		return null, err
	}

	tpl := template.New("_root")
	tpl.Funcs(StaticFunctions) // TODO: reference runtime template functions

	// f["name"] = func() string { // TODO: remove this
	//
	// 	name := []string{
	// 		r.Environment,
	// 		r.Name,
	// 	}
	//
	// 	return strings.Join(name, "-")
	//
	// }

	tpl, err = tpl.Parse(string(text))

	if err != nil {
		return null, err
	}

	buf := new(bytes.Buffer)

	err = tpl.Execute(buf, d)

	return buf.String(), err

}

func (d Document) Types() []string {

	const _type = "Type"

	var (
		nothing = struct{}{}
		types   = make(map[string]struct{})
	)

	stack, err := d.json()

	if err != nil {
		panic(err)
	}

	resources := stack["Resources"].(map[string]interface{})

	for _, resource := range resources {

		for key, value := range resource.(map[string]interface{}) {

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

func (d Document) json() (map[string]interface{}, error) {

	if stack, err := jsonify(d.stack); err != nil {
		return nil, err
	} else if casted, ok := stack.(map[string]interface{}); !ok {
		return nil, fmt.Errorf("unexpected stack structure (%T)", stack)
	} else {
		return casted, nil
	}

}

// jsonify casts the given map structure into a format suitable for json.Marshal
func jsonify(in interface{}) (interface{}, error) {

	switch cast := in.(type) {

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

	case float64, float32, int64, int, string, bool:
		return cast, nil

	default:
		return nil, fmt.Errorf("unhandled jsonify value %v (%T)", cast, cast)

	}

}
