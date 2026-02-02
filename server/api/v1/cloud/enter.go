package cloud

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	MirrorRepositoryApi
	ComputeNodeApi
	ProductSpecApi
	InstanceApi
}

var (
	mirrorService = service.ServiceGroupApp.CloudServiceGroup.MirrorRepositoryService
	nodeService   = service.ServiceGroupApp.CloudServiceGroup.ComputeNodeService
	specService   = service.ServiceGroupApp.CloudServiceGroup.ProductSpecService
	instService   = service.ServiceGroupApp.CloudServiceGroup.InstanceService
)
