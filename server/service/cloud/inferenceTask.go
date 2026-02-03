
package cloud

import (
	"context"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/cloud"
    cloudReq "github.com/flipped-aurora/gin-vue-admin/server/model/cloud/request"
)

type InferenceTaskService struct {}
// CreateInferenceTask 创建AI推理任务记录
// Author [yourname](https://github.com/yourname)
func (inferenceService *InferenceTaskService) CreateInferenceTask(ctx context.Context, inference *cloud.InferenceTask) (err error) {
	err = global.GVA_DB.Create(inference).Error
	return err
}

// DeleteInferenceTask 删除AI推理任务记录
// Author [yourname](https://github.com/yourname)
func (inferenceService *InferenceTaskService)DeleteInferenceTask(ctx context.Context, ID string) (err error) {
	err = global.GVA_DB.Delete(&cloud.InferenceTask{},"id = ?",ID).Error
	return err
}

// DeleteInferenceTaskByIds 批量删除AI推理任务记录
// Author [yourname](https://github.com/yourname)
func (inferenceService *InferenceTaskService)DeleteInferenceTaskByIds(ctx context.Context, IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]cloud.InferenceTask{},"id in ?",IDs).Error
	return err
}

// UpdateInferenceTask 更新AI推理任务记录
// Author [yourname](https://github.com/yourname)
func (inferenceService *InferenceTaskService)UpdateInferenceTask(ctx context.Context, inference cloud.InferenceTask) (err error) {
	err = global.GVA_DB.Model(&cloud.InferenceTask{}).Where("id = ?",inference.ID).Updates(&inference).Error
	return err
}

// GetInferenceTask 根据ID获取AI推理任务记录
// Author [yourname](https://github.com/yourname)
func (inferenceService *InferenceTaskService)GetInferenceTask(ctx context.Context, ID string) (inference cloud.InferenceTask, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&inference).Error
	return
}
// GetInferenceTaskInfoList 分页获取AI推理任务记录
// Author [yourname](https://github.com/yourname)
func (inferenceService *InferenceTaskService)GetInferenceTaskInfoList(ctx context.Context, info cloudReq.InferenceTaskSearch) (list []cloud.InferenceTask, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
    // 创建db
	db := global.GVA_DB.Model(&cloud.InferenceTask{})
    var inferences []cloud.InferenceTask
    // 如果有条件搜索 下方会自动创建搜索语句
    if len(info.CreatedAtRange) == 2 {
     db = db.Where("created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
    }
    
    if info.TaskName != nil && *info.TaskName != "" {
        db = db.Where("task_name LIKE ?", "%"+ *info.TaskName+"%")
    }
    if info.ModelPath != nil && *info.ModelPath != "" {
        db = db.Where("model_path LIKE ?", "%"+ *info.ModelPath+"%")
    }
    if info.Status != nil && *info.Status != "" {
        db = db.Where("status = ?", *info.Status)
    }
	err = db.Count(&total).Error
	if err!=nil {
    	return
    }

	if limit != 0 {
       db = db.Limit(limit).Offset(offset)
    }

	err = db.Find(&inferences).Error
	return  inferences, total, err
}
func (inferenceService *InferenceTaskService)GetInferenceTaskPublic(ctx context.Context) {
    // 此方法为获取数据源定义的数据
    // 请自行实现
}
