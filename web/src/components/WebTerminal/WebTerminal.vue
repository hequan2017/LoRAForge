<template>
  <div class="web-terminal-container" ref="containerRef">
    <!-- 终端工具栏 -->
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
        <p>点击下方按钮连接到容器终端</p>
        <el-button type="primary" size="large" @click="connect">
          <el-icon style="margin-right: 5px"><Link /></el-icon>
          连接终端
        </el-button>
      </div>
    </div>

    <!-- 快捷键提示 -->
    <div class="shortcut-hints" v-if="connected && showHints">
      <div class="hint-item">Ctrl+L 清屏</div>
      <div class="hint-item">Ctrl+C 中断</div>
      <div class="hint-item">双击选中文本</div>
      <el-button link type="primary" size="small" @click="showHints = false">隐藏</el-button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, nextTick, watch, computed } from 'vue'
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
    default: 'Terminal'
  },
  cols: {
    type: Number,
    default: 80
  },
  rows: {
    type: Number,
    default: 24
  },
  autoConnect: {
    type: Boolean,
    default: true
  },
  reconnect: {
    type: Boolean,
    default: true
  },
  reconnectInterval: {
    type: Number,
    default: 3000
  }
})

const emit = defineEmits(['connected', 'disconnected', 'error'])

const containerRef = ref(null)
const terminalRef = ref(null)
const connected = ref(false)
const connecting = ref(false)
const error = ref('')
const showRetry = ref(false)
const showHints = ref(true)
const ws = ref(null)
const fontSize = ref(14)
const reconnectTimer = ref(null)
let terminal = null
let fitAddon = null
let webLinksAddon = null
let searchAddon = null

// 初始化终端
const initTerminal = () => {
  if (terminal) return

  terminal = new Terminal({
    cursorBlink: true,
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
    cols: props.cols,
    rows: props.rows,
    scrollback: 1000,
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

      // 欢迎消息
      terminal.writeln('\x1b[1;32m欢迎使用 LoRAForge Web 终端\x1b[0m')
      terminal.writeln('提示: 双击可快速复制文本，Ctrl+L 清屏\r\n')
    }
  })

  // 终端输入处理
  terminal.onData(data => {
    if (ws.value && connected.value) {
      ws.value.send(data)
    }
  })

  // 终端尺寸变化
  terminal.onResize(({ cols, rows }) => {
    sendResize(cols, rows)
  })

  // 监听窗口大小变化
  window.addEventListener('resize', handleResize)

  // 键盘快捷键
  setupKeyboardShortcuts()
}

// 设置键盘快捷键
const setupKeyboardShortcuts = () => {
  const handleKeyDown = (e) => {
    // Ctrl+L 清屏
    if (e.ctrlKey && e.key === 'l') {
      e.preventDefault()
      if (terminal) {
        terminal.clear()
      }
    }
  }

  window.addEventListener('keydown', handleKeyDown)

  // 存储清理函数
  window._terminalKeyHandler = handleKeyDown
}

// 处理窗口大小变化
const handleResize = () => {
  if (fitAddon && terminalRef.value) {
    fitAddon.fit()
  }
}

// 发送终端尺寸调整
const sendResize = (cols, rows) => {
  if (!ws.value || !connected.value || !terminal) return

  try {
    ws.value.send(JSON.stringify({
      type: 'resize',
      cols: cols || terminal.cols,
      rows: rows || terminal.rows
    }))
  } catch (e) {
    console.error('发送尺寸调整失败:', e)
  }
}

// 改变字体大小
const changeFontSize = (delta) => {
  const newSize = fontSize.value + delta
  if (newSize >= 10 && newSize <= 24) {
    fontSize.value = newSize
    if (terminal) {
      terminal.options.fontSize = newSize
    }
  }
}

// 清屏
const clearTerminal = () => {
  if (terminal) {
    terminal.clear()
  }
}

// 强制重连
const forceReconnect = () => {
  disconnect()
  nextTick(() => {
    connect()
  })
}

// 连接 WebSocket
const connect = () => {
  if (connected.value || connecting.value) return

  connecting.value = true
  error.value = ''
  showRetry.value = false

  try {
    ws.value = new WebSocket(props.url)

    ws.value.onopen = () => {
      connected.value = true
      connecting.value = false
      error.value = ''
      showRetry.value = false

      // 清除重连定时器
      if (reconnectTimer.value) {
        clearTimeout(reconnectTimer.value)
        reconnectTimer.value = null
      }

      emit('connected')

      if (terminal) {
        terminal.clear()
        terminal.writeln('\x1b[1;32m✓ 已连接到容器 ' + props.title + '\x1b[0m')
        terminal.writeln('\x1b[2m提示: Ctrl+C 中断命令，Ctrl+L 清屏\x1b[0m\r\n')
      }

      // 发送初始尺寸
      nextTick(() => {
        if (terminal) {
          sendResize(terminal.cols, terminal.rows)
        }
      })
    }

    ws.value.onmessage = (event) => {
      if (terminal) {
        terminal.write(event.data)
      }
    }

    ws.value.onerror = (event) => {
      console.error('WebSocket 错误:', event)
      connecting.value = false
      showRetry.value = true
      emit('error', event)
    }

    ws.value.onclose = (event) => {
      const wasConnected = connected.value
      connected.value = false
      connecting.value = false
      ws.value = null

      emit('disconnected', event)

      if (terminal && wasConnected) {
        terminal.writeln('\r\n\x1b[1;31m✗ 连接已断开\x1b[0m')
        if (props.reconnect) {
          terminal.writeln('\x1b[33m正在尝试重连...\x1b[0m')
          scheduleReconnect()
        }
      } else {
        error.value = '连接失败，请检查容器状态或网络'
      }
    }
  } catch (e) {
    error.value = '连接失败: ' + e.message
    connecting.value = false
    showRetry.value = true
    emit('error', e)
  }
}

// 安排重连
const scheduleReconnect = () => {
  if (!props.reconnect || reconnectTimer.value) return

  reconnectTimer.value = setTimeout(() => {
    if (!connected.value) {
      console.log('尝试重新连接...')
      connect()
    }
  }, props.reconnectInterval)
}

// 断开连接
const disconnect = () => {
  // 清除重连定时器
  if (reconnectTimer.value) {
    clearTimeout(reconnectTimer.value)
    reconnectTimer.value = null
  }

  if (ws.value) {
    ws.value.close()
    ws.value = null
  }
  connected.value = false
  connecting.value = false
}

// 清理
const cleanup = () => {
  disconnect()

  // 移除键盘事件监听
  if (window._terminalKeyHandler) {
    window.removeEventListener('keydown', window._terminalKeyHandler)
    delete window._terminalKeyHandler
  }

  window.removeEventListener('resize', handleResize)

  if (terminal) {
    terminal.dispose()
    terminal = null
  }
}

// 监听 URL 变化
watch(() => props.url, () => {
  disconnect()
  if (props.url && props.autoConnect) {
    nextTick(() => connect())
  }
})

onMounted(() => {
  initTerminal()
  if (props.url && props.autoConnect) {
    connect()
  }
})

onUnmounted(() => {
  cleanup()
})

defineExpose({
  connect,
  disconnect,
  clearTerminal,
  forceReconnect,
  terminal
})
</script>

<style scoped>
.web-terminal-container {
  width: 100%;
  height: 100%;
  background-color: #1e1e1e;
  border-radius: 4px;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.terminal-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 6px 12px;
  background-color: #252526;
  border-bottom: 1px solid #3e3e3e;
  user-select: none;
}

.toolbar-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.toolbar-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.terminal-title {
  color: #cccccc;
  font-size: 13px;
  font-weight: 500;
}

.terminal-content {
  flex: 1;
  overflow: hidden;
  position: relative;
}

.error-message {
  padding: 12px;
  background-color: #1e1e1e;
}

.connect-prompt {
  flex: 1;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: #1e1e1e;
}

.connect-content {
  text-align: center;
  color: #888888;
}

.connect-icon {
  color: #4CAF50;
  margin-bottom: 16px;
  opacity: 0.8;
}

.connect-content p {
  margin-bottom: 20px;
  font-size: 14px;
}

.shortcut-hints {
  display: flex;
  gap: 16px;
  padding: 6px 12px;
  background-color: #252526;
  border-top: 1px solid #3e3e3e;
  font-size: 12px;
  color: #888888;
}

.hint-item {
  display: flex;
  align-items: center;
  gap: 4px;
}

:deep(.xterm) {
  height: 100%;
}

:deep(.xterm .xterm-viewport) {
  overflow-y: auto;
}

:deep(.xterm .xterm-screen) {
  padding: 4px;
}

/* 滚动条样式优化 */
:deep(.xterm .xterm-viewport::-webkit-scrollbar) {
  width: 8px;
}

:deep(.xterm .xterm-viewport::-webkit-scrollbar-track) {
  background: #1e1e1e;
}

:deep(.xterm .xterm-viewport::-webkit-scrollbar-thumb) {
  background: #4a4a4a;
  border-radius: 4px;
}

:deep(.xterm .xterm-viewport::-webkit-scrollbar-thumb:hover) {
  background: #5a5a5a;
}

/* 选中样式优化 */
:deep(.xterm .xterm-selection) {
  background-color: rgba(255, 255, 255, 0.2) !important;
}

/* 按钮样式 */
:deep(.el-button-group .el-button) {
  background-color: #3c3c3c;
  border-color: #4a4a4a;
  color: #cccccc;
}

:deep(.el-button-group .el-button:hover) {
  background-color: #4a4a4a;
  border-color: #5a5a5a;
  color: #ffffff;
}

:deep(.el-button-group .el-button:disabled) {
  background-color: #2c2c2c;
  border-color: #3a3a3a;
  color: #666666;
}
</style>
