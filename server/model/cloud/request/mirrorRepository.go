
package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
)

type MirrorRepositorySearch struct{
    CreatedAtRange []time.Time `json:"createdAtRange" form:"createdAtRange[]"`
    Name           string      `json:"name" form:"name"`
    Source         string      `json:"source" form:"source"`
    IsListed       *bool       `json:"isListed" form:"isListed"`
    request.PageInfo
}
