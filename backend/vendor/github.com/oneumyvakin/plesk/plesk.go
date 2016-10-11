package plesk

import (
    "log"
    "os/exec"
    "syscall"
)

type Plesk struct {
	Config map[string]string
	Log    *log.Logger
}

func New(log *log.Logger) (plesk Plesk, err error) {

	plesk = Plesk{
		Log: log,
	}

	plesk.Config, err = plesk.getSettings()
	return
}

func execute(log *log.Logger, command string, args ...string) (output string, outputBytes []byte, code int, err error) {
	log.Printf("%s %s", command, args)

	cmd := exec.Command(command, args...)
	var waitStatus syscall.WaitStatus

	if outputBytes, err = cmd.CombinedOutput(); err != nil {
		// Did the command fail because of an unsuccessful exit code
		if exitError, ok := err.(*exec.ExitError); ok {
			waitStatus = exitError.Sys().(syscall.WaitStatus)
			code = waitStatus.ExitStatus()
		}
	} else {
		// Command was successful
		waitStatus = cmd.ProcessState.Sys().(syscall.WaitStatus)
		code = waitStatus.ExitStatus()
	}

	output = string(outputBytes)

	log.Println("output: ", output, "err: ", err, "code: ", code)

	return
}