package initialize

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/task"

	"github.com/robfig/cron/v3"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

func Timer() {
	go func() {
		var option []cron.Option
		option = append(option, cron.WithSeconds())
		// 清理DB定时任务
		_, err := global.GVA_Timer.AddTaskByFunc("ClearDB", "@daily", func() {
			err := task.ClearTable(global.GVA_DB) // 定时任务方法定在task文件包中
			if err != nil {
				fmt.Println("timer error:", err)
			}
		}, "定时清理数据库【日志，黑名单】内容", option...)
		if err != nil {
			fmt.Println("add timer error:", err)
		}

		// 同步云资源状态定时任务 - 每1分钟执行一次
		_, err = global.GVA_Timer.AddTaskByFunc("SyncCloudStatus", "@every 1m", func() {
			fmt.Println("[定时任务] 开始执行云资源状态同步...")
			err := task.SyncAllCloudStatus()
			if err != nil {
				fmt.Println("[定时任务] 云资源状态同步失败:", err)
			} else {
				fmt.Println("[定时任务] 云资源状态同步完成")
			}
		}, "定时同步云资源状态【容器、镜像】", option...)
		if err != nil {
			fmt.Println("add sync cloud status timer error:", err)
		}

		// 其他定时任务定在这里 参考上方使用方法

		//_, err := global.GVA_Timer.AddTaskByFunc("定时任务标识", "corn表达式", func() {
		//	具体执行内容...
		//  ......
		//}, option...)
		//if err != nil {
		//	fmt.Println("add timer error:", err)
		//}
	}()
}
