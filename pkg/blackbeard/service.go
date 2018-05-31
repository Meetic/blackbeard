package blackbeard

//Service represent a kubernetes service
type Service struct {
	Name  string `json:"name"`
	Ports []Port `json:"ports"`
	Addr  string `json:"addr"`
}

//Port represent a kubernetes service port.
//This mean an internal port and a exposed port
type Port struct {
	Port        int32 `json:"port"`
	ExposedPort int32 `json:"exposedPort"`
}

type ServiceService interface {
	ListExposed(string) ([]Service, error)
}

type ServiceRepository interface {
	ListNodePort(string) ([]Service, error)
	ListIngress(string) ([]Service, error)
}

type serviceService struct {
	services ServiceRepository
}

func NewServiceService(services ServiceRepository) ServiceService {
	return &serviceService{
		services: services,
	}
}

//ListExposed find services exposed as NodePort and ingress configuration and return
//an array of services containing an URL, the exposed port and the service name.
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
