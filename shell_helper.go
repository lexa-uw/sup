package sup

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	shlex "github.com/flynn/go-shlex"
)

func CommandOutput(rawCmd string) (string, error) {
	parts, err := shlex.Split(rawCmd)

	if err != nil {
		return "", err
	}

	cmd := prepareCommand(parts[0], parts[1:]...)
	rawOut, err := cmd.Output()
	if err != nil {
		return "", err
	}
	out := strings.TrimRight(string(rawOut), "\n")
	return out, err
}

func prepareCommand(bin string, args ...string) *exec.Cmd {
	cmd := exec.Command(bin, args...)
	if os.Getenv("DEBUG") != "" {
		fmt.Println("running >>", bin, args)
	}
	return cmd
}
