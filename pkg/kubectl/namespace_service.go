package kubectl

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/Meetic/blackbeard/pkg/blackbeard"
)

const (
	timeout = 30 * time.Second
)

//NamespaceService is used to managed kubernetes namespace
type NamespaceService struct {
	configPath string
}

//Ensure that NamespaceService implements the interface
var _ blackbeard.NamespaceService = (*NamespaceService)(nil)

//Create create a namespace
func (ns *NamespaceService) Create(inv blackbeard.Inventory) error {

	err := execute("kubectl create namespace "+inv.Namespace, timeout)
	if err != nil {
		return fmt.Errorf("The namespace %s could not be created because the either the namespace already exist or the command timed out : %s", inv.Namespace, err.Error())
	}

	return nil
}

//Apply load configuration files into kubernetes
func (ns *NamespaceService) Apply(inv blackbeard.Inventory) error {

	err := execute("kubectl apply -f "+ns.configPath+"/"+inv.Namespace+" -n "+inv.Namespace, 10*time.Second)
	if err != nil {
		return fmt.Errorf("The namespace could not be configured : %s", err.Error())
	}

	return nil
}

func execute(c string, t time.Duration) error {

	cmd := exec.Command("/bin/sh", "-c", c)

	//Start process. Exit code 127 if process fail to start.
	if err := cmd.Start(); err != nil {
		log.Println("error at start")
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
					err = errors.New("The command has timedout but the process could not be kiled.")
				} else {
					err = errors.New("The command timed out.")
				}
			}
		}(timer, cmd)
	}

	err := cmd.Wait()

	if t > 0 {
		timer.Stop()
	}

	if err != nil {
		return errors.New("The command did not succed")
	}

	return nil
}
