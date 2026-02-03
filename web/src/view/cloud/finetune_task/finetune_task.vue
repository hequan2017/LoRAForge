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
        <el-form-item label="任务名称" prop="taskName">
          <el-input v-model="searchInfo.taskName" placeholder="搜索任务名称" />
        </el-form-item>
        <template v-if="showAllQuery">
          <el-form-item label="任务状态" prop="taskStatus">
            <el-select v-model="searchInfo.taskStatus" placeholder="请选择状态" clearable>
              <el-option label="等待中" value="pending" />
              <el-option label="运行中" value="running" />
              <el-option label="已完成" value="completed" />
              <el-option label="失败" value="failed" />
              <el-option label="已停止" value="stopped" />
              <el-option label="已取消" value="cancelled" />
            </el-select>
          </el-form-item>
          <el-form-item label="模型类型" prop="modelType">
            <el-select v-model="searchInfo.modelType" placeholder="请选择模型类型" clearable>
              <el-option label="LLM" value="llm" />
              <el-option label="WA" value="wa" />
              <el-option label="Audio" value="audio" />
            </el-select>
          </el-form-item>
          <el-form-item label="节点" prop="nodeId">
            <el-select v-model="searchInfo.nodeId" placeholder="请选择节点" filterable clearable>
              <el-option v-for="(item,key) in dataSource.nodeId" :key="key" :label="item.label" :value="item.value" />
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
        <el-button type="primary" icon="plus" @click="openDialog()">新建任务</el-button>
        <el-button icon="delete" style="margin-left: 10px;" :disabled="!multipleSelection.length" @click="onDelete">批量删除</el-button>
        <el-button icon="refresh" style="margin-left: 10px;" @click="syncAllStatus">同步状态</el-button>
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
        <el-table-column align="left" label="任务名称" prop="taskName" width="200" />
        <el-table-column align="left" label="模型类型" prop="modelType" width="100">
          <template #default="scope">
            <el-tag v-if="scope.row.modelType === 'llm'" type="primary" size="small">LLM</el-tag>
            <el-tag v-else-if="scope.row.modelType === 'wa'" type="success" size="small">WA</el-tag>
            <el-tag v-else-if="scope.row.modelType === 'audio'" type="warning" size="small">Audio</el-tag>
            <el-tag v-else type="info" size="small">{{ scope.row.modelType }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="模型" prop="modelPath" width="200">
          <template #default="scope">
            <el-tooltip :content="scope.row.modelPath" placement="top">
              <span class="text-xs text-gray-500 truncate block max-w-[180px]">{{ scope.row.modelPath }}</span>
            </el-tooltip>
          </template>
        </el-table-column>
        <el-table-column align="left" label="进度" prop="progress" width="150">
          <template #default="scope">
            <div class="flex items-center gap-2">
              <el-progress :percentage="Math.round(scope.row.progress || 0)" :stroke-width="8" />
            </div>
            <div class="text-xs text-gray-400 mt-1">
              {{ scope.row.currentStep || 0 }} / {{ scope.row.totalSteps || '-' }} 步
            </div>
          </template>
        </el-table-column>
        <el-table-column align="left" label="状态" prop="taskStatus" width="100">
          <template #default="scope">
            <el-tag v-if="scope.row.taskStatus === 'running'" type="success">运行中</el-tag>
            <el-tag v-else-if="scope.row.taskStatus === 'pending'" type="info">等待中</el-tag>
            <el-tag v-else-if="scope.row.taskStatus === 'completed'" type="success">已完成</el-tag>
            <el-tag v-else-if="scope.row.taskStatus === 'failed'" type="danger">失败</el-tag>
            <el-tag v-else-if="scope.row.taskStatus === 'stopped'" type="warning">已停止</el-tag>
            <el-tag v-else-if="scope.row.taskStatus === 'cancelled'" type="info">已取消</el-tag>
            <el-tag v-else type="info">{{ scope.row.taskStatus }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="操作" fixed="right" :min-width="appStore.operateMinWith">
          <template #default="scope">
            <el-button v-if="scope.row.taskStatus === 'pending'" type="success" link icon="VideoPlay" @click="startTaskFunc(scope.row)">启动</el-button>
            <el-button v-if="['failed', 'stopped', 'cancelled'].includes(scope.row.taskStatus)" type="success" link icon="RefreshRight" @click="restartTaskFunc(scope.row)">重新启动</el-button>
            <el-button v-if="scope.row.taskStatus === 'running'" type="warning" link icon="SwitchButton" @click="stopTaskFunc(scope.row)">停止</el-button>
            <el-button v-if="scope.row.taskStatus === 'running'" type="info" link class="table-button" @click="openLog(scope.row)"><el-icon style="margin-right: 5px"><Document /></el-icon>日志</el-button>
            <el-button type="primary" link class="table-button" @click="getDetails(scope.row)"><el-icon style="margin-right: 5px"><InfoFilled /></el-icon>详情</el-button>
            <el-button type="primary" link icon="edit" class="table-button" @click="updateTaskFunc(scope.row)">编辑</el-button>
            <el-button type="primary" link icon="delete" @click="deleteRow(scope.row)">删除</el-button>
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
    <el-drawer destroy-on-close :size="700" v-model="dialogFormVisible" :show-close="false" :before-close="closeDialog">
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg font-semibold">{{type==='create'?'创建微调任务':'编辑微调任务'}}</span>
          <div>
            <el-button :loading="btnLoading" type="primary" @click="enterDialog">{{type==='create'?'创建':'保存'}}</el-button>
            <el-button @click="closeDialog">取消</el-button>
          </div>
        </div>
      </template>

      <el-form :model="formData" label-position="top" ref="elFormRef" :rules="rule" label-width="120px">
        <el-tabs v-model="activeTab" class="w-full">
          <!-- 基础配置 -->
          <el-tab-pane label="基础配置" name="basic">
            <el-row :gutter="16">
              <el-col :span="12">
                <el-form-item label="任务名称" prop="taskName">
                  <el-input v-model="formData.taskName" placeholder="请输入任务名称" clearable />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="计算节点" prop="nodeId">
                  <el-select v-model="formData.nodeId" placeholder="请选择节点" filterable style="width:100%">
                    <el-option v-for="(item,key) in dataSource.nodeId" :key="key" :label="item.label" :value="item.value" />
                  </el-select>
                </el-form-item>
              </el-col>
            </el-row>

            <el-form-item label="任务描述">
              <el-input v-model="formData.taskDescription" type="textarea" :rows="2" placeholder="请输入任务描述" />
            </el-form-item>
          </el-tab-pane>

          <!-- 模型配置 -->
          <el-tab-pane label="模型配置" name="model">
            <el-row :gutter="16">
              <el-col :span="12">
                <el-form-item label="模型类型" prop="modelType">
                  <el-select v-model="formData.modelType" placeholder="请选择模型类型" style="width:100%">
                    <el-option label="LLM (大语言模型)" value="llm" />
                    <el-option label="WA (文生图)" value="wa" />
                    <el-option label="Audio (语音)" value="audio" />
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="模型路径" prop="modelPath">
                  <el-input v-model="formData.modelPath" placeholder="模型ID或路径，如: qwen/Qwen-7B" clearable />
                </el-form-item>
              </el-col>
            </el-row>

            <el-row :gutter="16">
              <el-col :span="8">
                <el-form-item label="LoRA秩">
                  <el-input-number v-model="formData.loraRank" :min="1" :max="256" :step="1" controls-position="right" style="width:100%" />
                </el-form-item>
              </el-col>
              <el-col :span="8">
                <el-form-item label="LoRA Alpha">
                  <el-input-number v-model="formData.loraAlpha" :min="1" :max="512" :step="1" controls-position="right" style="width:100%" />
                </el-form-item>
              </el-col>
              <el-col :span="8">
                <el-form-item label="LoRA Dropout">
                  <el-input-number v-model="formData.loraDropRate" :min="0" :max="0.5" :step="0.01" :precision="2" controls-position="right" style="width:100%" />
                </el-form-item>
              </el-col>
            </el-row>

            <el-form-item label="LoRA目标模块">
              <el-input v-model="formData.loraTargetModules" placeholder="如: c_attn,q_attn，默认自动选择" clearable />
            </el-form-item>
          </el-tab-pane>

          <!-- 数据集配置 -->
          <el-tab-pane label="数据集配置" name="dataset">
            <el-form-item label="训练数据集" prop="trainDataset">
              <el-input v-model="formData.trainDataset" placeholder="数据集ID或路径，如: dataset/xxx" clearable>
                <template #append>
                  <el-tooltip content="支持ModelScope数据集ID或本地路径" placement="top">
                    <el-icon><QuestionFilled /></el-icon>
                  </el-tooltip>
                </template>
              </el-input>
            </el-form-item>

            <el-form-item label="验证数据集">
              <el-input v-model="formData.valDataset" placeholder="验证数据集ID或路径（可选）" clearable />
            </el-form-item>

            <el-row :gutter="16">
              <el-col :span="12">
                <el-form-item label="数据集类型">
                  <el-select v-model="formData.datasetType" placeholder="请选择数据集类型" style="width:100%">
                    <el-option label="自定义" value="custom" />
                    <el-option label="Alpaca" value="alpaca" />
                    <el-option label="ShareGPT" value="sharegpt" />
                    <el-option label="MOSS" value="moss" />
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="最大样本数">
                  <el-input-number v-model="formData.maxSamples" :min="0" :step="1000" controls-position="right" style="width:100%" />
                </el-form-item>
              </el-col>
            </el-row>
          </el-tab-pane>

          <!-- 训练参数 -->
          <el-tab-pane label="训练参数" name="training">
            <el-row :gutter="16">
              <el-col :span="8">
                <el-form-item label="学习率">
                  <el-input-number v-model="formData.learningRate" :min="0.00001" :max="0.01" :step="0.0001" :precision="5" controls-position="right" style="width:100%" />
                </el-form-item>
              </el-col>
              <el-col :span="8">
                <el-form-item label="批处理大小">
                  <el-input-number v-model="formData.batchSize" :min="1" :max="256" :step="1" controls-position="right" style="width:100%" />
                </el-form-item>
              </el-col>
              <el-col :span="8">
                <el-form-item label="梯度累积步数">
                  <el-input-number v-model="formData.gradientAccSteps" :min="1" :max="32" :step="1" controls-position="right" style="width:100%" />
                </el-form-item>
              </el-col>
            </el-row>

            <el-row :gutter="16">
              <el-col :span="8">
                <el-form-item label="训练轮数">
                  <el-input-number v-model="formData.numEpochs" :min="1" :max="100" :step="1" controls-position="right" style="width:100%" />
                </el-form-item>
              </el-col>
              <el-col :span="8">
                <el-form-item label="最大步数">
                  <el-input-number v-model="formData.maxSteps" :min="0" :step="100" controls-position="right" style="width:100%" />
                </el-form-item>
              </el-col>
              <el-col :span="8">
                <el-form-item label="预热比例">
                  <el-input-number v-model="formData.warmupRatio" :min="0" :max="0.5" :step="0.05" :precision="2" controls-position="right" style="width:100%" />
                </el-form-item>
              </el-col>
            </el-row>

            <el-row :gutter="16">
              <el-col :span="8">
                <el-form-item label="优化器">
                  <el-select v-model="formData.optimizer" placeholder="请选择优化器" style="width:100%">
                    <el-option label="AdamW" value="adamw" />
                    <el-option label="Adam" value="adam" />
                    <el-option label="Lion" value="lion" />
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="8">
                <el-form-item label="权重衰减">
                  <el-input-number v-model="formData.weightDecay" :min="0" :max="0.1" :step="0.001" :precision="3" controls-position="right" style="width:100%" />
                </el-form-item>
              </el-col>
              <el-col :span="8">
                <el-form-item label="梯度裁剪">
                  <el-input-number v-model="formData.maxGradNorm" :min="0" :max="10" :step="0.1" :precision="1" controls-position="right" style="width:100%" />
                </el-form-item>
              </el-col>
            </el-row>
          </el-tab-pane>

          <!-- 环境配置 -->
          <el-tab-pane label="环境配置" name="environment">
            <el-row :gutter="16">
              <el-col :span="12">
                <el-form-item label="GPU数量">
                  <el-input-number v-model="formData.gpuCount" :min="0" :max="8" :step="1" controls-position="right" style="width:100%" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="GPU类型">
                  <el-select v-model="formData.gpuType" placeholder="请选择GPU类型" style="width:100%">
                    <el-option label="A100" value="A100" />
                    <el-option label="A800" value="A800" />
                    <el-option label="H800" value="H800" />
                    <el-option label="H100" value="H100" />
                    <el-option label="V100" value="V100" />
                    <el-option label="RTX 4090" value="RTX4090" />
                    <el-option label="RTX 3090" value="RTX3090" />
                  </el-select>
                </el-form-item>
              </el-col>
            </el-row>

            <el-row :gutter="16">
              <el-col :span="8">
                <el-form-item label="训练精度">
                  <el-select v-model="formData.precision" placeholder="请选择精度" style="width:100%">
                    <el-option label="BF16" value="bf16" />
                    <el-option label="FP16" value="fp16" />
                    <el-option label="FP32" value="fp32" />
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="8">
                <el-form-item label="序列长度">
                  <el-input-number v-model="formData.sequenceLength" :min="128" :max="32768" :step="128" controls-position="right" style="width:100%" />
                </el-form-item>
              </el-col>
              <el-col :span="8">
                <el-form-item label="量化位数">
                  <el-select v-model="formData.quantizationBit" placeholder="可选" style="width:100%">
                    <el-option label="不量化" :value="0" />
                    <el-option label="4位" :value="4" />
                    <el-option label="8位" :value="8" />
                  </el-select>
                </el-form-item>
              </el-col>
            </el-row>

            <el-row :gutter="16">
              <el-col :span="24">
                <el-form-item label="高级选项">
                  <el-checkbox v-model="formData.flashAttention">启用 Flash Attention</el-checkbox>
                  <el-checkbox v-model="formData.deepspeed" style="margin-left: 16px">启用 DeepSpeed</el-checkbox>
                  <el-checkbox v-model="formData.onlySaveModel" style="margin-left: 16px">仅保存最终模型</el-checkbox>
                </el-form-item>
              </el-col>
            </el-row>
          </el-tab-pane>

          <!-- 输出配置 -->
          <el-tab-pane label="输出配置" name="output">
            <el-form-item label="输出目录">
              <el-input v-model="formData.outputDir" placeholder="如: /data/output/finetune/task001" clearable />
            </el-form-item>

            <el-row :gutter="16">
              <el-col :span="12">
                <el-form-item label="检查点保存步数">
                  <el-input-number v-model="formData.checkpointStep" :min="0" :step="100" controls-position="right" style="width:100%" />
                </el-form-item>
              </el-col>
            </el-row>

            <el-form-item label="备注">
              <el-input v-model="formData.remark" type="textarea" :rows="2" placeholder="请输入备注信息" />
            </el-form-item>
          </el-tab-pane>
        </el-tabs>
      </el-form>
    </el-drawer>

    <!-- 详情抽屉 -->
    <el-drawer destroy-on-close :size="600" v-model="detailShow" :show-close="true" :before-close="closeDetailShow" title="任务详情">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="任务名称" :span="2">
          <span class="font-semibold">{{ detailForm.taskName || '-' }}</span>
        </el-descriptions-item>
        <el-descriptions-item label="任务状态">
          <el-tag v-if="detailForm.taskStatus === 'running'" type="success" size="small">运行中</el-tag>
          <el-tag v-else-if="detailForm.taskStatus === 'pending'" type="info" size="small">等待中</el-tag>
          <el-tag v-else-if="detailForm.taskStatus === 'completed'" type="success" size="small">已完成</el-tag>
          <el-tag v-else-if="detailForm.taskStatus === 'failed'" type="danger" size="small">失败</el-tag>
          <el-tag v-else-if="detailForm.taskStatus === 'stopped'" type="warning" size="small">已停止</el-tag>
          <el-tag v-else type="info" size="small">{{ detailForm.taskStatus }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="进度">
          {{ detailForm.progress ? detailForm.progress.toFixed(1) + '%' : '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="模型类型">
          {{ detailForm.modelType || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="模型路径">
          {{ detailForm.modelPath || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="容器ID" :span="2">
          <span class="font-mono text-xs">{{ detailForm.containerId || '-' }}</span>
        </el-descriptions-item>
        <el-descriptions-item label="训练数据集" :span="2">
          {{ detailForm.trainDataset || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="学习率">
          {{ detailForm.learningRate || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="批处理大小">
          {{ detailForm.batchSize || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="当前步数">
          {{ detailForm.currentStep || 0 }} / {{ detailForm.totalSteps || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="验证损失">
          {{ detailForm.validationLoss || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="错误信息" :span="2" v-if="detailForm.errorMessage">
          <el-text type="danger">{{ detailForm.errorMessage }}</el-text>
        </el-descriptions-item>
        <el-descriptions-item label="备注" :span="2">
          {{ detailForm.remark || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="创建时间" :span="2">
          {{ formatDate(detailForm.CreatedAt) }}
        </el-descriptions-item>
      </el-descriptions>
    </el-drawer>

    <!-- 日志弹窗 -->
    <el-dialog
      v-model="logVisible"
      :fullscreen="false"
      :show-close="true"
      @close="closeLog"
      width="90%"
      destroy-on-close
      title="训练日志">
      <div class="log-wrapper">
        <el-scrollbar height="500px">
          <pre class="log-content">{{ logContent || '暂无日志' }}</pre>
        </el-scrollbar>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import {
  createFineTuneTask,
  deleteFineTuneTask,
  deleteFineTuneTaskByIds,
  updateFineTuneTask,
  findFineTuneTask,
  getFineTuneTaskList,
  startFineTuneTask,
  stopFineTuneTask,
  restartFineTuneTask,
  getFineTuneTaskLogs,
  getFineTuneTaskDataSource,
  syncAllFineTuneTaskStatus
} from '@/api/cloud/finetune_task'

import { formatDate, filterDataSource } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive, computed, onMounted, onUnmounted } from 'vue'
import { useAppStore } from "@/pinia"

defineOptions({
  name: 'FineTuneTask'
})

const appStore = useAppStore()
const btnLoading = ref(false)
const showAllQuery = ref(false)
const activeTab = ref('basic')
const logVisible = ref(false)
const logContent = ref('')
const currentTask = ref(null)

// 自动刷新定时器
let refreshTimer = null

// 表单数据
const formData = ref({
  taskName: '',
  taskDescription: '',
  nodeId: undefined,
  modelType: 'llm',
  modelPath: '',
  loraTargetModules: '',
  loraRank: 8,
  loraAlpha: 32,
  loraDropRate: 0.05,
  trainDataset: '',
  valDataset: '',
  datasetType: 'alpaca',
  datasetProb: 1.0,
  maxSamples: undefined,
  learningRate: 0.0002,
  batchSize: 8,
  gradientAccSteps: 1,
  numEpochs: 1,
  maxSteps: undefined,
  warmupRatio: 0.1,
  optimizer: 'adamw',
  weightDecay: 0.01,
  maxGradNorm: 1.0,
  gpuCount: 1,
  gpuType: '',
  precision: 'bf16',
  quantizationBit: 0,
  flashAttention: true,
  deepspeed: false,
  sequenceLength: 2048,
  outputDir: '',
  checkpointStep: undefined,
  onlySaveModel: false,
  remark: ''
})

const dataSource = ref({})
const getDataSourceFunc = async () => {
  const res = await getFineTuneTaskDataSource()
  if (res.code === 0) {
    dataSource.value = res.data
  }
}
getDataSourceFunc()

// 验证规则
const rule = reactive({
  taskName: [{
    required: true,
    message: '请输入任务名称',
    trigger: ['blur'],
  }],
  nodeId: [{
    required: true,
    message: '请选择计算节点',
    trigger: ['change', 'blur'],
  }],
  modelType: [{
    required: true,
    message: '请选择模型类型',
    trigger: ['change', 'blur'],
  }],
  modelPath: [{
    required: true,
    message: '请输入模型路径',
    trigger: ['blur'],
  }],
  trainDataset: [{
    required: true,
    message: '请输入训练数据集',
    trigger: ['blur'],
  }],
})

const elFormRef = ref()
const elSearchFormRef = ref()

// 表格控制
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})

const onReset = () => {
  searchInfo.value = {}
  getTableData()
}

const onSubmit = () => {
  page.value = 1
  getTableData()
}

const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

const getTableData = async () => {
  const table = await getFineTuneTaskList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

getTableData()

// 多选
const multipleSelection = ref([])
const handleSelectionChange = (val) => {
  multipleSelection.value = val
}

// 删除
const deleteRow = (row) => {
  ElMessageBox.confirm('确定要删除吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    deleteTaskFunc(row)
  })
}

const onDelete = async () => {
  ElMessageBox.confirm('确定要删除吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
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
    const res = await deleteFineTuneTaskByIds({ ids: IDs })
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

const deleteTaskFunc = async (row) => {
  const res = await deleteFineTuneTask({ id: row.ID })
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

const type = ref('')

const updateTaskFunc = async (row) => {
  const res = await findFineTuneTask({ id: row.ID })
  type.value = 'update'
  if (res.code === 0) {
    formData.value = res.data
    // 处理checkbox字段
    if (res.data.flashAttention === null) formData.value.flashAttention = true
    if (res.data.deepspeed === null) formData.value.deepspeed = false
    if (res.data.onlySaveModel === null) formData.value.onlySaveModel = false
    dialogFormVisible.value = true
  }
}

const dialogFormVisible = ref(false)

const openDialog = () => {
  type.value = 'create'
  activeTab.value = 'basic'
  dialogFormVisible.value = true
}

const closeDialog = () => {
  dialogFormVisible.value = false
  formData.value = {
    taskName: '',
    taskDescription: '',
    nodeId: undefined,
    modelType: 'llm',
    modelPath: '',
    loraTargetModules: '',
    loraRank: 8,
    loraAlpha: 32,
    loraDropRate: 0.05,
    trainDataset: '',
    valDataset: '',
    datasetType: 'alpaca',
    datasetProb: 1.0,
    maxSamples: undefined,
    learningRate: 0.0002,
    batchSize: 8,
    gradientAccSteps: 1,
    numEpochs: 1,
    maxSteps: undefined,
    warmupRatio: 0.1,
    optimizer: 'adamw',
    weightDecay: 0.01,
    maxGradNorm: 1.0,
    gpuCount: 1,
    gpuType: '',
    precision: 'bf16',
    quantizationBit: 0,
    flashAttention: true,
    deepspeed: false,
    sequenceLength: 2048,
    outputDir: '',
    checkpointStep: undefined,
    onlySaveModel: false,
    remark: ''
  }
}

const enterDialog = async () => {
  btnLoading.value = true
  elFormRef.value?.validate(async (valid) => {
    if (!valid) return btnLoading.value = false
    let res
    switch (type.value) {
      case 'create':
        res = await createFineTuneTask(formData.value)
        break
      case 'update':
        res = await updateFineTuneTask(formData.value)
        break
      default:
        res = await createFineTuneTask(formData.value)
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
const detailShow = ref(false)

const openDetailShow = () => {
  detailShow.value = true
}

const getDetails = async (row) => {
  const res = await findFineTuneTask({ id: row.ID })
  if (res.code === 0) {
    detailForm.value = res.data
    openDetailShow()
  }
}

const closeDetailShow = () => {
  detailShow.value = false
  detailForm.value = {}
}

// 启动任务
const startTaskFunc = (row) => {
  ElMessageBox.confirm('确定要启动此微调任务吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'success'
  }).then(async () => {
    const res = await startFineTuneTask({ id: row.ID })
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: '任务启动成功'
      })
      getTableData()
    }
  })
}

// 停止任务
const stopTaskFunc = (row) => {
  ElMessageBox.confirm('确定要停止此微调任务吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    const res = await stopFineTuneTask({ id: row.ID, reason: '用户手动停止' })
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: '任务已停止'
      })
      getTableData()
    }
  })
}

// 重启任务
const restartTaskFunc = (row) => {
  const statusMap = {
    'failed': '失败',
    'stopped': '已停止',
    'cancelled': '已取消'
  }
  const statusText = statusMap[row.taskStatus] || row.taskStatus

  ElMessageBox.confirm(`确定要重新启动此${statusText}的微调任务吗?`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'info'
  }).then(async () => {
    const res = await restartFineTuneTask({ id: row.ID })
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: '任务正在重新启动'
      })
      getTableData()
    }
  })
}

// 同步状态
const syncAllStatus = async () => {
  const res = await syncAllFineTuneTaskStatus()
  if (res.code === 0) {
    ElMessage.success('状态同步成功')
    getTableData()
  }
}

// 查看日志
const openLog = async (row) => {
  currentTask.value = row
  logVisible.value = true
  logContent.value = '正在加载日志...'
  const res = await getFineTuneTaskLogs({ id: row.ID, offset: 0, limit: 500 })
  if (res.code === 0) {
    logContent.value = res.data.join('\n') || '暂无日志'
  } else {
    logContent.value = '日志加载失败'
  }
}

const closeLog = () => {
  logVisible.value = false
  currentTask.value = null
  logContent.value = ''
}

// 自动刷新运行中的任务
const startAutoRefresh = () => {
  refreshTimer = setInterval(() => {
    const hasRunning = tableData.value.some(task => task.taskStatus === 'running')
    if (hasRunning) {
      getTableData()
    }
  }, 10000) // 每10秒刷新一次
}

const stopAutoRefresh = () => {
  if (refreshTimer) {
    clearInterval(refreshTimer)
    refreshTimer = null
  }
}

onMounted(() => {
  startAutoRefresh()
})

onUnmounted(() => {
  stopAutoRefresh()
})
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

.log-wrapper {
  background-color: #1e1e1e;
  padding: 16px;
  border-radius: 4px;
}

.log-content {
  color: #d4d4d4;
  font-family: 'Courier New', monospace;
  font-size: 12px;
  white-space: pre-wrap;
  word-break: break-all;
  margin: 0;
}

.truncate {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
