package cloud

import api "github.com/flipped-aurora/gin-vue-admin/server/api/v1"

type RouterGroup struct {
	MirrorRepositoryRouter
	ComputeNodeRouter
	ProductSpecRouter
	InstanceRouter
}

var (
	mirrorApi = api.ApiGroupApp.CloudApiGroup.MirrorRepositoryApi
	nodeApi   = api.ApiGroupApp.CloudApiGroup.ComputeNodeApi
	specApi   = api.ApiGroupApp.CloudApiGroup.ProductSpecApi
	instApi   = api.ApiGroupApp.CloudApiGroup.InstanceApi
)
