package resource

type Deployments []Deployment

type Deployment struct {
	Name   string
	Status DeploymentStatus
}

type DeploymentStatus string

const (
	DeploymentReady    DeploymentStatus = "Ready"
	DeploymentNotReady DeploymentStatus = "NotReady"
)

type DeploymentRepository interface {
	List(namespace string) (Deployments, error)
}
