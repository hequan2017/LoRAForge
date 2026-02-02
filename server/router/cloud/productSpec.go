package cloud

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type ProductSpecRouter struct {}

// InitProductSpecRouter 初始化 产品规格 路由信息
func (s *ProductSpecRouter) InitProductSpecRouter(Router *gin.RouterGroup,PublicRouter *gin.RouterGroup) {
	specRouter := Router.Group("spec").Use(middleware.OperationRecord())
	specRouterWithoutRecord := Router.Group("spec")
	specRouterWithoutAuth := PublicRouter.Group("spec")
	{
		specRouter.POST("createProductSpec", specApi.CreateProductSpec)   // 新建产品规格
		specRouter.DELETE("deleteProductSpec", specApi.DeleteProductSpec) // 删除产品规格
		specRouter.DELETE("deleteProductSpecByIds", specApi.DeleteProductSpecByIds) // 批量删除产品规格
		specRouter.PUT("updateProductSpec", specApi.UpdateProductSpec)    // 更新产品规格
	}
	{
		specRouterWithoutRecord.GET("findProductSpec", specApi.FindProductSpec)        // 根据ID获取产品规格
		specRouterWithoutRecord.GET("getProductSpecList", specApi.GetProductSpecList)  // 获取产品规格列表
	}
	{
	    specRouterWithoutAuth.GET("getProductSpecPublic", specApi.GetProductSpecPublic)  // 产品规格开放接口
	}
}
