package kubectl

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/Meetic/blackbeard/pkg/blackbeard"
)

const (
	timeout = 60 * time.Second
)

//NamespaceService is used to managed kubernetes namespace
type NamespaceConfigurationService struct {
	configPath string
}

//Ensure that NamespaceService implements the interface
var _ blackbeard.NamespaceConfigurationService = (*NamespaceConfigurationService)(nil)

//Apply load configuration files into kubernetes
func (ns *NamespaceConfigurationService) Apply(inv blackbeard.Inventory) error {

	err := execute(fmt.Sprintf("kubectl apply -f %s -n %s", filepath.Join(ns.configPath, inv.Namespace), inv.Namespace), timeout)
	if err != nil {
		return fmt.Errorf("the namespace could not be configured : %v", err)
	}

	return nil
}

func execute(c string, t time.Duration) error {

	cmd := exec.Command("/bin/sh", "-c", c)

	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd")
		return err
	}

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	//Start process. Exit code 127 if process fail to start.
	if err := cmd.Start(); err != nil {
		fmt.Fprintln(os.Stderr, "Error stating Cmd")
		return err
	}
	var timer *time.Timer
	if t > 0 {
		timer = time.NewTimer(t)
		var err error
		go func(timer *time.Timer, cmd *exec.Cmd) {
			for _ = range timer.C {
				e := cmd.Process.Kill()
				if e != nil {
					err = errors.New("the command has timeout but the process could not be killed")
				} else {
					err = errors.New("the command timed out")
				}
			}
		}(timer, cmd)
	}

	err = cmd.Wait()

	if t > 0 {
		timer.Stop()
	}

	if err != nil {
		return errors.New("the command did not succeed")
	}

	return nil
}
