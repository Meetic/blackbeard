package resource

type Statefulsets []Statefulset

type Statefulset struct {
	Name   string
	Status StatefulsetStatus
}

type StatefulsetStatus string

const (
	StatefulsetReady    StatefulsetStatus = "Ready"
	StatefulsetNotReady StatefulsetStatus = "NotReady"
)

type StatefulsetRepository interface {
	List(namespace string) (Statefulsets, error)
}
