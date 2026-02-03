<template>
  <el-dialog
    v-model="visible"
    title="容器监控"
    width="800px"
    :before-close="handleClose"
    :destroy-on-close="true"
  >
    <div class="monitor-container">
      <el-row :gutter="20">
        <el-col :span="12">
          <div ref="cpuChartRef" class="chart-box"></div>
        </el-col>
        <el-col :span="12">
          <div ref="memChartRef" class="chart-box"></div>
        </el-col>
      </el-row>
      <el-row :gutter="20" class="mt-4">
        <el-col :span="12">
          <div ref="netChartRef" class="chart-box"></div>
        </el-col>
        <el-col :span="12">
          <div ref="blkChartRef" class="chart-box"></div>
        </el-col>
      </el-row>
    </div>
  </el-dialog>
</template>

<script setup>
import { ref, watch, onUnmounted, nextTick } from 'vue'
import * as echarts from 'echarts'
import { useUserStore } from '@/pinia'

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

const cpuChartRef = ref(null)
const memChartRef = ref(null)
const netChartRef = ref(null)
const blkChartRef = ref(null)

let cpuChart = null
let memChart = null
let netChart = null
let blkChart = null
let eventSource = null

// 数据缓冲
const maxPoints = 60
const timestamps = ref([])
const cpuData = ref([])
const memData = ref([])
const netRxData = ref([])
const netTxData = ref([])
const blkReadData = ref([])
const blkWriteData = ref([])

watch(() => props.modelValue, (val) => {
  visible.value = val
  if (val && props.instance?.ID) {
    nextTick(() => {
      initCharts()
      startMonitoring()
    })
  } else {
    stopMonitoring()
  }
})

const handleClose = () => {
  emit('update:modelValue', false)
}

const initCharts = () => {
  if (cpuChartRef.value) {
    cpuChart = echarts.init(cpuChartRef.value)
    cpuChart.setOption(getChartOption('CPU 使用率 (%)', '#409EFF'))
  }
  if (memChartRef.value) {
    memChart = echarts.init(memChartRef.value)
    memChart.setOption(getChartOption('内存使用率 (%)', '#67C23A'))
  }
  if (netChartRef.value) {
    netChart = echarts.init(netChartRef.value)
    netChart.setOption(getDualChartOption('网络流量 (KB/s)', ['接收', '发送'], ['#E6A23C', '#F56C6C']))
  }
  if (blkChartRef.value) {
    blkChart = echarts.init(blkChartRef.value)
    blkChart.setOption(getDualChartOption('磁盘 I/O (KB/s)', ['读取', '写入'], ['#409EFF', '#67C23A']))
  }
}

const getChartOption = (title, color) => ({
  title: { text: title, left: 'center', textStyle: { fontSize: 14, fontWeight: 'normal' } },
  tooltip: { trigger: 'axis' },
  grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
  xAxis: { type: 'category', boundaryGap: false, data: [] },
  yAxis: { type: 'value', min: 0, max: 100, splitLine: { lineStyle: { type: 'dashed' } } },
  series: [{ type: 'line', smooth: true, showSymbol: false, data: [], itemStyle: { color }, areaStyle: { color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [{ offset: 0, color }, { offset: 1, color: '#fff' }]) } }]
})

const getDualChartOption = (title, names, colors) => ({
  title: { text: title, left: 'center', textStyle: { fontSize: 14, fontWeight: 'normal' } },
  tooltip: { trigger: 'axis' },
  legend: { data: names, top: 25 },
  grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
  xAxis: { type: 'category', boundaryGap: false, data: [] },
  yAxis: { type: 'value', splitLine: { lineStyle: { type: 'dashed' } } },
  series: [
    { name: names[0], type: 'line', smooth: true, showSymbol: false, data: [], itemStyle: { color: colors[0] }, areaStyle: { opacity: 0.1, color: colors[0] } },
    { name: names[1], type: 'line', smooth: true, showSymbol: false, data: [], itemStyle: { color: colors[1] }, areaStyle: { opacity: 0.1, color: colors[1] } }
  ]
})

const startMonitoring = () => {
  // 清空数据
  timestamps.value = []
  cpuData.value = []
  memData.value = []
  netRxData.value = []
  netTxData.value = []

  const baseUrl = import.meta.env.VITE_BASE_API || ''
  const url = `${baseUrl}/inst/stats?nodeId=${props.instance.NodeID}&containerId=${props.instance.DockerContainer}&x-token=${userStore.token}`
  
  eventSource = new EventSource(url)
  
  let lastRx = 0
  let lastTx = 0
  let lastRead = 0
  let lastWrite = 0
  let isFirst = true

  eventSource.onmessage = (event) => {
    try {
      const data = JSON.parse(event.data)
      const timeStr = new Date(data.timestamp * 1000).toLocaleTimeString()
      
      // 维护队列长度
      if (timestamps.value.length > maxPoints) {
        timestamps.value.shift()
        cpuData.value.shift()
        memData.value.shift()
        netRxData.value.shift()
        netTxData.value.shift()
        blkReadData.value.shift()
        blkWriteData.value.shift()
      }

      timestamps.value.push(timeStr)
      cpuData.value.push(data.cpu_percent.toFixed(2))
      memData.value.push(data.mem_percent.toFixed(2))

      // 计算速率 (简单差值)
      const currentRx = data.net_rx
      const currentTx = data.net_tx
      const currentRead = data.blk_read
      const currentWrite = data.blk_write

      if (isFirst) {
        netRxData.value.push(0)
        netTxData.value.push(0)
        blkReadData.value.push(0)
        blkWriteData.value.push(0)
        isFirst = false
      } else {
        // 假设大约1秒一次，直接转为 KB/s
        netRxData.value.push(Math.max(0, ((currentRx - lastRx) / 1024)).toFixed(2))
        netTxData.value.push(Math.max(0, ((currentTx - lastTx) / 1024)).toFixed(2))
        blkReadData.value.push(Math.max(0, ((currentRead - lastRead) / 1024)).toFixed(2))
        blkWriteData.value.push(Math.max(0, ((currentWrite - lastWrite) / 1024)).toFixed(2))
      }
      lastRx = currentRx
      lastTx = currentTx
      lastRead = currentRead
      lastWrite = currentWrite

      // 更新图表
      cpuChart?.setOption({ xAxis: { data: timestamps.value }, series: [{ data: cpuData.value }] })
      memChart?.setOption({ xAxis: { data: timestamps.value }, series: [{ data: memData.value }] })
      netChart?.setOption({ xAxis: { data: timestamps.value }, series: [{ data: netRxData.value }, { data: netTxData.value }] })
      blkChart?.setOption({ xAxis: { data: timestamps.value }, series: [{ data: blkReadData.value }, { data: blkWriteData.value }] })

    } catch (e) {
      console.error(e)
    }
  }

  eventSource.onerror = () => {
    // console.error('SSE Error')
    eventSource?.close()
  }
}

const stopMonitoring = () => {
  if (eventSource) {
    eventSource.close()
    eventSource = null
  }
  if (cpuChart) {
    cpuChart.dispose()
    cpuChart = null
  }
  if (memChart) {
    memChart.dispose()
    memChart = null
  }
  if (netChart) {
    netChart.dispose()
    netChart = null
  }
  if (blkChart) {
    blkChart.dispose()
    blkChart = null
  }
}

onUnmounted(() => {
  stopMonitoring()
})
</script>

<style scoped>
.monitor-container {
  padding: 10px;
}
.chart-box {
  height: 300px;
  width: 100%;
}
.mt-4 {
  margin-top: 16px;
}
</style>
