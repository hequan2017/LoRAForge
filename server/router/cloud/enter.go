package cloud

import api "github.com/flipped-aurora/gin-vue-admin/server/api/v1"

type RouterGroup struct {
	MirrorRepositoryRouter
	ComputeNodeRouter
	ProductSpecRouter
	InstanceRouter
	ImageRouter
	NetworkRouter
	VolumeRouter
	FineTuneTaskRouter
	InferenceTaskRouter
	SwiftWebUIRouter
}

var (
	mirrorApi       = api.ApiGroupApp.CloudApiGroup.MirrorRepositoryApi
	nodeApi         = api.ApiGroupApp.CloudApiGroup.ComputeNodeApi
	specApi         = api.ApiGroupApp.CloudApiGroup.ProductSpecApi
	instApi         = api.ApiGroupApp.CloudApiGroup.InstanceApi
	fineTuneTaskApi = api.ApiGroupApp.CloudApiGroup.FineTuneTaskApi
	inferenceApi    = api.ApiGroupApp.CloudApiGroup.InferenceTaskApi
	SwiftWebUIApi   = api.ApiGroupApp.CloudApiGroup.SwiftWebUIApi
)
