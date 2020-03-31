package resource

type JobService interface {
	Delete(namespace, resourceName string) error
}

type JobRepository interface {
	Delete(namespace, resourceName string) error
}

type jobService struct {
	job JobRepository
}

func NewJobService(job JobRepository) JobService {
	return &jobService{
		job: job,
	}
}

func (js *jobService) Delete(namespace, resourceName string) error {
	return js.job.Delete(namespace, resourceName)
}
