package kubernetes

import (
	"encoding/json"
	"io/ioutil"
	"os/exec"

	"github.com/Meetic/blackbeard/pkg/resource"
)

type ClusterRepository struct{}

func NewClusterRepository() resource.ClusterRepository {
	return &ClusterRepository{}
}

func (r ClusterRepository) GetVersion() (*resource.Version, error) {
	cmd := exec.Command("/bin/sh", "-c", "kubectl version --output json")
	cmdReader, _ := cmd.StdoutPipe()

	err := cmd.Start()

	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(cmdReader)

	if err != nil {
		return nil, err
	}

	var v resource.Version

	err = json.Unmarshal(data, &v)

	if err != nil {
		return nil, err
	}

	return &v, nil
}
