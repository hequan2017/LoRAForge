
package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
)

type SwiftWebUISearch struct{
    CreatedAtRange []time.Time `json:"createdAtRange" form:"createdAtRange[]"`
      TaskName  *string `json:"taskName" form:"taskName"` 
    request.PageInfo
}
