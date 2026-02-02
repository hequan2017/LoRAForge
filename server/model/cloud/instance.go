
// 自动生成模板Instance
package cloud
import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// 实例管理 结构体  Instance
type Instance struct {
    global.GVA_MODEL
  MirrorID  *int64 `json:"mirrorId" form:"mirrorId" gorm:"comment:镜像;column:mirror_id;" binding:"required"`  //镜像
  TemplateID  *int64 `json:"templateId" form:"templateId" gorm:"comment:模版;column:template_id;" binding:"required"`  //模版
  UserID  *int64 `json:"userId" form:"userId" gorm:"comment:用户ID;column:user_id;"`  //用户ID
  NodeID  *int64 `json:"nodeId" form:"nodeId" gorm:"comment:节点;column:node_id;" binding:"required"`  //节点
  DockerContainer  *string `json:"dockerContainer" form:"dockerContainer" gorm:"comment:Docker容器;column:docker_container;"`  //Docker容器
  InstanceName  *string `json:"instanceName" form:"instanceName" gorm:"comment:实例名称;column:instance_name;" binding:"required"`  //实例名称
  ContainerStatus  *string `json:"containerStatus" form:"containerStatus" gorm:"comment:容器状态;column:container_status;"`  //容器状态
  Remark  *string `json:"remark" form:"remark" gorm:"comment:备注;column:remark;"`  //备注
}


// TableName 实例管理 Instance自定义表名 instances
func (Instance) TableName() string {
    return "instances"
}





