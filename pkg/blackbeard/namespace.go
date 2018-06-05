package blackbeard

const (
	running = "Running"
)

//NamespaceService defined the way namespace should be managed.
type NamespaceService interface {
	Create(string) error
	ApplyConfig(string, string) error
	Delete(string) error
	GetStatus(string) (int, error)
	GetPods(namespace string) (Pods, error)
}

type NamespaceRepository interface {
	Create(string) error
	ApplyConfig(string, string) error
	Delete(string) error
}

type namespaceService struct {
	namespaces NamespaceRepository
	pods       PodRepository
}

func NewNamespaceService(namespaces NamespaceRepository, pods PodRepository) NamespaceService {
	return &namespaceService{
		namespaces: namespaces,
		pods:       pods,
	}
}

func (ns *namespaceService) Create(n string) error {
	return ns.namespaces.Create(n)
}

func (ns *namespaceService) ApplyConfig(namespace, configPath string) error {
	return ns.namespaces.ApplyConfig(namespace, configPath)
}

func (ns *namespaceService) Delete(namespace string) error {
	return ns.namespaces.Delete(namespace)
}

func (ns *namespaceService) GetStatus(namespace string) (int, error) {

	pods, err := ns.pods.List(namespace)
	if err != nil {
		return 0, err
	}

	totalPods := len(pods)

	if totalPods == 0 {
		return 0, nil
	}

	var i int

	for _, pod := range pods {
		if pod.Status == running {
			i++
		}
	}

	status := i * 100 / totalPods

	return status, nil

}

func (ns *namespaceService) GetPods(namespace string) (Pods, error) {
	return ns.pods.List(namespace)
}
