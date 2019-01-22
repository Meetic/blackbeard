package resource

type clusterService struct {
	client ClusterRepository
}

type ClusterRepository interface {
	GetVersion() (*Version, error)
}

type ClusterService interface {
	GetVersion() (*Version, error)
}

func NewClusterService(client ClusterRepository) ClusterService {
	return &clusterService{
		client: client,
	}
}

type Version struct {
	ClientVersion struct {
		Major string `json:"major"`
		Minor string `json:"minor"`
	} `json:"clientVersion"`
	ServerVersion struct {
		Major string `json:"major"`
		Minor string `json:"minor"`
	} `json:"serverVersion"`
}

func (cs *clusterService) GetVersion() (*Version, error) {
	return cs.client.GetVersion()
}
