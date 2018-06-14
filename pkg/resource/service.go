package resource

// Service represent a kubernetes service
// Name is the service name
// Addr is the domain name where the service can be reach from outside the kubernetes cluster.
// For ingress exposed services it is the domain name declared in the ingress configuration
// for node port exposed services it is the ip / domain name of the cluster.
type Service struct {
	Name  string `json:"name"`
	Ports []Port `json:"ports"`
	Addr  string `json:"addr"`
}

// Port represent a kubernetes service port.
// This mean an internal port and a exposed port
type Port struct {
	Port        int32 `json:"port"`
	ExposedPort int32 `json:"exposedPort"`
}

// ServiceService defines the way kubernetes services are managed
type ServiceService interface {
	ListExposed(namespace string) ([]Service, error)
}

// ServiceRepository defines the way to interact with Kubernetes
type ServiceRepository interface {
	ListNodePort(n string) ([]Service, error)
	ListIngress(n string) ([]Service, error)
}

type serviceService struct {
	services ServiceRepository
}

// NewServiceService returns a ServicesService
func NewServiceService(services ServiceRepository) ServiceService {
	return &serviceService{
		services: services,
	}
}

// ListExposed find services exposed as NodePort and ingress configuration and return
// an array of services containing an URL, the exposed port and the service name.
func (ss *serviceService) ListExposed(namespace string) ([]Service, error) {

	var (
		services []Service
		err      error
	)

	services, err = ss.services.ListNodePort(namespace)
	if err != nil {
		return nil, err
	}

	ingress, err := ss.services.ListIngress(namespace)
	if err != nil {
		return nil, err
	}

	services = append(services, ingress...)

	return services, nil
}
