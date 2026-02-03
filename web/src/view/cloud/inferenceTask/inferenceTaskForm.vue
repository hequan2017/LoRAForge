
<template>
  <div>
    <div class="gva-form-box">
      <el-form :model="formData" ref="elFormRef" label-position="right" :rules="rule" label-width="80px">
        <el-form-item label="任务名称:" prop="taskName">
    <el-input v-model="formData.taskName" :clearable="false" placeholder="请输入任务名称" />
</el-form-item>
        <el-form-item label="模型路径:" prop="modelPath">
    <el-input v-model="formData.modelPath" :clearable="false" placeholder="请输入模型路径" />
</el-form-item>
        <el-form-item label="正向提示词:" prop="prompt">
    <el-input v-model="formData.prompt" :clearable="false" placeholder="请输入正向提示词" />
</el-form-item>
        <el-form-item label="反向提示词:" prop="negativePrompt">
    <el-input v-model="formData.negativePrompt" :clearable="false" placeholder="请输入反向提示词" />
</el-form-item>
        <el-form-item label="采样步数:" prop="steps">
    <el-input v-model.number="formData.steps" :clearable="false" placeholder="请输入采样步数" />
</el-form-item>
        <el-form-item label="引导系数:" prop="cfgScale">
    <el-input-number v-model="formData.cfgScale" style="width:100%" :precision="2" :clearable="false" />
</el-form-item>
        <el-form-item label="随机种子:" prop="seed">
    <el-input v-model.number="formData.seed" :clearable="false" placeholder="请输入随机种子" />
</el-form-item>
        <el-form-item label="宽度:" prop="width">
    <el-input v-model.number="formData.width" :clearable="false" placeholder="请输入宽度" />
</el-form-item>
        <el-form-item label="高度:" prop="height">
    <el-input v-model.number="formData.height" :clearable="false" placeholder="请输入高度" />
</el-form-item>
        <el-form-item label="采样器:" prop="sampler">
    <el-input v-model="formData.sampler" :clearable="false" placeholder="请输入采样器" />
</el-form-item>
        <el-form-item label="状态:" prop="status">
    <el-tree-select v-model="formData.status" placeholder="请选择状态" :data="inference_statusOptions" style="width:100%" filterable :clearable="false" check-strictly></el-tree-select>
</el-form-item>
        <el-form-item label="生成结果:" prop="resultImage">
    <SelectImage
     v-model="formData.resultImage"
     file-type="image"
    />
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
  createInferenceTask,
  updateInferenceTask,
  findInferenceTask
} from '@/api/cloud/inferenceTask'

defineOptions({
    name: 'InferenceTaskForm'
})

// 自动获取字典
import { getDictFunc } from '@/utils/format'
import { useRoute, useRouter } from "vue-router"
import { ElMessage } from 'element-plus'
import { ref, reactive } from 'vue'
// 图片选择组件
import SelectImage from '@/components/selectImage/selectImage.vue'


const route = useRoute()
const router = useRouter()

// 提交按钮loading
const btnLoading = ref(false)

const type = ref('')
const inference_statusOptions = ref([])
const formData = ref({
            taskName: '',
            modelPath: '',
            prompt: '',
            negativePrompt: '',
            steps: 0,
            cfgScale: 0,
            seed: 0,
            width: 0,
            height: 0,
            sampler: '',
            status: '',
            resultImage: "",
        })
// 验证规则
const rule = reactive({
               taskName : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               modelPath : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               prompt : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               steps : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               cfgScale : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               seed : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               width : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               height : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               sampler : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               status : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
})

const elFormRef = ref()

// 初始化方法
const init = async () => {
 // 建议通过url传参获取目标数据ID 调用 find方法进行查询数据操作 从而决定本页面是create还是update 以下为id作为url参数示例
    if (route.query.id) {
      const res = await findInferenceTask({ ID: route.query.id })
      if (res.code === 0) {
        formData.value = res.data
        type.value = 'update'
      }
    } else {
      type.value = 'create'
    }
    inference_statusOptions.value = await getDictFunc('inference_status')
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
               res = await createInferenceTask(formData.value)
               break
             case 'update':
               res = await updateInferenceTask(formData.value)
               break
             default:
               res = await createInferenceTask(formData.value)
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
