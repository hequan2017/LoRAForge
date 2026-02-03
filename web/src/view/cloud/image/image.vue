<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" class="demo-form-inline">
        <el-form-item label="计算节点">
          <el-select v-model="currentNodeId" placeholder="请选择节点" @change="handleNodeChange" style="width: 200px">
            <el-option
              v-for="item in nodeList"
              :key="item.ID"
              :label="item.name"
              :value="item.ID"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="refresh" @click="getTableData" :disabled="!currentNodeId">刷新</el-button>
          <el-button type="primary" icon="download" @click="openPullDialog" :disabled="!currentNodeId">拉取镜像</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="gva-table-box">
      <el-table
        :data="tableData"
        style="width: 100%"
        v-loading="loading"
        row-key="Id"
      >
        <el-table-column label="镜像ID" prop="Id" width="180">
          <template #default="scope">
            <el-tooltip :content="scope.row.Id" placement="top">
              <span>{{ scope.row.Id.substring(7, 19) }}</span>
            </el-tooltip>
          </template>
        </el-table-column>
        <el-table-column label="仓库/标签" min-width="250">
          <template #default="scope">
            <div v-if="scope.row.RepoTags && scope.row.RepoTags.length">
              <div v-for="tag in scope.row.RepoTags" :key="tag">{{ tag }}</div>
            </div>
            <div v-else class="text-gray-400">&lt;none&gt;</div>
          </template>
        </el-table-column>
        <el-table-column label="大小" width="120">
          <template #default="scope">
            {{ formatSize(scope.row.Size) }}
          </template>
        </el-table-column>
        <el-table-column label="创建时间" width="180">
          <template #default="scope">
            {{ formatTime(scope.row.Created) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120" fixed="right">
          <template #default="scope">
            <el-button type="primary" link icon="delete" @click="deleteRow(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <!-- 拉取镜像弹窗 -->
    <el-dialog v-model="pullDialogVisible" title="拉取镜像" width="600px" :close-on-click-modal="false">
      <el-form label-width="80px">
        <el-form-item label="镜像名称">
          <el-input v-model="pullImageName" placeholder="例如: nginx:latest 或 mysql:8.0" :disabled="pulling" />
        </el-form-item>
      </el-form>
      
      <!-- 拉取日志 -->
      <div v-if="pullLogs.length > 0" class="bg-black text-white p-4 rounded mt-4 h-64 overflow-y-auto font-mono text-xs">
        <div v-for="(log, index) in pullLogs" :key="index">{{ log }}</div>
      </div>

      <template #footer>
        <span class="dialog-footer">
          <el-button @click="closePullDialog" :disabled="pulling">关闭</el-button>
          <el-button type="primary" @click="startPull" :loading="pulling">开始拉取</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getComputeNodeList } from '@/api/cloud/compute_node'
import { getImages, removeImage } from '@/api/cloud/image'
import { useUserStore } from '@/pinia'

defineOptions({
  name: 'ImageManager'
})

const userStore = useUserStore()
const nodeList = ref([])
const currentNodeId = ref(undefined)
const tableData = ref([])
const loading = ref(false)

// 拉取相关
const pullDialogVisible = ref(false)
const pullImageName = ref('')
const pulling = ref(false)
const pullLogs = ref([])

// 初始化获取节点列表
const init = async () => {
  const res = await getComputeNodeList({ page: 1, pageSize: 100 })
  if (res.code === 0) {
    nodeList.value = res.data.list
    if (nodeList.value.length > 0) {
      currentNodeId.value = nodeList.value[0].ID
      getTableData()
    }
  }
}

const handleNodeChange = () => {
  getTableData()
}

const getTableData = async () => {
  if (!currentNodeId.value) return
  loading.value = true
  try {
    const res = await getImages({ nodeId: currentNodeId.value })
    if (res.code === 0) {
      tableData.value = res.data
    }
  } finally {
    loading.value = false
  }
}

// 格式化大小
const formatSize = (bytes) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

// 格式化时间
const formatTime = (timestamp) => {
  return new Date(timestamp * 1000).toLocaleString()
}

// 删除镜像
const deleteRow = (row) => {
  ElMessageBox.confirm(`确定要删除镜像 ${row.RepoTags ? row.RepoTags[0] : row.Id} 吗?`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    const res = await removeImage({ 
      nodeId: currentNodeId.value, 
      imageId: row.Id,
      force: false 
    })
    if (res.code === 0) {
      ElMessage.success('删除成功')
      getTableData()
    }
  })
}

// 打开拉取弹窗
const openPullDialog = () => {
  pullImageName.value = ''
  pullLogs.value = []
  pullDialogVisible.value = true
}

const closePullDialog = () => {
  if (pulling.value) return
  pullDialogVisible.value = false
}

// 开始拉取 (SSE)
const startPull = () => {
  if (!pullImageName.value) {
    ElMessage.warning('请输入镜像名称')
    return
  }

  pulling.value = true
  pullLogs.value = []
  
  const baseUrl = import.meta.env.VITE_BASE_API || ''
  const url = `${baseUrl}/cloud/image/pull?nodeId=${currentNodeId.value}&imageName=${encodeURIComponent(pullImageName.value)}&x-token=${userStore.token}`
  
  const eventSource = new EventSource(url)

  eventSource.onmessage = (event) => {
    // 尝试解析 JSON (Docker pull output usually is JSON stream)
    try {
      const data = JSON.parse(event.data)
      let msg = ''
      if (data.status) msg += data.status
      if (data.progress) msg += ' ' + data.progress
      if (data.id) msg = `${data.id}: ${msg}`
      
      if (msg) {
          // 简单的日志追加，实际可能需要更复杂的进度条处理
          // 如果是更新同一行（如下载进度），可以优化显示，这里简单追加
          pullLogs.value.push(msg)
          // 自动滚动到底部
          // nextTick...
      }
      
      if (data.error) {
        pullLogs.value.push(`Error: ${data.error}`)
        eventSource.close()
        pulling.value = false
        ElMessage.error(data.error)
      }
    } catch (e) {
      pullLogs.value.push(event.data)
    }
  }

  eventSource.onerror = (err) => {
    console.error('SSE Error:', err)
    // SSE 在连接关闭时也会触发 error，需要区分
    // 这里简单处理：如果正在拉取中遇到错误，认为是异常；如果是完成后的关闭则忽略
    if (eventSource.readyState === EventSource.CLOSED) {
        pullLogs.value.push('连接已关闭')
        pulling.value = false
        getTableData()
    } else {
        // 可能是网络错误
        // pullLogs.value.push('连接发生错误')
        // eventSource.close()
        // pulling.value = false
        // 很多时候 SSE 结束时会触发 error，我们可以在 onmessage 里判断完成标志，或者依赖后端关闭连接
        // 假设后端完成后会关闭连接
        eventSource.close()
        pulling.value = false
        getTableData()
    }
  }
}

onMounted(() => {
  init()
})
</script>
