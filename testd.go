package testd

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type Testd struct {
	daemon      *exec.Cmd
	output      chan []byte
	logFilePath string
	stopped     bool
}

// Stops the daemon, if it is not already not stopped
func (self *Testd) Stop() error {
	if self.stopped {
		return errors.New("already stopped")
	}
	self.daemon.Process.Kill()
	self.stopped = true

	dirSeparatorIndex := strings.LastIndex(self.logFilePath, "/")
	if dirSeparatorIndex != -1 {
		logDir := self.logFilePath[0:dirSeparatorIndex]
		os.MkdirAll(logDir, os.ModeDir|os.ModePerm)
	}
	return ioutil.WriteFile(self.logFilePath, <-self.output, os.ModePerm)
}

// New creates a new testd instance. It will call exec.Command using the
// given name and arg and start the daemon.
// All logs (stdout and stderr) of the daemon will be saved at the given logFilePath.
// If a directory on logFilePath does not exist, it will attempt to
// create the directory for you.
// You must always call Testd.Stop, even when you know that the daemon already exited for some reason.
func New(logFilePath string, name string, arg ...string) (*Testd, error) {
	self := Testd{
		output:      make(chan []byte),
		daemon:      exec.Command(name, arg...),
		logFilePath: logFilePath,
		stopped:     false,
	}

	stdout, err := self.daemon.StdoutPipe()
	if err != nil {
		return nil, err
	}

	stderr, err := self.daemon.StderrPipe()
	if err != nil {
		return nil, err
	}

	err = self.daemon.Start()
	if err != nil {
		return nil, err
	}
	go readDaemonOutput(stdout, stderr, self.output)
	return &self, nil
}

func readDaemonOutput(stdout io.Reader, stderr io.Reader, output chan<- []byte) {
	out := make([]byte, 0)
	errOutput := make([]byte, 0)

	readOutput := func(target *[]byte, source io.Reader) bool {
		r := make([]byte, 1024)
		n, err := source.Read(r)
		if n == 0 {
			if err != nil {
				return false
			}
			return true
		}
		*target = append(*target, r...)
		return true
	}

	for {
		if !readOutput(&out, stdout) &&
			!readOutput(&errOutput, stderr) {
			break
		}
	}

	output <- append(out, errOutput...)
}
