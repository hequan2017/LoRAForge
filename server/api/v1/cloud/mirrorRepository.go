package cloud

import (
	
	"github.com/flipped-aurora/gin-vue-admin/server/global"
    "github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
    "github.com/flipped-aurora/gin-vue-admin/server/model/cloud"
    cloudReq "github.com/flipped-aurora/gin-vue-admin/server/model/cloud/request"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

type MirrorRepositoryApi struct {}



// CreateMirrorRepository 创建镜像库
// @Tags MirrorRepository
// @Summary 创建镜像库
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body cloud.MirrorRepository true "创建镜像库"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /mirror/createMirrorRepository [post]
func (mirrorApi *MirrorRepositoryApi) CreateMirrorRepository(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	var mirror cloud.MirrorRepository
	err := c.ShouldBindJSON(&mirror)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = mirrorService.CreateMirrorRepository(ctx,&mirror)
	if err != nil {
        global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:" + err.Error(), c)
		return
	}
    response.OkWithMessage("创建成功", c)
}

// DeleteMirrorRepository 删除镜像库
// @Tags MirrorRepository
// @Summary 删除镜像库
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body cloud.MirrorRepository true "删除镜像库"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /mirror/deleteMirrorRepository [delete]
func (mirrorApi *MirrorRepositoryApi) DeleteMirrorRepository(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	ID := c.Query("ID")
	err := mirrorService.DeleteMirrorRepository(ctx,ID)
	if err != nil {
        global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteMirrorRepositoryByIds 批量删除镜像库
// @Tags MirrorRepository
// @Summary 批量删除镜像库
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /mirror/deleteMirrorRepositoryByIds [delete]
func (mirrorApi *MirrorRepositoryApi) DeleteMirrorRepositoryByIds(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	IDs := c.QueryArray("IDs[]")
	err := mirrorService.DeleteMirrorRepositoryByIds(ctx,IDs)
	if err != nil {
        global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateMirrorRepository 更新镜像库
// @Tags MirrorRepository
// @Summary 更新镜像库
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body cloud.MirrorRepository true "更新镜像库"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /mirror/updateMirrorRepository [put]
func (mirrorApi *MirrorRepositoryApi) UpdateMirrorRepository(c *gin.Context) {
    // 从ctx获取标准context进行业务行为
    ctx := c.Request.Context()

	var mirror cloud.MirrorRepository
	err := c.ShouldBindJSON(&mirror)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = mirrorService.UpdateMirrorRepository(ctx,mirror)
	if err != nil {
        global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindMirrorRepository 用id查询镜像库
// @Tags MirrorRepository
// @Summary 用id查询镜像库
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "用id查询镜像库"
// @Success 200 {object} response.Response{data=cloud.MirrorRepository,msg=string} "查询成功"
// @Router /mirror/findMirrorRepository [get]
func (mirrorApi *MirrorRepositoryApi) FindMirrorRepository(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	ID := c.Query("ID")
	remirror, err := mirrorService.GetMirrorRepository(ctx,ID)
	if err != nil {
        global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:" + err.Error(), c)
		return
	}
	response.OkWithData(remirror, c)
}
// GetMirrorRepositoryList 分页获取镜像库列表
// @Tags MirrorRepository
// @Summary 分页获取镜像库列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query cloudReq.MirrorRepositorySearch true "分页获取镜像库列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /mirror/getMirrorRepositoryList [get]
func (mirrorApi *MirrorRepositoryApi) GetMirrorRepositoryList(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	var pageInfo cloudReq.MirrorRepositorySearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := mirrorService.GetMirrorRepositoryInfoList(ctx,pageInfo)
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

// GetMirrorRepositoryPublic 不需要鉴权的镜像库接口
// @Tags MirrorRepository
// @Summary 不需要鉴权的镜像库接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /mirror/getMirrorRepositoryPublic [get]
func (mirrorApi *MirrorRepositoryApi) GetMirrorRepositoryPublic(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

    // 此接口不需要鉴权
    // 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
    mirrorService.GetMirrorRepositoryPublic(ctx)
    response.OkWithDetailed(gin.H{
       "info": "不需要鉴权的镜像库接口信息",
    }, "获取成功", c)
}
