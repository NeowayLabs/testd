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
	daemon  *exec.Cmd
	output  chan []byte
	logFile string
	stopped bool
}

func (self *Testd) Stop() error {
	if self.stopped {
		return errors.New("already stopped")
	}
	err := self.daemon.Process.Kill()
	if err != nil {
		return err
	}
	self.stopped = true

	dirSeparatorIndex := strings.LastIndex(self.logFile, "/")
	if dirSeparatorIndex != -1 {
		logDir := self.logFile[0:dirSeparatorIndex]
		os.MkdirAll(logDir, os.ModeDir|os.ModePerm)
	}
	return ioutil.WriteFile(self.logFile, <-self.output, os.ModePerm)
}

func New(logFile string, name string, arg ...string) (*Testd, error) {
	self := Testd{
		output:  make(chan []byte),
		daemon:  exec.Command(name, arg...),
		logFile: logFile,
		stopped: false,
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
