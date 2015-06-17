package tacks

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
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

// name generates a stack name out of a number of non-empty parts
func name(parts ...string) string {

	const null = ""

	var notNull []string

	for _, value := range parts {

		if value != null {
			notNull = append(notNull, value)
		}

	}

	return strings.Join(notNull, "-")

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
			value.StackName = name(key, t.Name, t.Version)
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

	env_vars := make(map[string]string)
	env_vars["_REGION"] = env.Region
	env_vars["_NAME"] = env.Name
	env_vars["_MODE"] = env.Mode
	env_vars["_STACKNAME"] = env.StackName

	for key, value := range env.Tags {
		env_vars["_TAG_"+key] = value
	}

	logger.Debugf("Setting environment variables from stack metadata %q", env_vars)
	for key, value := range env_vars {
		os.Setenv(key, value)
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
