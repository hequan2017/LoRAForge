
package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
)

type InferenceTaskSearch struct{
    CreatedAtRange []time.Time `json:"createdAtRange" form:"createdAtRange[]"`
      TaskName  *string `json:"taskName" form:"taskName"` 
      ModelPath  *string `json:"modelPath" form:"modelPath"` 
      Status  *string `json:"status" form:"status"` 
    request.PageInfo
}
