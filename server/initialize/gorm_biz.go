package initialize

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/computenode"
	"github.com/flipped-aurora/gin-vue-admin/server/model/imageregistry"
	"github.com/flipped-aurora/gin-vue-admin/server/model/instance"
	"github.com/flipped-aurora/gin-vue-admin/server/model/modeltraining"
	"github.com/flipped-aurora/gin-vue-admin/server/model/product"
)

func bizModel() error {
	db := global.GVA_DB
	err := db.AutoMigrate(imageregistry.ImageRegistry{}, computenode.ComputeNode{}, product.ProductSpec{}, instance.Instance{}, modeltraining.Dataset{}, modeltraining.DatasetVersion{}, modeltraining.TrainingTask{}, modeltraining.TrainingParam{}, modeltraining.ModelTestHistory{})
	if err != nil {
		return err
	}
	return nil
}
