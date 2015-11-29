package testd

import (
	"io"
	"os/exec"
	"testing"
)

type Testd struct {
	t      *testing.T
	daemon *exec.Cmd
	output chan string
}

func (self *Testd) Stop() error {
	err := self.daemon.Process.Kill()
	if err != nil {
		return err
	}
	//output := <-self.output
	//TODO: Save logs
	return nil
}

func New(t *testing.T, logsBaseDir string, name string, arg ...string) (*Testd, error) {
	self := Testd{
		t:      t,
		output: make(chan string),
		daemon: exec.Command(name, arg...),
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

func readDaemonOutput(stdout io.Reader, stderr io.Reader, output chan<- string) {
	out := make([]byte, 0)
	errOutput := make([]byte, 0)

	readMoreData := func(target *[]byte, source io.Reader) bool {
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
		if !readMoreData(&out, stdout) &&
			!readMoreData(&errOutput, stderr) {
			break
		}
	}
	output <- "stdout:" + string(out) + "\n\nstderr:" + string(errOutput)
}
