package tacks

import (
	"fmt"
	"io"
	"io/ioutil"
	"runtime"
	"strings"

	"gopkg.in/yaml.v2"
)

type Callback func(Document) error

type Template struct {
	Environments map[string]*Environment
	Name         string
	Stack        map[interface{}]interface{}
	Version      string
}

type Environment struct {
	Ask             bool
	DeleteOnFailure bool `yaml:"delete_on_failure"`
	Mode            string
	Name            string `yaml:"-"`
	Pre             []Command
	Post            []Command
	Region          string
	StackName       string
	Timeout         uint8
	Tags            map[string]string
	Variables       map[string]struct {
		Cmd      Command
		Constant string
		Env      Env
	}
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func NewTemplateFromReader(r io.Reader) (*Template, error) {

	const null = ""

	var t Template

	in, err := ioutil.ReadAll(r)

	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(in, &t); err != nil {
		return nil, err
	}

	for key, value := range t.Environments {

		if value.StackName == null {
			value.StackName = strings.Join([]string{
				key,
				t.Name,
				t.Version,
			}, "-")
		}

		value.Name = key

	}

	return &t, nil

}

func (t *Template) runHooks(label string, hooks []Command) error {

	for idx, cmd := range hooks {

		if out, err := cmd.Value(); err != nil {
			return err
		} else {
			logger.Debugf("Running %s-hook %d: %s -> %q", label, idx, cmd, out)
		}

	}

	return nil

}

func (t *Template) Evaluate(environment string, cb Callback) error {

	const null = ""

	env, ok := t.Environments[environment]

	if !ok {
		return fmt.Errorf("no such environment %q", environment)
	}

	if err := t.runHooks("pre", env.Pre); err != nil {
		return err
	}

	var (
		values    = make(Values, len(env.Variables))
		variables = make(map[string]interface{}, len(values))
	)

	for key, value := range env.Variables {

		logger.Debugf("Preparing variable %q", key)

		if con := value.Constant; con != null {
			values[key] = Constant(con)
		} else if cmd := value.Cmd; cmd != null {
			values[key] = Command(cmd)
		} else if env := value.Env; env != null {
			values[key] = Env(env)
		}

	}

	// evaluate
	for key, value := range values.Evaluate() {

		logger.Debugf("Evaluated variable %q: %v", key, value)

		// return on first error
		if err := value.Error; err != nil {
			return err
		}

		variables[key] = value.Value

	}

	document := Document{
		stack:       t.Stack,
		Environment: env,
		Variables:   variables,
	}

	logger.Infof("Running stack with variables %v", variables)

	if err := cb(document); err != nil {
		return err
	}

	if err := t.runHooks("post", env.Post); err != nil {
		return err
	}

	return nil

}
