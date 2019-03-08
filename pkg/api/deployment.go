package api

func (api *api) Kill(namespace string, deployments []string) []error {

	var errors []error

	for _, d := range deployments {
		pod, err := api.Pods().Find(namespace, d)
		if err != nil {
			errors = append(errors, err)
			continue
		}

		if err := api.Pods().Delete(namespace, pod); err != nil {
			errors = append(errors, err)
			continue
		}
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}
