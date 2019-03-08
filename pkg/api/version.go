package api

import (
	"strings"

	"github.com/Meetic/blackbeard/pkg/version"
)

type Version struct {
	Blackbeard string `json:"blackbeard"`
	Kubernetes string `json:"kubernetes"`
	Kubectl    string `json:"kubectl"`
}

func (api *api) GetVersion() (*Version, error) {
	v, err := api.cluster.GetVersion()

	if err != nil {
		return nil, err
	}

	return &Version{
		Blackbeard: version.GetVersion(),
		Kubernetes: strings.Join([]string{v.ClientVersion.Major, v.ClientVersion.Minor}, "."),
		Kubectl:    strings.Join([]string{v.ServerVersion.Major, v.ServerVersion.Minor}, "."),
	}, nil
}
