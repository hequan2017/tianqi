package modeltraining

type ServiceGroup struct {
	DatasetService
	DatasetVersionService
	TrainingTaskService
	TrainingParamService
	ModelTestHistoryService
}