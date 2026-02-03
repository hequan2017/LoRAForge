package cloud

type ServiceGroup struct {
	MirrorRepositoryService
	ComputeNodeService
	ProductSpecService
	InstanceService
	ImageService
	NetworkService
	VolumeService
	FineTuneTaskService
	InferenceTaskService
	SwiftWebUIService
}
