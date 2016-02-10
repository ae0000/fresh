package runner

import (
	"io"
	"os/exec"
	"runtime"
)

func run() bool {
	runnerLog("Running...")

	// Notify only on linux
	if runtime.GOOS == "linux" {
		notify := exec.Command("notify-send", "--expire-time=10", "--urgency=low", "Built")
		notify.Run()
	}

	cmd := exec.Command(buildPath())

	stderr, err := cmd.StderrPipe()
	if err != nil {
		fatal(err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fatal(err)
	}

	err = cmd.Start()
	if err != nil {
		fatal(err)
	}

	go io.Copy(appLogWriter{}, stderr)
	go io.Copy(appLogWriter{}, stdout)

	go func() {
		<-stopChannel
		pid := cmd.Process.Pid
		runnerLog("Killing PID %d", pid)
		cmd.Process.Kill()
	}()

	return true
}
