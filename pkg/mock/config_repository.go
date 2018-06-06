package mock

import (
	"text/template"

	"github.com/Meetic/blackbeard/pkg/blackbeard"
)

type configRepository struct{}

const (
	tpl = `
	{{range .Values.microservices}}
  ---
  kind: Deployment
  apiVersion: extensions/v1beta1
  metadata:
  name: {{.name}}
  spec:
  replicas: 1
  template:
  metadata:
  labels:
  app: fpm-{{.name}}
  spec:
  containers:
  - name: {{.name}}
  image:  docker.io/{{.name}}:{{.version}}
{{end}}
`
)

// NewConfigRepository returns a new Mock ConfigRepository
func NewConfigRepository() blackbeard.ConfigRepository {
	return &configRepository{}
}

func (cr *configRepository) GetTemplate() ([]blackbeard.ConfigTemplate, error) {

	var templates []blackbeard.ConfigTemplate

	templates = append(templates, blackbeard.ConfigTemplate{
		Name:     "template.yml",
		Template: template.Must(template.New("tpl").Parse(tpl)),
	})

	return templates, nil
}

func (cr *configRepository) Save(namespace string, configs []blackbeard.Config) error {
	return nil
}

func (cr *configRepository) Delete(namespace string) error {
	return nil
}
