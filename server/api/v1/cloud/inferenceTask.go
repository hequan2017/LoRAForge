package cloud

import (
	
	"github.com/flipped-aurora/gin-vue-admin/server/global"
    "github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
    "github.com/flipped-aurora/gin-vue-admin/server/model/cloud"
    cloudReq "github.com/flipped-aurora/gin-vue-admin/server/model/cloud/request"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

type InferenceTaskApi struct {}



// CreateInferenceTask 创建AI推理任务
// @Tags InferenceTask
// @Summary 创建AI推理任务
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body cloud.InferenceTask true "创建AI推理任务"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /inference/createInferenceTask [post]
func (inferenceApi *InferenceTaskApi) CreateInferenceTask(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	var inference cloud.InferenceTask
	err := c.ShouldBindJSON(&inference)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = inferenceService.CreateInferenceTask(ctx,&inference)
	if err != nil {
        global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:" + err.Error(), c)
		return
	}
    response.OkWithMessage("创建成功", c)
}

// DeleteInferenceTask 删除AI推理任务
// @Tags InferenceTask
// @Summary 删除AI推理任务
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body cloud.InferenceTask true "删除AI推理任务"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /inference/deleteInferenceTask [delete]
func (inferenceApi *InferenceTaskApi) DeleteInferenceTask(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	ID := c.Query("ID")
	err := inferenceService.DeleteInferenceTask(ctx,ID)
	if err != nil {
        global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteInferenceTaskByIds 批量删除AI推理任务
// @Tags InferenceTask
// @Summary 批量删除AI推理任务
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /inference/deleteInferenceTaskByIds [delete]
func (inferenceApi *InferenceTaskApi) DeleteInferenceTaskByIds(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	IDs := c.QueryArray("IDs[]")
	err := inferenceService.DeleteInferenceTaskByIds(ctx,IDs)
	if err != nil {
        global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateInferenceTask 更新AI推理任务
// @Tags InferenceTask
// @Summary 更新AI推理任务
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body cloud.InferenceTask true "更新AI推理任务"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /inference/updateInferenceTask [put]
func (inferenceApi *InferenceTaskApi) UpdateInferenceTask(c *gin.Context) {
    // 从ctx获取标准context进行业务行为
    ctx := c.Request.Context()

	var inference cloud.InferenceTask
	err := c.ShouldBindJSON(&inference)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = inferenceService.UpdateInferenceTask(ctx,inference)
	if err != nil {
        global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindInferenceTask 用id查询AI推理任务
// @Tags InferenceTask
// @Summary 用id查询AI推理任务
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "用id查询AI推理任务"
// @Success 200 {object} response.Response{data=cloud.InferenceTask,msg=string} "查询成功"
// @Router /inference/findInferenceTask [get]
func (inferenceApi *InferenceTaskApi) FindInferenceTask(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	ID := c.Query("ID")
	reinference, err := inferenceService.GetInferenceTask(ctx,ID)
	if err != nil {
        global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:" + err.Error(), c)
		return
	}
	response.OkWithData(reinference, c)
}
// GetInferenceTaskList 分页获取AI推理任务列表
// @Tags InferenceTask
// @Summary 分页获取AI推理任务列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query cloudReq.InferenceTaskSearch true "分页获取AI推理任务列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /inference/getInferenceTaskList [get]
func (inferenceApi *InferenceTaskApi) GetInferenceTaskList(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	var pageInfo cloudReq.InferenceTaskSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := inferenceService.GetInferenceTaskInfoList(ctx,pageInfo)
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

// GetInferenceTaskPublic 不需要鉴权的AI推理任务接口
// @Tags InferenceTask
// @Summary 不需要鉴权的AI推理任务接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /inference/getInferenceTaskPublic [get]
func (inferenceApi *InferenceTaskApi) GetInferenceTaskPublic(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

    // 此接口不需要鉴权
    // 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
    inferenceService.GetInferenceTaskPublic(ctx)
    response.OkWithDetailed(gin.H{
       "info": "不需要鉴权的AI推理任务接口信息",
    }, "获取成功", c)
}
