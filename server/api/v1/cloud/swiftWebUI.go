package cloud

import (
	
	"github.com/flipped-aurora/gin-vue-admin/server/global"
    "github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
    "github.com/flipped-aurora/gin-vue-admin/server/model/cloud"
    cloudReq "github.com/flipped-aurora/gin-vue-admin/server/model/cloud/request"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

type SwiftWebUIApi struct {}



// CreateSwiftWebUI 创建Swift WebUI管理
// @Tags SwiftWebUI
// @Summary 创建Swift WebUI管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body cloud.SwiftWebUI true "创建Swift WebUI管理"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /SwiftWebUI/createSwiftWebUI [post]
func (SwiftWebUIApi *SwiftWebUIApi) CreateSwiftWebUI(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	var SwiftWebUI cloud.SwiftWebUI
	err := c.ShouldBindJSON(&SwiftWebUI)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = SwiftWebUIService.CreateSwiftWebUI(ctx,&SwiftWebUI)
	if err != nil {
        global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:" + err.Error(), c)
		return
	}
    response.OkWithMessage("创建成功", c)
}

// DeleteSwiftWebUI 删除Swift WebUI管理
// @Tags SwiftWebUI
// @Summary 删除Swift WebUI管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body cloud.SwiftWebUI true "删除Swift WebUI管理"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /SwiftWebUI/deleteSwiftWebUI [delete]
func (SwiftWebUIApi *SwiftWebUIApi) DeleteSwiftWebUI(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	ID := c.Query("ID")
	err := SwiftWebUIService.DeleteSwiftWebUI(ctx,ID)
	if err != nil {
        global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteSwiftWebUIByIds 批量删除Swift WebUI管理
// @Tags SwiftWebUI
// @Summary 批量删除Swift WebUI管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /SwiftWebUI/deleteSwiftWebUIByIds [delete]
func (SwiftWebUIApi *SwiftWebUIApi) DeleteSwiftWebUIByIds(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	IDs := c.QueryArray("IDs[]")
	err := SwiftWebUIService.DeleteSwiftWebUIByIds(ctx,IDs)
	if err != nil {
        global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateSwiftWebUI 更新Swift WebUI管理
// @Tags SwiftWebUI
// @Summary 更新Swift WebUI管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body cloud.SwiftWebUI true "更新Swift WebUI管理"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /SwiftWebUI/updateSwiftWebUI [put]
func (SwiftWebUIApi *SwiftWebUIApi) UpdateSwiftWebUI(c *gin.Context) {
    // 从ctx获取标准context进行业务行为
    ctx := c.Request.Context()

	var SwiftWebUI cloud.SwiftWebUI
	err := c.ShouldBindJSON(&SwiftWebUI)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = SwiftWebUIService.UpdateSwiftWebUI(ctx,SwiftWebUI)
	if err != nil {
        global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindSwiftWebUI 用id查询Swift WebUI管理
// @Tags SwiftWebUI
// @Summary 用id查询Swift WebUI管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "用id查询Swift WebUI管理"
// @Success 200 {object} response.Response{data=cloud.SwiftWebUI,msg=string} "查询成功"
// @Router /SwiftWebUI/findSwiftWebUI [get]
func (SwiftWebUIApi *SwiftWebUIApi) FindSwiftWebUI(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	ID := c.Query("ID")
	reSwiftWebUI, err := SwiftWebUIService.GetSwiftWebUI(ctx,ID)
	if err != nil {
        global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:" + err.Error(), c)
		return
	}
	response.OkWithData(reSwiftWebUI, c)
}
// GetSwiftWebUIList 分页获取Swift WebUI管理列表
// @Tags SwiftWebUI
// @Summary 分页获取Swift WebUI管理列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query cloudReq.SwiftWebUISearch true "分页获取Swift WebUI管理列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /SwiftWebUI/getSwiftWebUIList [get]
func (SwiftWebUIApi *SwiftWebUIApi) GetSwiftWebUIList(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	var pageInfo cloudReq.SwiftWebUISearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := SwiftWebUIService.GetSwiftWebUIInfoList(ctx,pageInfo)
	if err != nil {
	    global.GVA_LOG.Error("获取失败!", zap.Error(err))
        response.FailWithMessage("获取失败:" + err.Error(), c)
        return
    }
    response.OkWithDetailed(response.PageResult{
        List:     list,
        Total:    total,
        Page:     pageInfo.Page,
        PageSize: pageInfo.PageSize,
    }, "获取成功", c)
}
// GetSwiftWebUIDataSource 获取SwiftWebUI的数据源
// @Tags SwiftWebUI
// @Summary 获取SwiftWebUI的数据源
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "查询成功"
// @Router /SwiftWebUI/getSwiftWebUIDataSource [get]
func (SwiftWebUIApi *SwiftWebUIApi) GetSwiftWebUIDataSource(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

    // 此接口为获取数据源定义的数据
    dataSource, err := SwiftWebUIService.GetSwiftWebUIDataSource(ctx)
    if err != nil {
        global.GVA_LOG.Error("查询失败!", zap.Error(err))
   		response.FailWithMessage("查询失败:" + err.Error(), c)
   		return
    }
   response.OkWithData(dataSource, c)
}

// GetSwiftWebUIPublic 不需要鉴权的Swift WebUI管理接口
// @Tags SwiftWebUI
// @Summary 不需要鉴权的Swift WebUI管理接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /SwiftWebUI/getSwiftWebUIPublic [get]
func (SwiftWebUIApi *SwiftWebUIApi) GetSwiftWebUIPublic(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

    // 此接口不需要鉴权
    // 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
    SwiftWebUIService.GetSwiftWebUIPublic(ctx)
    response.OkWithDetailed(gin.H{
       "info": "不需要鉴权的Swift WebUI管理接口信息",
    }, "获取成功", c)
}
