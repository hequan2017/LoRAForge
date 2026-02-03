
<template>
  <div>
    <div class="gva-form-box">
      <el-form :model="formData" ref="elFormRef" label-position="right" :rules="rule" label-width="80px">
        <el-form-item label="任务名称:" prop="taskName">
    <el-input v-model="formData.taskName" :clearable="false" placeholder="请输入任务名称" />
</el-form-item>
        <el-form-item label="计算节点:" prop="nodeId">
    <el-select v-model="formData.nodeId" placeholder="请选择计算节点" filterable style="width:100%" :clearable="false">
        <el-option v-for="(item,key) in dataSource.nodeId" :key="key" :label="item.label" :value="item.value" />
    </el-select>
</el-form-item>
        <el-form-item label="语言:" prop="language">
    <el-tree-select v-model="formData.language" placeholder="请选择语言" :data="swift_languageOptions" style="width:100%" filterable :clearable="false" check-strictly></el-tree-select>
</el-form-item>
        <el-form-item label="端口:" prop="port">
    <el-input v-model.number="formData.port" :clearable="false" placeholder="请输入端口" />
</el-form-item>
        <el-form-item label="状态:" prop="status">
    <el-tree-select v-model="formData.status" placeholder="请选择状态" :data="swift_statusOptions" style="width:100%" filterable :clearable="false" check-strictly></el-tree-select>
</el-form-item>
        <el-form-item label="访问地址:" prop="accessUrl">
    <el-input v-model="formData.accessUrl" :clearable="false" placeholder="请输入访问地址" />
</el-form-item>
        <el-form-item>
          <el-button :loading="btnLoading" type="primary" @click="save">保存</el-button>
          <el-button type="primary" @click="back">返回</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import {
    getSwiftWebUIDataSource,
  createSwiftWebUI,
  updateSwiftWebUI,
  findSwiftWebUI
} from '@/api/cloud/swift_web_ui'

defineOptions({
    name: 'SwiftWebUIForm'
})

// 自动获取字典
import { getDictFunc } from '@/utils/format'
import { useRoute, useRouter } from "vue-router"
import { ElMessage } from 'element-plus'
import { ref, reactive } from 'vue'


const route = useRoute()
const router = useRouter()

// 提交按钮loading
const btnLoading = ref(false)

const type = ref('')
const swift_languageOptions = ref([])
const swift_statusOptions = ref([])
const formData = ref({
            taskName: '',
            nodeId: undefined,
            language: '',
            port: 0,
            status: '',
            accessUrl: '',
        })
// 验证规则
const rule = reactive({
               taskName : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               nodeId : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
})

const elFormRef = ref()
  const dataSource = ref([])
  const getDataSourceFunc = async()=>{
    const res = await getSwiftWebUIDataSource()
    if (res.code === 0) {
      dataSource.value = res.data
    }
  }
  getDataSourceFunc()

// 初始化方法
const init = async () => {
 // 建议通过url传参获取目标数据ID 调用 find方法进行查询数据操作 从而决定本页面是create还是update 以下为id作为url参数示例
    if (route.query.id) {
      const res = await findSwiftWebUI({ ID: route.query.id })
      if (res.code === 0) {
        formData.value = res.data
        type.value = 'update'
      }
    } else {
      type.value = 'create'
    }
    swift_languageOptions.value = await getDictFunc('swift_language')
    swift_statusOptions.value = await getDictFunc('swift_status')
}

init()
// 保存按钮
const save = async() => {
      btnLoading.value = true
      elFormRef.value?.validate( async (valid) => {
         if (!valid) return btnLoading.value = false
            let res
           switch (type.value) {
             case 'create':
               res = await createSwiftWebUI(formData.value)
               break
             case 'update':
               res = await updateSwiftWebUI(formData.value)
               break
             default:
               res = await createSwiftWebUI(formData.value)
               break
           }
           btnLoading.value = false
           if (res.code === 0) {
             ElMessage({
               type: 'success',
               message: '创建/更改成功'
             })
           }
       })
}

// 返回按钮
const back = () => {
    router.go(-1)
}

</script>

<style>
</style>
