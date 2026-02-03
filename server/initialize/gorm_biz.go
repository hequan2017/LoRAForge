package initialize

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/cloud"
)

func bizModel() error {
	db := global.GVA_DB
	err := db.AutoMigrate(
		cloud.MirrorRepository{},
		cloud.ComputeNode{},
		cloud.ProductSpec{},
		cloud.Instance{},
		cloud.FineTuneTask{},
		cloud.FineTuneTaskSnapshot{},
	)
	if err != nil {
		return err
	}
	return nil
}
