package mock

import (
	"encoding/json"
	"text/template"

	"github.com/Meetic/blackbeard/pkg/playbook"
)

const (
	def = `{
  "namespace": "default",
  "values": {
    "microservices": [
      {
        "name": "api-advertising",
        "version": "latest",
        "urls": [
          "api-advertising"
        ]
      },
      {
        "name": "api-algo",
        "version": "latest",
        "urls": [
          "api-algo"
        ]
      }
    ]
  }
}`

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

type playbooks struct{}

func NewPlaybookRepository() playbook.PlaybookRepository {
	return &playbooks{}
}

func (p *playbooks) GetTemplate() ([]playbook.ConfigTemplate, error) {

	var templates []playbook.ConfigTemplate

	templates = append(templates, playbook.ConfigTemplate{
		Name:     "template.yml",
		Template: template.Must(template.New("tpl").Parse(tpl)),
	})

	return templates, nil
}

func (p *playbooks) GetDefault() (playbook.Inventory, error) {

	var inventory playbook.Inventory

	if err := json.Unmarshal([]byte(def), &inventory); err != nil {
		return playbook.Inventory{}, playbook.NewErrorReadingDefaultsFile(err)
	}

	return inventory, nil
}
