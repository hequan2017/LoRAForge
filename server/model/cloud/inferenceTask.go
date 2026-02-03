
// 自动生成模板InferenceTask
package cloud
import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// AI推理任务 结构体  InferenceTask
type InferenceTask struct {
    global.GVA_MODEL
  TaskName  *string `json:"taskName" form:"taskName" gorm:"comment:任务名称;column:task_name;size:100;" binding:"required"`  //任务名称
  ModelPath  *string `json:"modelPath" form:"modelPath" gorm:"comment:模型路径;column:model_path;size:255;" binding:"required"`  //模型路径
  Prompt  *string `json:"prompt" form:"prompt" gorm:"comment:正向提示词;column:prompt;size:text;" binding:"required"`  //正向提示词
  NegativePrompt  *string `json:"negativePrompt" form:"negativePrompt" gorm:"comment:反向提示词;column:negative_prompt;size:text;"`  //反向提示词
  Steps  *int64 `json:"steps" form:"steps" gorm:"default:20;comment:采样步数;column:steps;" binding:"required"`  //采样步数
  CfgScale  *float64 `json:"cfgScale" form:"cfgScale" gorm:"default:7.0;comment:引导系数;column:cfg_scale;" binding:"required"`  //引导系数
  Seed  *int64 `json:"seed" form:"seed" gorm:"default:-1;comment:随机种子;column:seed;" binding:"required"`  //随机种子
  Width  *int64 `json:"width" form:"width" gorm:"default:512;comment:宽度;column:width;" binding:"required"`  //宽度
  Height  *int64 `json:"height" form:"height" gorm:"default:512;comment:高度;column:height;" binding:"required"`  //高度
  Sampler  *string `json:"sampler" form:"sampler" gorm:"default:Euler a;comment:采样器;column:sampler;size:50;" binding:"required"`  //采样器
  Status  *string `json:"status" form:"status" gorm:"comment:状态;column:status;size:20;" binding:"required"`  //状态
  ResultImage  string `json:"resultImage" form:"resultImage" gorm:"comment:生成结果;column:result_image;"`  //生成结果
  GenerationTime  *float64 `json:"generationTime" form:"generationTime" gorm:"comment:耗时(秒);column:generation_time;"`  //耗时(秒)
}


// TableName AI推理任务 InferenceTask自定义表名 cloud_inference_tasks
func (InferenceTask) TableName() string {
    return "cloud_inference_tasks"
}





