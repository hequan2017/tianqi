package modeltraining

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	DatasetApi
	TrainingTaskApi
	ModelTestHistoryApi
}

var datasetService = service.ServiceGroupApp.ModeltrainingServiceGroup.DatasetService
var datasetVersionService = service.ServiceGroupApp.ModeltrainingServiceGroup.DatasetVersionService
var trainingTaskService = service.ServiceGroupApp.ModeltrainingServiceGroup.TrainingTaskService
var trainingParamService = service.ServiceGroupApp.ModeltrainingServiceGroup.TrainingParamService
var modelTestHistoryService = service.ServiceGroupApp.ModeltrainingServiceGroup.ModelTestHistoryService