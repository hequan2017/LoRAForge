
// 自动生成模板ComputeNode
package cloud
import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// 算力节点 结构体  ComputeNode
type ComputeNode struct {
    global.GVA_MODEL
  Name  *string `json:"name" form:"name" gorm:"comment:名字;column:name;" binding:"required"`  //名字
  Region  *string `json:"region" form:"region" gorm:"comment:区域;column:region;"`  //区域
  CPU  *string `json:"cpu" form:"cpu" gorm:"comment:CPU;column:cpu;"`  //CPU
  Memory  *string `json:"memory" form:"memory" gorm:"comment:内存;column:memory;"`  //内存
  SystemDiskCapacity  *int64 `json:"systemDiskCapacity" form:"systemDiskCapacity" gorm:"comment:系统盘容量;column:system_disk_capacity;"`  //系统盘容量
  DataDiskCapacity  *int64 `json:"dataDiskCapacity" form:"dataDiskCapacity" gorm:"comment:数据盘容量;column:data_disk_capacity;"`  //数据盘容量
  PublicIP  *string `json:"publicIp" form:"publicIp" gorm:"comment:公网IP;column:public_ip;" binding:"required"`  //公网IP
  PrivateIP  *string `json:"privateIp" form:"privateIp" gorm:"comment:内网IP;column:private_ip;" binding:"required"`  //内网IP
  SSHPort  *int64 `json:"sshPort" form:"sshPort" gorm:"default:22;comment:SSH端口;column:ssh_port;" binding:"required"`  //SSH端口
  Username  *string `json:"username" form:"username" gorm:"comment:用户名;column:username;"`  //用户名
  Password  *string `json:"password" form:"password" gorm:"comment:密码;column:password;"`  //密码
  GPUName  *string `json:"gpuName" form:"gpuName" gorm:"comment:显卡名称;column:gpu_name;"`  //显卡名称
  GPUCount  *int64 `json:"gpuCount" form:"gpuCount" gorm:"comment:显卡数量;column:gpu_count;"`  //显卡数量
  DockerConnectAddress  *string `json:"dockerConnectAddress" form:"dockerConnectAddress" gorm:"comment:Docker连接地址;column:docker_connect_address;"`  //Docker连接地址
  UseTLS  *bool `json:"useTls" form:"useTls" gorm:"default:true;comment:使用TLS;column:use_tls;"`  //使用TLS
  CACertificate  *string `json:"caCertificate" form:"caCertificate" gorm:"comment:CA证书;column:ca_certificate;size:text;"`  //CA证书
  ClientCertificate  *string `json:"clientCertificate" form:"clientCertificate" gorm:"comment:客户端证书;column:client_certificate;size:text;"`  //客户端证书
  ClientKey  *string `json:"clientKey" form:"clientKey" gorm:"comment:客户端私钥;column:client_key;size:text;"`  //客户端私钥
  IsListed  *bool `json:"isListed" form:"isListed" gorm:"default:true;comment:是否上架;column:is_listed;" binding:"required"`  //是否上架
  Remark  *string `json:"remark" form:"remark" gorm:"comment:备注;column:remark;"`  //备注
}


// TableName 算力节点 ComputeNode自定义表名 compute_nodes
func (ComputeNode) TableName() string {
    return "compute_nodes"
}





