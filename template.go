package tacks

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"

	"text/template"
	"github.com/spf13/cast"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Environments map[string]Environment
	Name         string
	Stack        Stack
	Version      string
	values       map[string]interface{}
}

type Environment struct {
	Ask             bool
	DeleteOnFailure bool
	Mode            string
	Pre             []Command
	Post            []Command
	Timeout         uint8
	Tags            map[string]string
	Variables       map[string]Variable
}

type Runtime struct {
	Environment string
	Variables   map[string]interface{}
}

type Variable struct {
	Cast     string
	Cmd      Command
	Constant string
	Env      Env
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func NewConfigFromReader(r io.Reader) (*Config, error) {

	var c Config

	in, err := ioutil.ReadAll(r)

	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(in, &c); err != nil {
		return nil, err
	}

	return &c, nil

}

func (c *Config) runHooks(label string, hooks []Command) error {

	for idx, cmd := range hooks {

		if out, err := cmd.Value(); err != nil {
			return err
		} else {
			logger.Debugf("Running %s-hook %d: %s -> %q", label, idx, cmd, out)
		}

	}

	return nil

}

func (c *Config) runStack(r Runtime) (string, error) {

	const null = ""

	tpl := template.New("_root")
	tpl.Funcs(Functions)

	if out, err := c.Stack.MarshalText(); err != nil {
		return null, err
	} else if tpl, err := tpl.Parse(string(out)); err != nil {
		return null, err
	} else {

		var buf bytes.Buffer

		tpl.ParseName = c.Name
		err := tpl.Execute(&buf, r)

		return buf.String(), err

	}

}

func (c *Config) Evaluate(environment string) (Environment, error) {

	const null = ""

	env, ok := c.Environments[environment]

	if !ok {
		return env, fmt.Errorf("no such environment %q", environment)
	}

	if err := c.runHooks("pre", env.Pre); err != nil {
		return env, err
	}

	var values = make(Values, len(env.Variables))
	c.values = make(map[string]interface{}, len(values))

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

	// evaluate and ...
	for key, value := range values.Evaluate() {

		logger.Debugf("Evaluated variable %q: %v", key, value)

		// return on first error
		if err := value.Error; err != nil {
			return env, err
		}

		// ... cast
		var (
			casted interface{}
			target = env.Variables[key].Cast
		)

		if target == null {
			casted = value.Value
		} else if target == "int" {
			casted = cast.ToInt(value.Value)
		} else {
			return env, fmt.Errorf("unknown cast target %q for variable %q", target, key)
		}

		if target != null {
			logger.Debugf("Casted variable %q to %q: %v (%T)", key, target, casted, casted)
		}

		c.values[key] = casted

	}

	runtime := Runtime{
		Environment: environment,
		Variables:   c.values, // TODO: remove from config state
	}

	if stack, err := c.runStack(runtime); err != nil {
		return env, err
	} else {
		fmt.Fprintln(os.Stderr, stack)
	}

	if err := c.runHooks("post", env.Post); err != nil {
		return env, err
	}

	return env, nil

}
