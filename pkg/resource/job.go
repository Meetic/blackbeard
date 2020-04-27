package resource

type JobService interface {
	Delete(namespace, resourceName string) error
}

type JobRepository interface {
	Delete(namespace, resourceName string) error
	List(namespace string) (Jobs, error)
}

type jobService struct {
	job JobRepository
}

type Jobs []Job

type Job struct {
	Name   string
	Status JobStatus
}

type JobStatus string

const (
	JobReady    JobStatus = "Ready"
	JobNotReady JobStatus = "NotReady"
)

func NewJobService(job JobRepository) JobService {
	return &jobService{
		job: job,
	}
}

func (js *jobService) Delete(namespace, resourceName string) error {
	return js.job.Delete(namespace, resourceName)
}
