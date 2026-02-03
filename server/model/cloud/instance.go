// 自动生成模板Instance
package cloud

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// 实例管理 结构体  Instance
type Instance struct {
	global.GVA_MODEL
	MirrorID        *int64  `json:"mirrorId" form:"mirrorId" gorm:"comment:镜像;column:mirror_id;"`               //镜像(旧)
	ImageName       *string `json:"imageName" form:"imageName" gorm:"comment:镜像名称;column:image_name;"`          //镜像名称(新)
	TemplateID      *int64  `json:"templateId" form:"templateId" gorm:"comment:模版;column:template_id;"`                            //模版
	UserID          *int64  `json:"userId" form:"userId" gorm:"comment:用户ID;column:user_id;"`                                      //用户ID
	NodeID          *int64  `json:"nodeId" form:"nodeId" gorm:"comment:节点;column:node_id;"`                     //节点
	DockerContainer *string `json:"dockerContainer" form:"dockerContainer" gorm:"comment:Docker容器;column:docker_container;"`       //Docker容器
	InstanceName    *string `json:"instanceName" form:"instanceName" gorm:"comment:实例名称;column:instance_name;"` //实例名称
	ContainerStatus *string `json:"containerStatus" form:"containerStatus" gorm:"comment:容器状态;column:container_status;"`           //容器状态

	// 容器配置参数
	Cpu          *float64 `json:"cpu" form:"cpu" gorm:"comment:CPU限制;column:cpu;"`                                      //CPU限制
	Memory       *int64   `json:"memory" form:"memory" gorm:"comment:内存限制(MB);column:memory;"`                          //内存限制(MB)
	GpuCount     *int64   `json:"gpuCount" form:"gpuCount" gorm:"comment:GPU数量;column:gpu_count;"`                      //GPU数量
	PortMapping  *string  `json:"portMapping" form:"portMapping" gorm:"type:text;comment:端口映射;column:port_mapping;"`    //端口映射
	VolumeMounts *string  `json:"volumeMounts" form:"volumeMounts" gorm:"type:text;comment:挂载目录;column:volume_mounts;"` //挂载目录
	EnvVars      *string  `json:"envVars" form:"envVars" gorm:"type:text;comment:环境变量;column:env_vars;"`                //环境变量
	Command      *string  `json:"command" form:"command" gorm:"comment:启动命令;column:command;"`                           //启动命令

	Remark *string `json:"remark" form:"remark" gorm:"comment:备注;column:remark;"` //备注
}

// TableName 实例管理 Instance自定义表名 instances
func (Instance) TableName() string {
	return "instances"
}
