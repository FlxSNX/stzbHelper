<script setup>
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import { NCard, NButton, NEmpty } from 'naive-ui'
import { GetLogs } from '../../wailsjs/go/main/App'
import { EventsOn, EventsOff } from '../../wailsjs/runtime/runtime'
import { Trash2, Terminal } from 'lucide-vue-next'

const logs = ref([])
const logContainer = ref(null)
let unsubscribe = null

const scrollToBottom = () => {
    nextTick(() => {
        if (logContainer.value) {
            logContainer.value.scrollTop = logContainer.value.scrollHeight
        }
    })
}

const loadLogs = () => {
    GetLogs().then(v => {
        let resp = JSON.parse(v)
        if (resp.code == 200) {
            logs.value = resp.data || []
            scrollToBottom()
        }
    })
}

const clearLogs = () => {
    logs.value = []
}

onMounted(() => {
    loadLogs()
    unsubscribe = EventsOn('log', (msg) => {
        logs.value.push(msg)
        scrollToBottom()
    })
})

onUnmounted(() => {
    if (unsubscribe) {
        EventsOff('log')
    }
})
</script>

<template>
    <div class="page-logs">
        <n-card class="page-card" embedded>
            <div class="page-header">
                <div class="page-header-info">
                    <h2 class="page-title">运行日志</h2>
                    <p class="page-desc">共 {{ logs.length }} 条日志</p>
                </div>
                <n-button size="small" @click="clearLogs">
                    <template #icon><Trash2 :size="14" /></template>
                    清空
                </n-button>
            </div>

            <div class="log-container" ref="logContainer" v-if="logs.length > 0">
                <div class="log-line" v-for="(line, i) in logs" :key="i">
                    <Terminal :size="12" class="log-icon" />
                    <span>{{ line }}</span>
                </div>
            </div>
            <n-empty v-else description="暂无日志" style="padding: 60px 0;" />
        </n-card>
    </div>
</template>

<style scoped lang="scss">
.page-logs {
    display: flex;
    flex-direction: column;
}

.page-card {
    border-radius: 12px;
}

.page-header {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    margin-bottom: 16px;
}

.page-title {
    font-size: 20px;
    font-weight: 600;
    color: var(--color-text);
    margin-bottom: 4px;
}

.page-desc {
    font-size: 13px;
    color: var(--color-text-secondary);
}

.log-container {
    max-height: calc(100vh - 200px);
    overflow-y: auto;
    background: #1e293b;
    border-radius: 8px;
    padding: 16px;
    font-family: 'Cascadia Code', 'Fira Code', 'Consolas', monospace;
    font-size: 13px;
    line-height: 1.7;
}

.log-line {
    display: flex;
    align-items: flex-start;
    gap: 8px;
    color: #e2e8f0;
    word-break: break-all;
}

.log-icon {
    flex-shrink: 0;
    color: var(--color-text-secondary);
    margin-top: 4px;
}
</style>
