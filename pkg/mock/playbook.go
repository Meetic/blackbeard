package mock

import (
	"encoding/json"
	"text/template"

	"github.com/Meetic/blackbeard/pkg/blackbeard"
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

func NewPlaybookRepository() blackbeard.PlaybookRepository {
	return &playbooks{}
}

func (p *playbooks) GetTemplate() ([]blackbeard.ConfigTemplate, error) {

	var templates []blackbeard.ConfigTemplate

	templates = append(templates, blackbeard.ConfigTemplate{
		Name:     "template.yml",
		Template: template.Must(template.New("tpl").Parse(tpl)),
	})

	return templates, nil
}

func (p *playbooks) GetDefault() (blackbeard.Inventory, error) {

	var inventory blackbeard.Inventory

	if err := json.Unmarshal([]byte(def), &inventory); err != nil {
		return blackbeard.Inventory{}, blackbeard.NewErrorReadingDefaultsFile(err)
	}

	return inventory, nil
}
