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
          <el-button type="primary" icon="plus" @click="openCreateDialog" :disabled="!currentNodeId">创建网络</el-button>
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
        <el-table-column label="网络ID" prop="Id" width="120">
          <template #default="scope">
            <el-tooltip :content="scope.row.Id" placement="top">
              <span>{{ scope.row.Id.substring(0, 12) }}</span>
            </el-tooltip>
          </template>
        </el-table-column>
        <el-table-column label="名称" prop="Name" min-width="150" />
        <el-table-column label="驱动" prop="Driver" width="120" />
        <el-table-column label="范围" prop="Scope" width="100" />
        <el-table-column label="创建时间" width="180">
          <template #default="scope">
            {{ formatTime(scope.row.Created) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120" fixed="right">
          <template #default="scope">
            <el-button type="primary" link icon="delete" @click="deleteRow(scope.row)" :disabled="scope.row.Name === 'bridge' || scope.row.Name === 'host' || scope.row.Name === 'none'">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <!-- 创建网络弹窗 -->
    <el-dialog v-model="createDialogVisible" title="创建网络" width="500px">
      <el-form ref="createFormRef" :model="createForm" :rules="createRules" label-width="100px">
        <el-form-item label="网络名称" prop="name">
          <el-input v-model="createForm.name" placeholder="请输入网络名称" />
        </el-form-item>
        <el-form-item label="驱动" prop="driver">
          <el-select v-model="createForm.driver" placeholder="请选择驱动">
            <el-option label="bridge" value="bridge" />
            <el-option label="overlay" value="overlay" />
            <el-option label="macvlan" value="macvlan" />
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
import { getNetworks, createNetwork, removeNetwork } from '@/api/cloud/network'

defineOptions({
  name: 'NetworkManager'
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
  driver: 'bridge'
})
const createRules = {
  name: [{ required: true, message: '请输入网络名称', trigger: 'blur' }],
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
    const res = await getNetworks({ nodeId: currentNodeId.value })
    if (res.code === 0) {
      tableData.value = res.data
    }
  } finally {
    loading.value = false
  }
}

// 格式化时间
const formatTime = (timestamp) => {
  return new Date(timestamp).toLocaleString()
}

// 删除网络
const deleteRow = (row) => {
  ElMessageBox.confirm(`确定要删除网络 ${row.Name} 吗?`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    const res = await removeNetwork({ 
      nodeId: currentNodeId.value, 
      networkId: row.Id
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
  createForm.driver = 'bridge'
  createDialogVisible.value = true
}

// 提交创建
const submitCreate = async () => {
  if (!createFormRef.value) return
  await createFormRef.value.validate(async (valid) => {
    if (valid) {
      creating.value = true
      try {
        const res = await createNetwork({
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
