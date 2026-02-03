
// 自动生成模板SwiftWebUI
package cloud
import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// Swift WebUI管理 结构体  SwiftWebUI
type SwiftWebUI struct {
    global.GVA_MODEL
  TaskName  *string `json:"taskName" form:"taskName" gorm:"comment:任务名称;column:task_name;size:255;" binding:"required"`  //任务名称
  NodeID  *int64 `json:"nodeId" form:"nodeId" gorm:"comment:节点ID;column:node_id;" binding:"required"`  //计算节点
  Language  *string `json:"language" form:"language" gorm:"default:zh;comment:语言;column:language;size:32;"`  //语言
  Port  *int64 `json:"port" form:"port" gorm:"default:7860;comment:端口;column:port;"`  //端口
  Status  *string `json:"status" form:"status" gorm:"default:pending;comment:状态;column:status;size:32;"`  //状态
  AccessUrl  *string `json:"accessUrl" form:"accessUrl" gorm:"comment:访问地址;column:access_url;size:255;"`  //访问地址
  ContainerId  *string `json:"containerId" form:"containerId" gorm:"comment:容器ID;column:container_id;size:64;"`  //容器ID
  Pid  *string `json:"pid" form:"pid" gorm:"comment:进程ID;column:pid;size:32;"`  //进程ID
}


// TableName Swift WebUI管理 SwiftWebUI自定义表名 swift_web_uis
func (SwiftWebUI) TableName() string {
    return "swift_web_uis"
}





