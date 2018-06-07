package osops

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
)

type Config struct {
	Log logrus.FieldLogger
}

func New(c Config) *osOperations {
	return &osOperations{
		Config: c,
	}
}

//go:generate mockery -name OSOperations -inpkg

// OSOperations operations against systemd
type OSOperations interface {
	// ExecCommand will execute the command
	ExecCommand(command string, args ...string) (string, error)
}

type osOperations struct {
	Config
}

func (o *osOperations) ExecCommand(command string, args ...string) (string, error) {
	out, err := exec.Command(command, args...).CombinedOutput()
	output := strings.TrimSpace(string(out))
	if err != nil {
		return output, fmt.Errorf("failed executing %s %v , error %s", command, args, err)
	}
	o.Log.Debug("Command executed", "command", command, "arguments", args, "output", output)
	return output, err
}
