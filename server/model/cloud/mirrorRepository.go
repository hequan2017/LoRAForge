
// 自动生成模板MirrorRepository
package cloud
import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// 镜像库 结构体  MirrorRepository
type MirrorRepository struct {
    global.GVA_MODEL
  Name  *string `json:"name" form:"name" gorm:"comment:名字;column:name;" binding:"required"`  //名字
  Address  *string `json:"address" form:"address" gorm:"comment:地址;column:address;" binding:"required"`  //地址
  Description  *string `json:"description" form:"description" gorm:"comment:描述;column:description;"`  //描述
  Source  *string `json:"source" form:"source" gorm:"comment:来源;column:source;"`  //来源
  IsListed  *bool `json:"isListed" form:"isListed" gorm:"default:true;comment:是否上架;column:is_listed;" binding:"required"`  //是否上架
  Remark  *string `json:"remark" form:"remark" gorm:"comment:备注;column:remark;"`  //备注
}


// TableName 镜像库 MirrorRepository自定义表名 mirror_repositories
func (MirrorRepository) TableName() string {
    return "mirror_repositories"
}





