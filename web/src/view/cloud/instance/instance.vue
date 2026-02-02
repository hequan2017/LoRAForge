
<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" @keyup.enter="onSubmit">
      <el-form-item label="创建日期" prop="createdAtRange">
      <template #label>
        <span>
          创建日期
          <el-tooltip content="搜索范围是开始日期（包含）至结束日期（不包含）">
            <el-icon><QuestionFilled /></el-icon>
          </el-tooltip>
        </span>
      </template>

      <el-date-picker
            v-model="searchInfo.createdAtRange"
            class="!w-380px"
            type="datetimerange"
            range-separator="至"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
          />
       </el-form-item>
        <el-form-item label="实例名称" prop="instanceName">
            <el-input v-model="searchInfo.instanceName" placeholder="搜索实例名称" />
       </el-form-item>

       <template v-if="showAllQuery">
       <!-- 将需要控制显示状态的查询条件添加到此范围内 -->
       <el-form-item label="镜像" prop="mirrorId">
            <el-select v-model="searchInfo.mirrorId" placeholder="请选择镜像" filterable clearable>
                <el-option v-for="(item,key) in dataSource.mirrorId" :key="key" :label="item.label" :value="item.value" />
            </el-select>
       </el-form-item>
       <el-form-item label="模版" prop="templateId">
            <el-select v-model="searchInfo.templateId" placeholder="请选择模版" filterable clearable>
                <el-option v-for="(item,key) in dataSource.templateId" :key="key" :label="item.label" :value="item.value" />
            </el-select>
       </el-form-item>
       <el-form-item label="节点" prop="nodeId">
            <el-select v-model="searchInfo.nodeId" placeholder="请选择节点" filterable clearable>
                <el-option v-for="(item,key) in dataSource.nodeId" :key="key" :label="item.label" :value="item.value" />
            </el-select>
       </el-form-item>
       <el-form-item label="容器状态" prop="containerStatus">
            <el-select v-model="searchInfo.containerStatus" placeholder="请选择状态" clearable>
                <el-option label="运行中" value="Running" />
                <el-option label="已关闭" value="Closed" />
                <el-option label="已停止" value="Stopped" />
            </el-select>
       </el-form-item>
       </template>

        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
          <el-button link type="primary" icon="arrow-down" @click="showAllQuery=true" v-if="!showAllQuery">展开</el-button>
          <el-button link type="primary" icon="arrow-up" @click="showAllQuery=false" v-else>收起</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
        <div class="gva-btn-list">
            <el-button  type="primary" icon="plus" @click="openDialog()">新增实例</el-button>
            <el-button  icon="delete" style="margin-left: 10px;" :disabled="!multipleSelection.length" @click="onDelete">批量删除</el-button>

        </div>
        <el-table
        ref="multipleTable"
        style="width: 100%"
        tooltip-effect="dark"
        :data="tableData"
        row-key="ID"
        @selection-change="handleSelectionChange"
        >
        <el-table-column type="selection" width="55" />

        <el-table-column sortable align="left" label="日期" prop="CreatedAt" width="180">
            <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>

            <el-table-column align="left" label="实例名称" prop="instanceName" width="150" />

            <el-table-column align="left" label="镜像" prop="mirrorId" width="140">
    <template #default="scope">
        <span>{{ filterDataSource(dataSource.mirrorId,scope.row.mirrorId) }}</span>
    </template>
</el-table-column>
            <el-table-column align="left" label="节点" prop="nodeId" width="120">
    <template #default="scope">
        <span>{{ filterDataSource(dataSource.nodeId,scope.row.nodeId) }}</span>
    </template>
</el-table-column>

            <el-table-column align="left" label="配置" width="180">
                <template #default="scope">
                    <el-tag size="small" v-if="scope.row.cpu" type="info">CPU: {{ scope.row.cpu }}核</el-tag>
                    <el-tag size="small" v-if="scope.row.memory" type="success" style="margin-left: 4px">{{ formatMemory(scope.row.memory) }}</el-tag>
                    <el-tag size="small" v-if="scope.row.gpuCount" type="warning" style="margin-left: 4px">GPU: {{ scope.row.gpuCount }}</el-tag>
                </template>
            </el-table-column>

            <el-table-column align="left" label="容器ID" prop="dockerContainer" width="100">
                <template #default="scope">
                    <el-tooltip v-if="scope.row.dockerContainer" :content="scope.row.dockerContainer" placement="top">
                        <span class="text-xs text-gray-500">{{ scope.row.dockerContainer?.substring(0, 8) }}...</span>
                    </el-tooltip>
                </template>
            </el-table-column>

            <el-table-column align="left" label="状态" prop="containerStatus" width="100">
                <template #default="scope">
                    <el-tag v-if="scope.row.containerStatus === 'Running'" type="success">运行中</el-tag>
                    <el-tag v-else-if="scope.row.containerStatus === 'Closed'" type="info">已关闭</el-tag>
                    <el-tag v-else-if="scope.row.containerStatus === 'Stopped'" type="danger">已停止</el-tag>
                    <el-tag v-else type="info">{{ scope.row.containerStatus || '未知' }}</el-tag>
                </template>
            </el-table-column>

        <el-table-column align="left" label="操作" fixed="right" :min-width="appStore.operateMinWith">
            <template #default="scope">
            <el-button  type="success" link class="table-button" @click="openTerminal(scope.row)"><el-icon style="margin-right: 5px"><Monitor /></el-icon>终端</el-button>
            <el-button  type="info" link class="table-button" @click="openLog(scope.row)"><el-icon style="margin-right: 5px"><Document /></el-icon>日志</el-button>
            <el-button  type="primary" link class="table-button" @click="getDetails(scope.row)"><el-icon style="margin-right: 5px"><InfoFilled /></el-icon>查看</el-button>
            <el-button  type="primary" link icon="edit" class="table-button" @click="updateInstanceFunc(scope.row)">编辑</el-button>
            <el-button   type="primary" link icon="delete" @click="deleteRow(scope.row)">删除</el-button>
            <el-button type="success" link icon="VideoPlay" @click="startInstanceFunc(scope.row)">启动</el-button>
            <el-button type="warning" link icon="SwitchButton" @click="closeInstanceFunc(scope.row)">关闭</el-button>
            <el-button type="success" link icon="Refresh" @click="restartInstanceFunc(scope.row)">重启</el-button>
            </template>
        </el-table-column>
        </el-table>
        <div class="gva-pagination">
            <el-pagination
            layout="total, sizes, prev, pager, next, jumper"
            :current-page="page"
            :page-size="pageSize"
            :page-sizes="[10, 30, 50, 100]"
            :total="total"
            @current-change="handleCurrentChange"
            @size-change="handleSizeChange"
            />
        </div>
    </div>

    <!-- 创建/编辑抽屉 -->
    <el-drawer destroy-on-close :size="650" v-model="dialogFormVisible" :show-close="false" :before-close="closeDialog">
       <template #header>
              <div class="flex justify-between items-center">
                <span class="text-lg font-semibold">{{type==='create'?'创建新实例':'编辑实例'}}</span>
                <div>
                  <el-button :loading="btnLoading" type="primary" @click="enterDialog">{{type==='create'?'创建':'保存'}}</el-button>
                  <el-button @click="closeDialog">取消</el-button>
                </div>
              </div>
            </template>

          <el-form :model="formData" label-position="top" ref="elFormRef" :rules="rule" label-width="100px">
            <el-tabs v-model="activeTab" class="w-full">
                <!-- 基础配置 -->
                <el-tab-pane label="基础配置" name="basic">
                    <el-row :gutter="16">
                        <el-col :span="12">
                            <el-form-item label="实例名称" prop="instanceName">
                                <el-input v-model="formData.instanceName" placeholder="请输入实例名称" clearable />
                            </el-form-item>
                        </el-col>
                        <el-col :span="12">
                            <el-form-item label="节点" prop="nodeId">
                                <el-select v-model="formData.nodeId" placeholder="请选择节点" filterable style="width:100%">
                                    <el-option v-for="(item,key) in dataSource.nodeId" :key="key" :label="item.label" :value="item.value" />
                                </el-select>
                            </el-form-item>
                        </el-col>
                    </el-row>

                    <el-row :gutter="16">
                        <el-col :span="12">
                            <el-form-item label="镜像" prop="mirrorId">
                                <el-select v-model="formData.mirrorId" placeholder="请选择镜像" filterable style="width:100%">
                                    <el-option v-for="(item,key) in dataSource.mirrorId" :key="key" :label="item.label" :value="item.value" />
                                </el-select>
                            </el-form-item>
                        </el-col>
                        <el-col :span="12">
                            <el-form-item label="规格模版">
                                <el-select v-model="formData.templateId" placeholder="选择模版自动填充配置" filterable clearable style="width:100%" @change="onTemplateChange">
                                    <el-option v-for="(item,key) in dataSource.templateId" :key="key" :label="item.label" :value="item.value" />
                                </el-select>
                            </el-form-item>
                        </el-col>
                    </el-row>

                    <el-form-item label="启动命令">
                        <el-input v-model="formData.command" placeholder="例如: /bin/bash 或 python app.py" clearable>
                            <template #append>
                                <el-tooltip content="容器启动后执行的命令，留空则使用镜像默认命令" placement="top">
                                    <el-icon><QuestionFilled /></el-icon>
                                </el-tooltip>
                            </template>
                        </el-input>
                    </el-form-item>

                    <el-form-item label="备注">
                        <el-input v-model="formData.remark" type="textarea" :rows="2" placeholder="请输入备注信息" />
                    </el-form-item>
                </el-tab-pane>

                <!-- 资源限制 -->
                <el-tab-pane label="资源限制" name="resource">
                    <el-row :gutter="16">
                        <el-col :span="8">
                            <el-form-item label="CPU核数">
                                <el-input-number v-model="formData.cpu" :min="0.1" :max="64" :step="0.1" :precision="1" controls-position="right" style="width:100%" />
                                <div class="text-xs text-gray-400 mt-1">0.1 = 10%，1 = 1核</div>
                            </el-form-item>
                        </el-col>
                        <el-col :span="8">
                            <el-form-item label="内存(MB)">
                                <el-input-number v-model="formData.memory" :min="128" :max="102400" :step="128" controls-position="right" style="width:100%" />
                                <div class="text-xs text-gray-400 mt-1">建议最小 512MB</div>
                            </el-form-item>
                        </el-col>
                        <el-col :span="8">
                            <el-form-item label="GPU数量">
                                <el-input-number v-model="formData.gpuCount" :min="0" :max="8" :step="1" controls-position="right" style="width:100%" />
                                <div class="text-xs text-gray-400 mt-1">0表示不使用GPU</div>
                            </el-form-item>
                        </el-col>
                    </el-row>

                    <el-alert type="info" :closable="false" show-icon class="mb-4">
                        <template #title>
                            资源限制说明
                        </template>
                        <div class="text-xs">
                            • CPU：限制容器可使用的CPU核心数，0.5表示50%的一个核心<br>
                            • 内存：限制容器可使用的最大内存，超出会被OOM杀死<br>
                            • GPU：需要节点支持GPU且有nvidia-docker支持
                        </div>
                    </el-alert>
                </el-tab-pane>

                <!-- 网络配置 -->
                <el-tab-pane label="网络配置" name="network">
                    <el-form-item label="端口映射">
                        <el-input v-model="formData.portMapping" type="textarea" :rows="5" placeholder="每行一个映射，格式：主机端口:容器端口&#10;例如：&#10;8080:80&#10;3306:3306" />
                        <template #label>
                            <span>端口映射 </span>
                            <el-tooltip content="将容器内部端口映射到主机端口，每行一个，格式为 主机端口:容器端口" placement="top">
                                <el-icon class="ml-1"><QuestionFilled /></el-icon>
                            </el-tooltip>
                        </template>
                    </el-form-item>
                </el-tab-pane>

                <!-- 存储与环境 -->
                <el-tab-pane label="存储与环境" name="storage">
                    <el-form-item label="挂载目录">
                        <el-input v-model="formData.volumeMounts" type="textarea" :rows="4" placeholder="每行一个挂载，格式：主机路径:容器路径[:读写权限]&#10;例如：&#10;/data/host:/data/container&#10;/tmp/logs:/var/log/app:ro" />
                        <template #label>
                            <span>挂载目录 </span>
                            <el-tooltip content="将主机目录挂载到容器内，每行一个，格式为 主机路径:容器路径，可选:ro表示只读" placement="top">
                                <el-icon class="ml-1"><QuestionFilled /></el-icon>
                            </el-tooltip>
                        </template>
                    </el-form-item>

                    <el-form-item label="环境变量">
                        <el-input v-model="formData.envVars" type="textarea" :rows="5" placeholder="每行一个环境变量，格式：KEY=VALUE&#10;例如：&#10;NODE_ENV=production&#10;API_KEY=your_api_key" />
                        <template #label>
                            <span>环境变量 </span>
                            <el-tooltip content="设置容器的环境变量，每行一个，格式为 KEY=VALUE" placement="top">
                                <el-icon class="ml-1"><QuestionFilled /></el-icon>
                            </el-tooltip>
                        </template>
                    </el-form-item>
                </el-tab-pane>
            </el-tabs>
          </el-form>
    </el-drawer>

    <!-- 详情抽屉 -->
    <el-drawer destroy-on-close :size="600" v-model="detailShow" :show-close="true" :before-close="closeDetailShow" title="实例详情">
            <el-descriptions :column="2" border>
                <el-descriptions-item label="实例名称" :span="2">
                    <span class="font-semibold">{{ detailForm.instanceName || '-' }}</span>
                </el-descriptions-item>
                <el-descriptions-item label="镜像">
                    {{ filterDataSource(dataSource.mirrorId,detailForm.mirrorId) }}
                </el-descriptions-item>
                <el-descriptions-item label="节点">
                    {{ filterDataSource(dataSource.nodeId,detailForm.nodeId) }}
                </el-descriptions-item>
                <el-descriptions-item label="规格模版">
                    {{ filterDataSource(dataSource.templateId,detailForm.templateId) }}
                </el-descriptions-item>
                <el-descriptions-item label="容器状态">
                    <el-tag v-if="detailForm.containerStatus === 'Running'" type="success" size="small">运行中</el-tag>
                    <el-tag v-else-if="detailForm.containerStatus === 'Closed'" type="info" size="small">已关闭</el-tag>
                    <el-tag v-else-if="detailForm.containerStatus === 'Stopped'" type="danger" size="small">已停止</el-tag>
                    <el-tag v-else type="info" size="small">{{ detailForm.containerStatus || '未知' }}</el-tag>
                </el-descriptions-item>
                <el-descriptions-item label="Docker容器ID" :span="2">
                    <span class="font-mono text-xs">{{ detailForm.dockerContainer || '-' }}</span>
                </el-descriptions-item>

                <el-descriptions-item label="资源配置" :span="2">
                    <div class="flex items-center gap-2">
                        <el-tag v-if="detailForm.cpu" type="info" size="small">CPU: {{ detailForm.cpu }}核</el-tag>
                        <el-tag v-if="detailForm.memory" type="success" size="small">内存: {{ formatMemory(detailForm.memory) }}</el-tag>
                        <el-tag v-if="detailForm.gpuCount" type="warning" size="small">GPU: {{ detailForm.gpuCount }}个</el-tag>
                        <span v-if="!detailForm.cpu && !detailForm.memory && !detailForm.gpuCount" class="text-gray-400">未设置</span>
                    </div>
                </el-descriptions-item>

                <el-descriptions-item label="启动命令" :span="2">
                    <span class="font-mono text-xs">{{ detailForm.command || '使用镜像默认命令' }}</span>
                </el-descriptions-item>

                <el-descriptions-item label="端口映射" :span="2">
                    <pre class="text-xs bg-gray-50 p-2 rounded max-h-32 overflow-auto">{{ detailForm.portMapping || '无' }}</pre>
                </el-descriptions-item>

                <el-descriptions-item label="挂载目录" :span="2">
                    <pre class="text-xs bg-gray-50 p-2 rounded max-h-32 overflow-auto">{{ detailForm.volumeMounts || '无' }}</pre>
                </el-descriptions-item>

                <el-descriptions-item label="环境变量" :span="2">
                    <pre class="text-xs bg-gray-50 p-2 rounded max-h-32 overflow-auto">{{ detailForm.envVars || '无' }}</pre>
                </el-descriptions-item>

                <el-descriptions-item label="备注" :span="2">
                    {{ detailForm.remark || '-' }}
                </el-descriptions-item>

                <el-descriptions-item label="创建时间" :span="2">
                    {{ formatDate(detailForm.CreatedAt) }}
                </el-descriptions-item>
            </el-descriptions>
    </el-drawer>

    <!-- WebSSH 终端弹窗 -->
    <el-dialog
        v-model="terminalVisible"
        :fullscreen="false"
        :show-close="true"
        @close="closeTerminal"
        class="terminal-dialog"
        width="90%"
        destroy-on-close>
        <template #header>
            <div class="flex justify-between items-center w-full">
                <span class="text-lg font-semibold">
                    <el-icon class="mr-2"><Monitor /></el-icon>
                    Web终端 - {{ currentInstance?.instanceName || '容器' }}
                </span>
                <el-tag v-if="currentInstance" size="small" type="success">
                    {{ currentInstance.containerStatus === 'Running' ? '运行中' : currentInstance.containerStatus }}
                </el-tag>
            </div>
        </template>
        <div class="terminal-wrapper">
            <WebTerminal
                v-if="terminalVisible && terminalUrl"
                :url="terminalUrl"
                :title="currentInstance?.instanceName || 'Terminal'"
                @connected="onTerminalConnected"
                @disconnected="onTerminalDisconnected"
                @error="onTerminalError" />
            <div v-else-if="!currentInstance || currentInstance.containerStatus !== 'Running'" class="flex items-center justify-center h-full">
                <el-empty description="容器状态不是运行中，无法连接终端" />
            </div>
        </div>
    </el-dialog>

    <!-- 容器日志弹窗 -->
    <el-dialog
        v-model="logVisible"
        :fullscreen="false"
        :show-close="true"
        @close="closeLog"
        class="terminal-dialog"
        width="90%"
        destroy-on-close>
        <template #header>
            <div class="flex justify-between items-center w-full">
                <span class="text-lg font-semibold">
                    <el-icon class="mr-2"><Document /></el-icon>
                    容器日志 - {{ currentInstance?.instanceName || '容器' }}
                </span>
                <el-tag v-if="currentInstance" size="small" type="success">
                    {{ currentInstance.containerStatus === 'Running' ? '运行中' : currentInstance.containerStatus }}
                </el-tag>
            </div>
        </template>
        <div class="terminal-wrapper">
            <WebLog
                v-if="logVisible && logUrl"
                :url="logUrl"
                :title="currentInstance?.instanceName || 'Logs'" />
        </div>
    </el-dialog>

  </div>
</template>

<script setup>
import {
    getInstanceDataSource,
  createInstance,
  deleteInstance,
  deleteInstanceByIds,
  updateInstance,
  findInstance,
  getInstanceList,
  closeInstance,
  restartInstance,
  startInstance,
  syncInstances
} from '@/api/cloud/instance'
import { getComputeNodeList } from '@/api/cloud/compute_node'
import WebTerminal from '@/components/WebTerminal/WebTerminal.vue'

// 全量引入格式化工具 请按需保留
import { getDictFunc, formatDate, formatBoolean, filterDict ,filterDataSource, returnArrImg, onDownloadFile } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive, computed } from 'vue'
import { useAppStore } from "@/pinia"
import { useUserStore } from "@/pinia"




defineOptions({
    name: 'Instance'
})

// 提交按钮loading
const btnLoading = ref(false)
const appStore = useAppStore()

// 当前激活的标签页
const activeTab = ref('basic')

// 控制更多查询条件显示/隐藏状态
const showAllQuery = ref(false)

// 算力节点相关
const nodeList = ref([])
const currentNodeId = ref(null)
const syncLoading = ref(false)

const getNodeList = async () => {
  const res = await getComputeNodeList({ page: 1, pageSize: 100 })
  if (res.code === 0) {
    nodeList.value = res.data.list
  }
}

const handleNodeClick = async (node) => {
  // 如果点击已选中的节点，仅刷新列表
  if (currentNodeId.value === node.ID) {
    // 也可以选择重新同步
    // return
  }
  
  currentNodeId.value = node.ID
  // 自动填充查询条件
  searchInfo.value.nodeId = node.ID
  // 刷新列表
  onSubmit()

  // 同步容器信息
  syncLoading.value = true
  try {
    const res = await syncInstances({ nodeId: node.ID })
    if (res.code === 0) {
      ElMessage.success('同步成功，列表已更新')
      onSubmit() // 同步完成后再次刷新以显示最新数据
    }
  } catch (error) {
    console.error('Sync failed:', error)
  } finally {
    syncLoading.value = false
  }
}

// 初始化时获取节点列表
getNodeList()

// ============ WebSSH 终端相关 ============
const terminalVisible = ref(false)
const currentInstance = ref(null)
const userStore = useUserStore()

// 计算 WebSocket URL
const terminalUrl = computed(() => {
  if (!currentInstance.value) return ''

  // 获取当前主机的协议和主机名
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  const host = window.location.host
  const baseUrl = import.meta.env.VITE_BASE_API || ''

  // 构建 WebSocket URL
  return `${protocol}//${host}${baseUrl}/inst/webssh?id=${currentInstance.value.ID}&x-token=${encodeURIComponent(userStore.token)}`
})

// 打开终端
const openTerminal = (row) => {
  if (row.containerStatus !== 'Running') {
    ElMessage.warning('容器状态不是运行中，无法连接终端')
    return
  }
  currentInstance.value = row
  terminalVisible.value = true
}

// 关闭终端
const closeTerminal = () => {
  terminalVisible.value = false
  currentInstance.value = null
}

// 终端连接成功回调
const onTerminalConnected = () => {
  ElMessage.success('终端连接成功')
}

// 终端断开连接回调
const onTerminalDisconnected = () => {
  console.log('终端已断开连接')
}

// 终端错误回调
const onTerminalError = (error) => {
  ElMessage.error('终端连接失败')
  console.error('终端错误:', error)
}
// ============ WebSSH 终端结束 ============

// 自动化生成的字典（可能为空）以及字段
const formData = ref({
            mirrorId: undefined,
            templateId: undefined,
            nodeId: undefined,
            instanceName: '',
            cpu: 1,
            memory: 1024,
            gpuCount: 0,
            portMapping: '',
            volumeMounts: '',
            envVars: '',
            command: '',
            remark: '',
        })
  const dataSource = ref([])
  const getDataSourceFunc = async()=>{
    const res = await getInstanceDataSource()
    if (res.code === 0) {
      dataSource.value = res.data
    }
  }
  getDataSourceFunc()


  // 模版变更时自动填充资源配置
  const onTemplateChange = async (templateId) => {
      if (!templateId) return
      // 这里可以根据模版ID获取模版详情并自动填充
      // 暂时简化处理，实际可以调用API获取模版详情
  }



// 验证规则
const rule = reactive({
               mirrorId : [{
                   required: true,
                   message: '请选择镜像',
                   trigger: ['change','blur'],
               },
              ],
               nodeId : [{
                   required: true,
                   message: '请选择节点',
                   trigger: ['change','blur'],
               },
              ],
               instanceName : [{
                   required: true,
                   message: '请输入实例名称',
                   trigger: ['blur'],
               },
               {
                   pattern: /^[a-zA-Z0-9\u4e00-\u9fa5_-]+$/,
                   message: '实例名称只能包含字母、数字、中文、下划线和连字符',
                   trigger: ['blur', 'change'],
               },
              ],
})

const elFormRef = ref()
const elSearchFormRef = ref()

// =========== 表格控制部分 ===========
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})
// 重置
const onReset = () => {
  searchInfo.value = {}
  getTableData()
}

// 搜索
const onSubmit = () => {
  elSearchFormRef.value?.validate(async(valid) => {
    if (!valid) return
    page.value = 1
    getTableData()
  })
}

// 分页
const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

// 修改页面容量
const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

// 查询
const getTableData = async() => {
  const table = await getInstanceList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

getTableData()

// ============== 表格控制部分结束 ===============

// 获取需要的字典 可能为空 按需保留
const setOptions = async () =>{
}

// 获取需要的字典 可能为空 按需保留
setOptions()


// 多选数据
const multipleSelection = ref([])
// 多选
const handleSelectionChange = (val) => {
    multipleSelection.value = val
}

// 删除行
const deleteRow = (row) => {
    ElMessageBox.confirm('确定要删除吗?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
    }).then(() => {
            deleteInstanceFunc(row)
        })
    }

// 多选删除
const onDelete = async() => {
  ElMessageBox.confirm('确定要删除吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async() => {
      const IDs = []
      if (multipleSelection.value.length === 0) {
        ElMessage({
          type: 'warning',
          message: '请选择要删除的数据'
        })
        return
      }
      multipleSelection.value &&
        multipleSelection.value.map(item => {
          IDs.push(item.ID)
        })
      const res = await deleteInstanceByIds({ IDs })
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: '删除成功'
        })
        if (tableData.value.length === IDs.length && page.value > 1) {
          page.value--
        }
        getTableData()
      }
    })
    }

// 行为控制标记（弹窗内部需要增还是改）
const type = ref('')

// 更新行
const updateInstanceFunc = async(row) => {
    const res = await findInstance({ ID: row.ID })
    type.value = 'update'
    if (res.code === 0) {
        formData.value = res.data
        dialogFormVisible.value = true
    }
}


// 删除行
const deleteInstanceFunc = async (row) => {
    const res = await deleteInstance({ ID: row.ID })
    if (res.code === 0) {
        ElMessage({
                type: 'success',
                message: '删除成功'
            })
            if (tableData.value.length === 1 && page.value > 1) {
            page.value--
        }
        getTableData()
    }
}

// 弹窗控制标记
const dialogFormVisible = ref(false)

// 打开弹窗
const openDialog = () => {
    type.value = 'create'
    activeTab.value = 'basic'
    dialogFormVisible.value = true
}

// 关闭弹窗
const closeDialog = () => {
    dialogFormVisible.value = false
    formData.value = {
        mirrorId: undefined,
        templateId: undefined,
        nodeId: undefined,
        instanceName: '',
        cpu: 1,
        memory: 1024,
        gpuCount: 0,
        portMapping: '',
        volumeMounts: '',
        envVars: '',
        command: '',
        remark: '',
        }
}
// 弹窗确定
const enterDialog = async () => {
     btnLoading.value = true
     elFormRef.value?.validate( async (valid) => {
             if (!valid) return btnLoading.value = false
              let res
              switch (type.value) {
                case 'create':
                  res = await createInstance(formData.value)
                  break
                case 'update':
                  res = await updateInstance(formData.value)
                  break
                default:
                  res = await createInstance(formData.value)
                  break
              }
              btnLoading.value = false
              if (res.code === 0) {
                ElMessage({
                  type: 'success',
                  message: type.value === 'create' ? '创建成功' : '保存成功'
                })
                closeDialog()
                getTableData()
              }
      })
}

const detailForm = ref({})

// 查看详情控制标记
const detailShow = ref(false)


// 打开详情弹窗
const openDetailShow = () => {
  detailShow.value = true
}


// 打开详情
const getDetails = async (row) => {
  // 打开弹窗
  const res = await findInstance({ ID: row.ID })
  if (res.code === 0) {
    detailForm.value = res.data
    openDetailShow()
  }
}


// 关闭详情弹窗
const closeDetailShow = () => {
  detailShow.value = false
  detailForm.value = {}
}

const startInstanceFunc = (row) => {
    ElMessageBox.confirm('确定要启动此实例吗?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'success'
    }).then(async () => {
        const res = await startInstance({ ID: row.ID })
        if (res.code === 0) {
            ElMessage({
                type: 'success',
                message: '启动成功'
            })
            getTableData()
        }
    })
}

const closeInstanceFunc = (row) => {
    ElMessageBox.confirm('确定要关闭此实例吗?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
    }).then(async () => {
        const res = await closeInstance({ ID: row.ID })
        if (res.code === 0) {
            ElMessage({
                type: 'success',
                message: '关闭成功'
            })
            getTableData()
        }
    })
}

const restartInstanceFunc = (row) => {
    ElMessageBox.confirm('确定要重启此实例吗?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
    }).then(async () => {
        const res = await restartInstance({ ID: row.ID })
        if (res.code === 0) {
            ElMessage({
                type: 'success',
                message: '重启成功'
            })
            getTableData()
        }
    })
}

// 格式化内存显示
const formatMemory = (mb) => {
    if (!mb) return '-'
    if (mb < 1024) return `${mb}MB`
    return `${(mb / 1024).toFixed(1)}GB`
}



</script>

<style scoped>
:deep(.el-drawer__header) {
    margin-bottom: 0;
    padding: 16px 20px;
    border-bottom: 1px solid #e5e7eb;
}
:deep(.el-drawer__body) {
    padding: 20px;
}
:deep(.el-tabs__content) {
    padding-top: 16px;
}
pre {
    white-space: pre-wrap;
    word-break: break-all;
}

/* 终端弹窗样式 */
:deep(.terminal-dialog .el-dialog__body) {
    padding: 0;
    background-color: #1e1e1e;
    min-height: 500px;
    max-height: 70vh;
}

:deep(.terminal-dialog .el-dialog__header) {
    padding: 12px 16px;
    background-color: #252526;
    border-bottom: 1px solid #3e3e3e;
}

.terminal-wrapper {
    height: 600px;
    display: flex;
    flex-direction: column;
}
</style>
