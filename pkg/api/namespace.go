package api

//Namespace represents a kubernetes namespace enrich with informations from the playbook.
type Namespace struct {
	//Name is the namespace name
	Name string
	//Phase is the namespace status phase. It could be "active" or "terminating"
	Phase string
	//Status is the namespace status. It is a percentage of runnning pods vs all pods in the namespace.
	Status int
	//Managed is true if the namespace as an associated inventory on the current playbook. False if not.
	Managed bool
}

// ListNamespaces returns a list of Namespace.
// For each kubernetes namespace, it checks if an associated inventory exists.
func (api *api) ListNamespaces() ([]Namespace, error) {
	nsList, err := api.namespaces.List()
	if err != nil {
		return nil, err
	}

	var namespaces []Namespace

	for _, ns := range nsList {

		namespace := Namespace{
			Name:    ns.Name,
			Phase:   ns.Phase,
			Status:  ns.Status,
			Managed: false,
		}

		if api.inventories.Exists(ns.Name) {
			namespace.Managed = true
		}

		namespaces = append(namespaces, namespace)
	}

	return namespaces, nil

}
