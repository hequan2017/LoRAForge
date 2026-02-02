
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
        <el-form-item label="名字" prop="name">
            <el-input v-model="searchInfo.name" placeholder="搜索名字" />
       </el-form-item>
       
       <template v-if="showAllQuery">
       <!-- 将需要控制显示状态的查询条件添加到此范围内 -->
       <el-form-item label="区域" prop="region">
            <el-input v-model="searchInfo.region" placeholder="搜索区域" />
       </el-form-item>
       <el-form-item label="公网IP" prop="publicIp">
            <el-input v-model="searchInfo.publicIp" placeholder="搜索公网IP" />
       </el-form-item>
       <el-form-item label="显卡名称" prop="gpuName">
            <el-input v-model="searchInfo.gpuName" placeholder="搜索显卡名称" />
       </el-form-item>
       <el-form-item label="是否上架" prop="isListed">
            <el-select v-model="searchInfo.isListed" placeholder="请选择" clearable>
               <el-option label="是" :value="true" />
               <el-option label="否" :value="false" />
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
            <el-button  type="primary" icon="plus" @click="openDialog()">新增</el-button>
            <el-button  icon="delete" style="margin-left: 10px;" :disabled="!multipleSelection.length" @click="onDelete">删除</el-button>
            
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
        
            <el-table-column align="left" label="名字" prop="name" width="120" />

            <el-table-column align="left" label="区域" prop="region" width="120" />

            <el-table-column align="left" label="CPU" prop="cpu" width="120" />

            <el-table-column align="left" label="内存" prop="memory" width="120" />

            <el-table-column align="left" label="系统盘容量" prop="systemDiskCapacity" width="120" />

            <el-table-column align="left" label="数据盘容量" prop="dataDiskCapacity" width="120" />

            <el-table-column align="left" label="公网IP" prop="publicIp" width="120" />

            <el-table-column align="left" label="内网IP" prop="privateIp" width="120" />

            <el-table-column align="left" label="SSH端口" prop="sshPort" width="120" />

            <el-table-column align="left" label="用户名" prop="username" width="120" />

            <el-table-column align="left" label="密码" prop="password" width="120" />

            <el-table-column align="left" label="显卡名称" prop="gpuName" width="120" />

            <el-table-column align="left" label="显卡数量" prop="gpuCount" width="120" />

            <el-table-column align="left" label="Docker连接地址" prop="dockerConnectAddress" width="120" />

            <el-table-column align="left" label="使用TLS" prop="useTls" width="120">
    <template #default="scope">{{ formatBoolean(scope.row.useTls) }}</template>
</el-table-column>
            <el-table-column align="left" label="是否上架" prop="isListed" width="120">
    <template #default="scope">{{ formatBoolean(scope.row.isListed) }}</template>
</el-table-column>
            <el-table-column align="left" label="备注" prop="remark" width="120" />

        <el-table-column align="left" label="操作" fixed="right" :min-width="appStore.operateMinWith">
            <template #default="scope">
            <el-button  type="primary" link class="table-button" @click="getDetails(scope.row)"><el-icon style="margin-right: 5px"><InfoFilled /></el-icon>查看</el-button>
            <el-button  type="primary" link icon="edit" class="table-button" @click="updateComputeNodeFunc(scope.row)">编辑</el-button>
            <el-button   type="primary" link icon="delete" @click="deleteRow(scope.row)">删除</el-button>
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
    <el-drawer destroy-on-close :size="appStore.drawerSize" v-model="dialogFormVisible" :show-close="false" :before-close="closeDialog">
       <template #header>
              <div class="flex justify-between items-center">
                <span class="text-lg">{{type==='create'?'新增':'编辑'}}</span>
                <div>
                  <el-button :loading="btnLoading" type="primary" @click="enterDialog">确 定</el-button>
                  <el-button @click="closeDialog">取 消</el-button>
                </div>
              </div>
            </template>

          <el-form :model="formData" label-position="top" ref="elFormRef" :rules="rule" label-width="80px">
            <el-form-item label="名字:" prop="name">
    <el-input v-model="formData.name" :clearable="false" placeholder="请输入名字" />
</el-form-item>
            <el-form-item label="区域:" prop="region">
    <el-input v-model="formData.region" :clearable="false" placeholder="请输入区域" />
</el-form-item>
            <el-form-item label="CPU:" prop="cpu">
    <el-input v-model="formData.cpu" :clearable="false" placeholder="请输入CPU" />
</el-form-item>
            <el-form-item label="内存:" prop="memory">
    <el-input v-model="formData.memory" :clearable="false" placeholder="请输入内存" />
</el-form-item>
            <el-form-item label="系统盘容量:" prop="systemDiskCapacity">
    <el-input v-model.number="formData.systemDiskCapacity" :clearable="false" placeholder="请输入系统盘容量" />
</el-form-item>
            <el-form-item label="数据盘容量:" prop="dataDiskCapacity">
    <el-input v-model.number="formData.dataDiskCapacity" :clearable="false" placeholder="请输入数据盘容量" />
</el-form-item>
            <el-form-item label="公网IP:" prop="publicIp">
    <el-input v-model="formData.publicIp" :clearable="false" placeholder="请输入公网IP" />
</el-form-item>
            <el-form-item label="内网IP:" prop="privateIp">
    <el-input v-model="formData.privateIp" :clearable="false" placeholder="请输入内网IP" />
</el-form-item>
            <el-form-item label="SSH端口:" prop="sshPort">
    <el-input v-model.number="formData.sshPort" :clearable="false" placeholder="请输入SSH端口" />
</el-form-item>
            <el-form-item label="用户名:" prop="username">
    <el-input v-model="formData.username" :clearable="false" placeholder="请输入用户名" />
</el-form-item>
            <el-form-item label="密码:" prop="password">
    <el-input v-model="formData.password" :clearable="false" placeholder="请输入密码" />
</el-form-item>
            <el-form-item label="显卡名称:" prop="gpuName">
    <el-input v-model="formData.gpuName" :clearable="false" placeholder="请输入显卡名称" />
</el-form-item>
            <el-form-item label="显卡数量:" prop="gpuCount">
    <el-input v-model.number="formData.gpuCount" :clearable="false" placeholder="请输入显卡数量" />
</el-form-item>
            <el-form-item label="Docker连接地址:" prop="dockerConnectAddress">
    <el-input v-model="formData.dockerConnectAddress" :clearable="false" placeholder="请输入Docker连接地址" />
</el-form-item>
            <el-form-item label="使用TLS:" prop="useTls">
    <el-switch v-model="formData.useTls" active-color="#13ce66" inactive-color="#ff4949" active-text="是" inactive-text="否" clearable ></el-switch>
</el-form-item>
            <el-form-item label="CA证书:" prop="caCertificate">
    <el-input v-model="formData.caCertificate" :clearable="false" placeholder="请输入CA证书" type="textarea" :rows="4" />
</el-form-item>
            <el-form-item label="客户端证书:" prop="clientCertificate">
    <el-input v-model="formData.clientCertificate" :clearable="false" placeholder="请输入客户端证书" type="textarea" :rows="4" />
</el-form-item>
            <el-form-item label="客户端私钥:" prop="clientKey">
    <el-input v-model="formData.clientKey" :clearable="false" placeholder="请输入客户端私钥" type="textarea" :rows="4" />
</el-form-item>
            <el-form-item label="是否上架:" prop="isListed">
    <el-switch v-model="formData.isListed" active-color="#13ce66" inactive-color="#ff4949" active-text="是" inactive-text="否" clearable ></el-switch>
</el-form-item>
            <el-form-item label="备注:" prop="remark">
    <el-input v-model="formData.remark" :clearable="false" placeholder="请输入备注" />
</el-form-item>
          </el-form>
    </el-drawer>

    <el-drawer destroy-on-close :size="appStore.drawerSize" v-model="detailShow" :show-close="true" :before-close="closeDetailShow" title="查看">
            <el-descriptions :column="1" border>
                    <el-descriptions-item label="名字">
                        {{ detailForm.name }}
                    </el-descriptions-item>
                    <el-descriptions-item label="区域">
                        {{ detailForm.region }}
                    </el-descriptions-item>
                    <el-descriptions-item label="CPU">
                        {{ detailForm.cpu }}
                    </el-descriptions-item>
                    <el-descriptions-item label="内存">
                        {{ detailForm.memory }}
                    </el-descriptions-item>
                    <el-descriptions-item label="系统盘容量">
                        {{ detailForm.systemDiskCapacity }}
                    </el-descriptions-item>
                    <el-descriptions-item label="数据盘容量">
                        {{ detailForm.dataDiskCapacity }}
                    </el-descriptions-item>
                    <el-descriptions-item label="公网IP">
                        {{ detailForm.publicIp }}
                    </el-descriptions-item>
                    <el-descriptions-item label="内网IP">
                        {{ detailForm.privateIp }}
                    </el-descriptions-item>
                    <el-descriptions-item label="SSH端口">
                        {{ detailForm.sshPort }}
                    </el-descriptions-item>
                    <el-descriptions-item label="用户名">
                        {{ detailForm.username }}
                    </el-descriptions-item>
                    <el-descriptions-item label="密码">
                        {{ detailForm.password }}
                    </el-descriptions-item>
                    <el-descriptions-item label="显卡名称">
                        {{ detailForm.gpuName }}
                    </el-descriptions-item>
                    <el-descriptions-item label="显卡数量">
                        {{ detailForm.gpuCount }}
                    </el-descriptions-item>
                    <el-descriptions-item label="Docker连接地址">
                        {{ detailForm.dockerConnectAddress }}
                    </el-descriptions-item>
                    <el-descriptions-item label="使用TLS">
                        {{ formatBoolean(detailForm.useTls) }}
                    </el-descriptions-item>
                    <el-descriptions-item label="CA证书">
                        {{ detailForm.caCertificate }}
                    </el-descriptions-item>
                    <el-descriptions-item label="客户端证书">
                        {{ detailForm.clientCertificate }}
                    </el-descriptions-item>
                    <el-descriptions-item label="客户端私钥">
                        {{ detailForm.clientKey }}
                    </el-descriptions-item>
                    <el-descriptions-item label="是否上架">
                        {{ formatBoolean(detailForm.isListed) }}
                    </el-descriptions-item>
                    <el-descriptions-item label="备注">
                        {{ detailForm.remark }}
                    </el-descriptions-item>
            </el-descriptions>
        </el-drawer>

  </div>
</template>

<script setup>
import {
  createComputeNode,
  deleteComputeNode,
  deleteComputeNodeByIds,
  updateComputeNode,
  findComputeNode,
  getComputeNodeList
} from '@/api/cloud/compute_node'

// 全量引入格式化工具 请按需保留
import { getDictFunc, formatDate, formatBoolean, filterDict ,filterDataSource, returnArrImg, onDownloadFile } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'
import { useAppStore } from "@/pinia"




defineOptions({
    name: 'ComputeNode'
})

// 提交按钮loading
const btnLoading = ref(false)
const appStore = useAppStore()

// 控制更多查询条件显示/隐藏状态
const showAllQuery = ref(false)

// 自动化生成的字典（可能为空）以及字段
const formData = ref({
            name: '',
            region: '',
            cpu: '',
            memory: '',
            systemDiskCapacity: 0,
            dataDiskCapacity: 0,
            publicIp: '',
            privateIp: '',
            sshPort: 0,
            username: '',
            password: '',
            gpuName: '',
            gpuCount: 0,
            dockerConnectAddress: '',
            useTls: false,
            caCertificate: '',
            clientCertificate: '',
            clientKey: '',
            isListed: false,
            remark: '',
        })



// 验证规则
const rule = reactive({
               name : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               },
               {
                   whitespace: true,
                   message: '不能只输入空格',
                   trigger: ['input', 'blur'],
              }
              ],
               publicIp : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               },
               {
                   whitespace: true,
                   message: '不能只输入空格',
                   trigger: ['input', 'blur'],
              }
              ],
               privateIp : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               },
               {
                   whitespace: true,
                   message: '不能只输入空格',
                   trigger: ['input', 'blur'],
              }
              ],
               sshPort : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               },
              ],
               isListed : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
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
    if (searchInfo.value.useTls === ""){
        searchInfo.value.useTls=null
    }
    if (searchInfo.value.isListed === ""){
        searchInfo.value.isListed=null
    }
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
  const table = await getComputeNodeList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
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
            deleteComputeNodeFunc(row)
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
      const res = await deleteComputeNodeByIds({ IDs })
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
const updateComputeNodeFunc = async(row) => {
    const res = await findComputeNode({ ID: row.ID })
    type.value = 'update'
    if (res.code === 0) {
        formData.value = res.data
        dialogFormVisible.value = true
    }
}


// 删除行
const deleteComputeNodeFunc = async (row) => {
    const res = await deleteComputeNode({ ID: row.ID })
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
    dialogFormVisible.value = true
}

// 关闭弹窗
const closeDialog = () => {
    dialogFormVisible.value = false
    formData.value = {
        name: '',
        region: '',
        cpu: '',
        memory: '',
        systemDiskCapacity: 0,
        dataDiskCapacity: 0,
        publicIp: '',
        privateIp: '',
        sshPort: 0,
        username: '',
        password: '',
        gpuName: '',
        gpuCount: 0,
        dockerConnectAddress: '',
        useTls: false,
        caCertificate: '',
        clientCertificate: '',
        clientKey: '',
        isListed: false,
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
                  res = await createComputeNode(formData.value)
                  break
                case 'update':
                  res = await updateComputeNode(formData.value)
                  break
                default:
                  res = await createComputeNode(formData.value)
                  break
              }
              btnLoading.value = false
              if (res.code === 0) {
                ElMessage({
                  type: 'success',
                  message: '创建/更改成功'
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
  const res = await findComputeNode({ ID: row.ID })
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


</script>

<style>

</style>
