-- 模型微调模块菜单SQL
-- 请在数据库中执行以下SQL来创建菜单和API权限

-- 1. 首先查询云资源管理父菜单的ID（假设父菜单ID为X，实际使用时需要替换）
-- SELECT id FROM sys_base_menus WHERE name = 'cloud' OR title = '云资源管理' LIMIT 1;

-- 2. 插入API权限记录（sys_apis表）
-- 需要先获取父级分组的最大ID，然后插入新的API

-- 微调任务API列表
INSERT INTO `sys_apis` (`created_at`, `updated_at`, `path`, `description`, `api_group`, `method`) VALUES
(NOW(), NOW(), '/finetune/createFineTuneTask', '创建微调任务', '微调任务', 'POST'),
(NOW(), NOW(), '/finetune/startFineTuneTask', '启动微调任务', '微调任务', 'POST'),
(NOW(), NOW(), '/finetune/stopFineTuneTask', '停止微调任务', '微调任务', 'POST'),
(NOW(), NOW(), '/finetune/cancelFineTuneTask', '取消微调任务', '微调任务', 'POST'),
(NOW(), NOW(), '/finetune/deleteFineTuneTask', '删除微调任务', '微调任务', 'DELETE'),
(NOW(), NOW(), '/finetune/deleteFineTuneTaskByIds', '批量删除微调任务', '微调任务', 'DELETE'),
(NOW(), NOW(), '/finetune/updateFineTuneTask', '更新微调任务', '微调任务', 'PUT'),
(NOW(), NOW(), '/finetune/findFineTuneTask', '查询微调任务', '微调任务', 'GET'),
(NOW(), NOW(), '/finetune/getFineTuneTaskList', '获取微调任务列表', '微调任务', 'GET'),
(NOW(), NOW(), '/finetune/getFineTuneTaskLogs', '获取任务日志', '微调任务', 'GET'),
(NOW(), NOW(), '/finetune/getFineTuneTaskMetrics', '获取训练指标', '微调任务', 'GET'),
(NOW(), NOW(), '/finetune/exportFineTuneTaskModel', '导出微调模型', '微调任务', 'POST'),
(NOW(), NOW(), '/finetune/getFineTuneTaskStatistics', '获取统计信息', '微调任务', 'GET'),
(NOW(), NOW(), '/finetune/createFineTuneTaskSnapshot', '创建训练快照', '微调任务', 'POST'),
(NOW(), NOW(), '/finetune/getFineTuneTaskSnapshots', '获取训练快照', '微调任务', 'GET'),
(NOW(), NOW(), '/finetune/syncFineTuneTaskStatus', '同步任务状态', '微调任务', 'POST'),
(NOW(), NOW(), '/finetune/syncAllFineTuneTaskStatus', '同步所有任务状态', '微调任务', 'POST'),
(NOW(), NOW(), '/finetune/getFineTuneTaskDataSource', '获取数据源', '微调任务', 'GET');

-- 3. 插入菜单记录（sys_base_menus表）
-- 注意：需要替换下面的 @parent_id 为实际的云资源管理父菜单ID

-- 微调任务管理菜单（父菜单下的子菜单）
-- 假设父菜单ID需要手动查询后替换，这里使用占位符
INSERT INTO `sys_base_menus` (
    `created_at`,
    `updated_at`,
    `parent_id`,
    `path`,
    `name`,
    `hidden`,
    `component`,
    `sort`,
    `title`,
    `icon`,
    `keep_alive`,
    `default_menu`,
    `close_tab`
) VALUES (
    NOW(),
    NOW(),
    0,  -- 需要替换为实际的云资源管理父菜单ID
    'finetuneTask',
    'finetuneTask',
    0,
    'view/cloud/finetune_task/finetune_task.vue',
    5,
    '模型微调',
    'cpu',
    0,
    0,
    0
);

-- 4. 获取刚插入的菜单ID，用于创建按钮权限
-- SET @menu_id = LAST_INSERT_ID();

-- 5. 插入按钮权限（sys_base_menu_btns表）
-- 注意：需要替换 @menu_id 为上面插入的菜单ID

INSERT INTO `sys_base_menu_btns` (`created_at`, `updated_at`, `name`, `menu_id`, `description`) VALUES
(NOW(), NOW(), 'create', 0, '新建任务'),
(NOW(), NOW(), 'start', 0, '启动任务'),
(NOW(), NOW(), 'stop', 0, '停止任务'),
(NOW(), NOW(), 'delete', 0, '删除任务'),
(NOW(), NOW(), 'export', 0, '导出模型'),
(NOW(), NOW(), 'sync', 0, '同步状态');

-- ============================================================
-- 使用说明：
-- 1. 执行此SQL前，请先查询云资源管理父菜单的ID：
--    SELECT id FROM sys_base_menus WHERE name LIKE '%cloud%' OR title LIKE '%云资源%' LIMIT 1;
-- 2. 将查询到的父菜单ID替换到 sys_base_menus 表插入语句中的 parent_id 字段
-- 3. 执行完成后，需要在后台管理系统中为对应的角色分配菜单权限
-- ============================================================

-- 完整示例（需要替换实际ID）：
-- 假设云资源管理父菜单ID为 100
-- INSERT INTO `sys_base_menus` VALUES (NULL, NOW(), NOW(), NULL, 100, 'finetuneTask', 'finetuneTask', 0, 'view/cloud/finetune_task/finetune_task.vue', 5, '{"keepAlive":false,"defaultMenu":false,"title":"模型微调","icon":"cpu","closeTab":false}', 0);
