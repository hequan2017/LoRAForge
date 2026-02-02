<template>
  <div class="web-log-container" ref="containerRef">
    <!-- 日志工具栏 -->
    <div class="terminal-toolbar">
      <div class="toolbar-left">
        <span class="terminal-title">{{ title }}</span>
        <el-tag v-if="connected" type="success" size="small">已连接</el-tag>
        <el-tag v-else-if="connecting" type="warning" size="small">连接中...</el-tag>
        <el-tag v-else type="info" size="small">未连接</el-tag>
      </div>
      <div class="toolbar-right">
        <!-- 字体大小调整 -->
        <el-button-group size="small">
          <el-button :disabled="!connected" @click="changeFontSize(-1)" title="减小字体">
            <el-icon><ZoomOut /></el-icon>
          </el-button>
          <el-button disabled style="width: 50px">{{ fontSize }}px</el-button>
          <el-button :disabled="!connected" @click="changeFontSize(1)" title="增大字体">
            <el-icon><ZoomIn /></el-icon>
          </el-button>
        </el-button-group>

        <!-- 功能按钮 -->
        <el-button-group size="small" style="margin-left: 8px">
          <el-tooltip content="清屏 (Ctrl+L)" placement="bottom">
            <el-button :disabled="!connected" @click="clearTerminal">
              <el-icon><Delete /></el-icon>
            </el-button>
          </el-tooltip>
          <el-tooltip content="强制重连" placement="bottom">
            <el-button :disabled="connected || connecting" @click="forceReconnect" type="warning">
              <el-icon><Refresh /></el-icon>
            </el-button>
          </el-tooltip>
          <el-tooltip content="断开连接" placement="bottom">
            <el-button :disabled="!connected" @click="disconnect" type="danger">
              <el-icon><SwitchButton /></el-icon>
            </el-button>
          </el-tooltip>
        </el-button-group>
      </div>
    </div>

    <!-- 错误提示 -->
    <div v-if="error" class="error-message">
      <el-alert :title="error" type="error" :closable="false">
        <template #default>
          <div style="display: flex; align-items: center; justify-content: space-between;">
            <span>{{ error }}</span>
            <el-button v-if="showRetry" type="primary" size="small" @click="connect">重试</el-button>
          </div>
        </template>
      </el-alert>
    </div>

    <!-- 终端内容区 -->
    <div v-show="connected || !error" ref="terminalRef" class="terminal-content"></div>

    <!-- 连接提示 -->
    <div v-if="!connected && !error && !connecting" class="connect-prompt">
      <div class="connect-content">
        <el-icon class="connect-icon" :size="48"><Monitor /></el-icon>
        <p>点击下方按钮连接到容器日志</p>
        <el-button type="primary" size="large" @click="connect">
          <el-icon style="margin-right: 5px"><Link /></el-icon>
          连接日志
        </el-button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, nextTick, watch } from 'vue'
import { Terminal } from '@xterm/xterm'
import { FitAddon } from '@xterm/addon-fit'
import { WebLinksAddon } from '@xterm/addon-web-links'
import { SearchAddon } from '@xterm/addon-search'
import '@xterm/xterm/css/xterm.css'
import { ElMessage } from 'element-plus'

const props = defineProps({
  url: {
    type: String,
    required: true
  },
  title: {
    type: String,
    default: 'Container Logs'
  },
  autoConnect: {
    type: Boolean,
    default: true
  }
})

const emit = defineEmits(['connected', 'disconnected', 'error'])

const containerRef = ref(null)
const terminalRef = ref(null)
const connected = ref(false)
const connecting = ref(false)
const error = ref('')
const showRetry = ref(false)
const ws = ref(null)
const fontSize = ref(14)
let terminal = null
let fitAddon = null
let webLinksAddon = null
let searchAddon = null

// 初始化终端
const initTerminal = () => {
  if (terminal) return

  terminal = new Terminal({
    cursorBlink: false, // 日志不需要光标闪烁
    disableStdin: true, // 禁止输入
    fontSize: fontSize.value,
    lineHeight: 1.2,
    fontFamily: 'Consolas, "Courier New", "SF Mono", "Menlo", "Monaco", monospace',
    theme: {
      background: '#1e1e1e',
      foreground: '#d4d4d4',
      cursor: '#ffffff',
      cursorAccent: '#ffffff',
      selectionBackground: 'rgba(255, 255, 255, 0.3)',
      black: '#000000',
      red: '#cd3131',
      green: '#0dbc79',
      yellow: '#e5e510',
      blue: '#2472c8',
      magenta: '#bc3fbc',
      cyan: '#11a8cd',
      white: '#e5e5e5',
      brightBlack: '#666666',
      brightRed: '#f14c4c',
      brightGreen: '#23d18b',
      brightYellow: '#f5f543',
      brightBlue: '#3b8eea',
      brightMagenta: '#d670d6',
      brightCyan: '#29b8db',
      brightWhite: '#ffffff'
    },
    scrollback: 5000, // 增加回滚行数
    convertEol: true
  })

  fitAddon = new FitAddon()
  webLinksAddon = new WebLinksAddon()
  searchAddon = new SearchAddon()

  terminal.loadAddon(fitAddon)
  terminal.loadAddon(webLinksAddon)
  terminal.loadAddon(searchAddon)

  nextTick(() => {
    if (terminalRef.value) {
      terminal.open(terminalRef.value)
      fitAddon.fit()
      
      // 监听 resize
      window.addEventListener('resize', fitTerminal)
    }
  })
}

const fitTerminal = () => {
  if (fitAddon) {
    fitAddon.fit()
  }
}

// 建立 WebSocket 连接
const connect = () => {
  if (connected.value || connecting.value) return
  if (!props.url) {
    error.value = '连接地址为空'
    return
  }

  connecting.value = true
  error.value = ''
  showRetry.value = false

  initTerminal()

  try {
    ws.value = new WebSocket(props.url)
    ws.value.binaryType = 'arraybuffer'

    ws.value.onopen = () => {
      connected.value = true
      connecting.value = false
      error.value = ''
      terminal.writeln('\x1b[1;32m连接成功\x1b[0m\r\n')
      emit('connected')
      fitTerminal()
    }

    ws.value.onmessage = (event) => {
      // 接收数据写入终端
      if (typeof event.data === 'string') {
        terminal.write(event.data)
      } else {
        terminal.write(new Uint8Array(event.data))
      }
    }

    ws.value.onerror = (e) => {
      console.error('WebSocket error:', e)
      if (!connected.value) {
        error.value = '连接失败，请检查网络或后端服务'
        showRetry.value = true
      }
      connecting.value = false
    }

    ws.value.onclose = (e) => {
      connected.value = false
      connecting.value = false
      if (e.code !== 1000) { // 非正常关闭
        terminal.writeln(`\r\n\x1b[1;31m连接断开 (Code: ${e.code})\x1b[0m\r\n`)
        if (!error.value) {
           // error.value = `连接断开 (Code: ${e.code})`
           // showRetry.value = true
        }
      } else {
        terminal.writeln('\r\n\x1b[1;33m连接已关闭\x1b[0m\r\n')
      }
      emit('disconnected')
    }

  } catch (e) {
    error.value = `连接异常: ${e.message}`
    connecting.value = false
    showRetry.value = true
  }
}

const disconnect = () => {
  if (ws.value) {
    ws.value.close()
    ws.value = null
  }
  connected.value = false
  connecting.value = false
}

const forceReconnect = () => {
  disconnect()
  setTimeout(() => {
    connect()
  }, 500)
}

const clearTerminal = () => {
  if (terminal) {
    terminal.clear()
  }
}

const changeFontSize = (delta) => {
  fontSize.value += delta
  if (fontSize.value < 10) fontSize.value = 10
  if (fontSize.value > 30) fontSize.value = 30
  
  if (terminal) {
    terminal.options.fontSize = fontSize.value
    fitTerminal()
  }
}

// 监听 URL 变化
watch(() => props.url, (newUrl) => {
  if (newUrl && props.autoConnect) {
    disconnect()
    connect()
  }
})

onMounted(() => {
  if (props.autoConnect && props.url) {
    connect()
  }
})

onUnmounted(() => {
  disconnect()
  if (terminal) {
    terminal.dispose()
  }
  window.removeEventListener('resize', fitTerminal)
})

defineExpose({
  connect,
  disconnect,
  clearTerminal
})
</script>

<style scoped>
.web-log-container {
  height: 100%;
  display: flex;
  flex-direction: column;
  background-color: #1e1e1e;
  overflow: hidden;
}

.terminal-toolbar {
  height: 40px;
  background-color: #252526;
  border-bottom: 1px solid #333;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 10px;
  flex-shrink: 0;
}

.toolbar-left {
  display: flex;
  align-items: center;
  gap: 10px;
}

.terminal-title {
  color: #cccccc;
  font-weight: bold;
  font-size: 14px;
}

.toolbar-right {
  display: flex;
  align-items: center;
}

.terminal-content {
  flex: 1;
  overflow: hidden;
  padding: 5px;
  /* xterm 容器需要相对定位 */
  position: relative; 
}

.error-message {
  padding: 10px;
}

.connect-prompt {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #cccccc;
}

.connect-content {
  text-align: center;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 20px;
}

:deep(.xterm-viewport) {
  overflow-y: auto;
}
</style>
