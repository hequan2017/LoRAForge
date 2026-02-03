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
          <el-button type="primary" icon="plus" @click="openCreateDialog" :disabled="!currentNodeId">创建数据卷</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="gva-table-box">
      <el-table
        :data="tableData"
        style="width: 100%"
        v-loading="loading"
        row-key="Name"
      >
        <el-table-column label="名称" prop="Name" min-width="150" />
        <el-table-column label="驱动" prop="Driver" width="120" />
        <el-table-column label="挂载点" prop="Mountpoint" min-width="200">
          <template #default="scope">
             <span class="text-xs text-gray-500">{{ scope.row.Mountpoint }}</span>
          </template>
        </el-table-column>
        <el-table-column label="创建时间" width="180">
          <template #default="scope">
            {{ formatTime(scope.row.CreatedAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120" fixed="right">
          <template #default="scope">
            <el-button type="primary" link icon="delete" @click="deleteRow(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <!-- 创建卷弹窗 -->
    <el-dialog v-model="createDialogVisible" title="创建数据卷" width="500px">
      <el-form ref="createFormRef" :model="createForm" :rules="createRules" label-width="100px">
        <el-form-item label="卷名称" prop="name">
          <el-input v-model="createForm.name" placeholder="请输入卷名称" />
        </el-form-item>
        <el-form-item label="驱动" prop="driver">
          <el-select v-model="createForm.driver" placeholder="请选择驱动">
            <el-option label="local" value="local" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="createDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitCreate" :loading="creating">创建</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getComputeNodeList } from '@/api/cloud/compute_node'
import { getVolumes, createVolume, removeVolume } from '@/api/cloud/volume'

defineOptions({
  name: 'VolumeManager'
})

const nodeList = ref([])
const currentNodeId = ref(undefined)
const tableData = ref([])
const loading = ref(false)

// 创建相关
const createDialogVisible = ref(false)
const creating = ref(false)
const createFormRef = ref(null)
const createForm = reactive({
  name: '',
  driver: 'local'
})
const createRules = {
  name: [{ required: true, message: '请输入卷名称', trigger: 'blur' }],
  driver: [{ required: true, message: '请选择驱动', trigger: 'change' }]
}

// 初始化
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
    const res = await getVolumes({ nodeId: currentNodeId.value })
    if (res.code === 0) {
      tableData.value = res.data
    }
  } finally {
    loading.value = false
  }
}

// 格式化时间
const formatTime = (timestamp) => {
  if (!timestamp) return '-'
  return new Date(timestamp).toLocaleString()
}

// 删除卷
const deleteRow = (row) => {
  ElMessageBox.confirm(`确定要删除数据卷 ${row.Name} 吗? (此操作不可逆)`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    const res = await removeVolume({ 
      nodeId: currentNodeId.value, 
      volumeId: row.Name,
      force: false
    })
    if (res.code === 0) {
      ElMessage.success('删除成功')
      getTableData()
    }
  })
}

// 打开创建弹窗
const openCreateDialog = () => {
  createForm.name = ''
  createForm.driver = 'local'
  createDialogVisible.value = true
}

// 提交创建
const submitCreate = async () => {
  if (!createFormRef.value) return
  await createFormRef.value.validate(async (valid) => {
    if (valid) {
      creating.value = true
      try {
        const res = await createVolume({
          nodeId: currentNodeId.value,
          name: createForm.name,
          driver: createForm.driver
        })
        if (res.code === 0) {
          ElMessage.success('创建成功')
          createDialogVisible.value = false
          getTableData()
        }
      } finally {
        creating.value = false
      }
    }
  })
}

onMounted(() => {
  init()
})
</script>
