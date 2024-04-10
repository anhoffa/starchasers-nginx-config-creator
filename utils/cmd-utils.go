package utils

import (
	logging "github.com/ipfs/go-log/v2"
	"os"
	"os/exec"
)

var log = logging.Logger("utils")

func RedirectCmdOutput(name string, arg ...string) error {
	log.Infow("Running command", "command", name, "args", arg)

	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
