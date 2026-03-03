package modeltraining

import api "github.com/flipped-aurora/gin-vue-admin/server/api/v1"

type RouterGroup struct {
	DatasetRouter
	TrainingTaskRouter
	ModelTestHistoryRouter
}

var datasetApi = api.ApiGroupApp.ModeltrainingApiGroup.DatasetApi
var trainingTaskApi = api.ApiGroupApp.ModeltrainingApiGroup.TrainingTaskApi
var modelTestHistoryApi = api.ApiGroupApp.ModeltrainingApiGroup.ModelTestHistoryApi