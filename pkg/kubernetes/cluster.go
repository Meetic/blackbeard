package kubernetes

import (
	"encoding/json"
	"os/exec"

	"github.com/Meetic/blackbeard/pkg/resource"
)

type ClusterRepository struct{}

func NewClusterRepository() resource.ClusterRepository {
	return &ClusterRepository{}
}

func (r ClusterRepository) GetVersion() (*resource.Version, error) {
	cmd := exec.Command("/bin/sh", "-c", "kubectl version --output json")
	result, err := cmd.Output()

	if err != nil {
		return nil, err
	}

	var v resource.Version

	err = json.Unmarshal(result, &v)

	if err != nil {
		return nil, err
	}

	return &v, nil
}
