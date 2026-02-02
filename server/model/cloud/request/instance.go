package request

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

type InstanceSearch struct {
	CreatedAtRange []time.Time `json:"createdAtRange" form:"createdAtRange[]"`
	InstanceName   string      `json:"instanceName" form:"instanceName"`
	MirrorId       *int        `json:"mirrorId" form:"mirrorId"`
	TemplateId     *int        `json:"templateId" form:"templateId"`
	NodeId         *int        `json:"nodeId" form:"nodeId"`
	request.PageInfo
}
