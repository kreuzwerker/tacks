package tacks

import (
	"fmt"
	"os/exec"
	"strings"
)

const NoCommandOutput = `""`

type Command string

func (c Command) Value() (interface{}, error) {

	const nothing = ""

	txt := func(out []byte) string {

		const nothing = ""

		if s := string(out); len(strings.TrimSpace(s)) == 0 {
			return NoCommandOutput
		} else {
			return strings.TrimRight(s, "\n")
		}

	}

	bash, err := exec.LookPath("bash")

	if err != nil {
		return nothing, fmt.Errorf("failed to find bash: %v", err)
	}

	cmd := exec.Command(bash,
		"-o",
		"pipefail",
		"-c",
		string(c))

	if out, err := cmd.CombinedOutput(); err != nil {
		return nothing, fmt.Errorf("failed to execute %q: %s (%v)", c, txt(out), err)
	} else {
		return txt(out), nil
	}

}
