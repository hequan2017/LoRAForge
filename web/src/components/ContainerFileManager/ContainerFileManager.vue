<template>
  <el-dialog
    v-model="visible"
    title="文件管理"
    width="900px"
    :before-close="handleClose"
    :destroy-on-close="true"
  >
    <div class="file-manager-container">
      <div class="flex items-center justify-between mb-4">
        <div class="flex items-center gap-2 flex-1">
          <el-button type="primary" link icon="ArrowUp" @click="goUp" :disabled="currentPath === '/'">上级</el-button>
          <el-input v-model="currentPath" @keyup.enter="loadFiles" placeholder="输入路径" style="max-width: 400px">
            <template #prepend>/</template>
          </el-input>
          <el-button icon="Refresh" @click="loadFiles" circle />
        </div>
        <div class="flex items-center gap-2">
          <el-upload
            :action="uploadUrl"
            :headers="uploadHeaders"
            :data="uploadData"
            :show-file-list="false"
            :on-success="handleUploadSuccess"
            :on-error="handleUploadError"
            :before-upload="beforeUpload"
          >
            <el-button type="primary" icon="Upload">上传文件</el-button>
          </el-upload>
        </div>
      </div>

      <el-table
        v-loading="loading"
        :data="fileList"
        style="width: 100%"
        height="400"
        @row-dblclick="handleRowDblClick"
      >
        <el-table-column width="40">
          <template #default="scope">
            <el-icon v-if="scope.row.isDir" class="text-yellow-500 text-lg"><Folder /></el-icon>
            <el-icon v-else class="text-gray-500 text-lg"><Document /></el-icon>
          </template>
        </el-table-column>
        <el-table-column label="名称" prop="name" min-width="200" show-overflow-tooltip>
          <template #default="scope">
            <span :class="scope.row.isDir ? 'text-blue-600 cursor-pointer hover:underline' : ''" @click="scope.row.isDir && enterDir(scope.row.name)">
              {{ scope.row.name }}
            </span>
          </template>
        </el-table-column>
        <el-table-column label="大小" prop="size" width="100" />
        <el-table-column label="权限" prop="perm" width="120" />
        <el-table-column label="修改时间" prop="date" width="180" />
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="scope">
            <div class="flex items-center gap-2">
              <el-button v-if="!scope.row.isDir" type="primary" link icon="Download" @click="downloadFile(scope.row)">下载</el-button>
              <el-button type="danger" link icon="Delete" @click="deleteFile(scope.row)">删除</el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>
    </div>
  </el-dialog>
</template>

<script setup>
import { ref, watch, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useUserStore } from '@/pinia'
import service from '@/utils/request'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  instance: {
    type: Object,
    default: () => ({})
  }
})

const emit = defineEmits(['update:modelValue'])
const visible = ref(false)
const userStore = useUserStore()
const loading = ref(false)
const currentPath = ref('/')
const fileList = ref([])

const baseUrl = import.meta.env.VITE_BASE_API || ''
const uploadUrl = computed(() => `${baseUrl}/inst/file/upload`)
const uploadHeaders = computed(() => ({
  'x-token': userStore.token
}))
const uploadData = computed(() => ({
  nodeId: props.instance.NodeID,
  containerId: props.instance.DockerContainer,
  path: currentPath.value === '/' ? '/' : currentPath.value + '/'
}))

watch(() => props.modelValue, (val) => {
  visible.value = val
  if (val && props.instance?.ID) {
    currentPath.value = '/'
    loadFiles()
  }
})

const handleClose = () => {
  emit('update:modelValue', false)
}

const loadFiles = async () => {
  if (!props.instance?.ID) return
  loading.value = true
  try {
    const path = currentPath.value === '/' ? '/' : currentPath.value
    const res = await service({
      url: '/inst/file/list',
      method: 'get',
      params: {
        nodeId: props.instance.NodeID,
        containerId: props.instance.DockerContainer,
        path
      }
    })
    if (res.code === 0) {
      fileList.value = res.data || []
    }
  } catch (e) {
    // console.error(e)
  } finally {
    loading.value = false
  }
}

const goUp = () => {
  if (currentPath.value === '/') return
  const parts = currentPath.value.split('/').filter(Boolean)
  parts.pop()
  currentPath.value = '/' + parts.join('/')
  loadFiles()
}

const enterDir = (name) => {
  if (currentPath.value === '/') {
    currentPath.value = '/' + name
  } else {
    currentPath.value = currentPath.value + '/' + name
  }
  loadFiles()
}

const handleRowDblClick = (row) => {
  if (row.isDir) {
    enterDir(row.name)
  }
}

const downloadFile = (row) => {
  const path = currentPath.value === '/' ? '/' + row.name : currentPath.value + '/' + row.name
  const url = `${baseUrl}/inst/file/download?nodeId=${props.instance.NodeID}&containerId=${props.instance.DockerContainer}&path=${encodeURIComponent(path)}&x-token=${userStore.token}`
  window.open(url, '_blank')
}

const deleteFile = (row) => {
  ElMessageBox.confirm(`确定要删除 ${row.name} 吗?`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    const path = currentPath.value === '/' ? '/' + row.name : currentPath.value + '/' + row.name
    const res = await service({
      url: '/inst/file/delete',
      method: 'post',
      data: {
        nodeId: props.instance.NodeID,
        containerId: props.instance.DockerContainer,
        path
      }
    })
    if (res.code === 0) {
      ElMessage.success('删除成功')
      loadFiles()
    }
  })
}

const beforeUpload = () => {
  loading.value = true
}

const handleUploadSuccess = (res) => {
  loading.value = false
  if (res.code === 0) {
    ElMessage.success('上传成功')
    loadFiles()
  } else {
    ElMessage.error(res.msg || '上传失败')
  }
}

const handleUploadError = () => {
  loading.value = false
  ElMessage.error('上传失败')
}
</script>

<style scoped>
.file-manager-container {
  padding: 10px;
}
</style>
