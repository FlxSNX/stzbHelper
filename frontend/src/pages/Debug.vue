<script setup>
import { ref } from 'vue'
import { NCard, NButton, NInput, NTag, NCollapse, NCollapseItem, NSpace } from 'naive-ui'
import * as App from '../../wailsjs/go/main/App'
import { Play, Trash2, Copy } from 'lucide-vue-next'
import { useMessage } from 'naive-ui'

const nmessage = useMessage()
const results = ref({})

const methods = ref(Object.keys(App).map(name => {
    const fn = App[name]
    const paramCount = fn.length
    const params = []
    for (let i = 0; i < paramCount; i++) {
        params.push({ name: `arg${i + 1}`, value: '' })
    }
    return { name, fn, params, paramCount }
}))

const callMethod = async (method) => {
    const args = method.params.map(p => {
        const v = p.value.trim()
        if (v === '') return v
        try { return JSON.parse(v) } catch { return v }
    })
    const start = performance.now()
    try {
        const raw = await method.fn(...args)
        const elapsed = (performance.now() - start).toFixed(1)
        let parsed
        try { parsed = JSON.parse(raw) } catch { parsed = raw }
        results.value[method.name] = {
            status: 'success',
            time: elapsed,
            raw: raw,
            parsed: parsed
        }
    } catch (e) {
        const elapsed = (performance.now() - start).toFixed(1)
        results.value[method.name] = {
            status: 'error',
            time: elapsed,
            raw: String(e),
            parsed: null
        }
    }
}

const clearResults = () => {
    results.value = {}
}

const copyResult = (name) => {
    const r = results.value[name]
    if (r) {
        navigator.clipboard.writeText(r.raw)
        nmessage.success('已复制')
    }
}
</script>

<template>
    <div class="page-debug">
        <n-card class="page-card" embedded>
            <div class="page-header">
                <div class="page-header-info">
                    <h2 class="page-title">API 调试</h2>
                    <p class="page-desc">共 {{ methods.length }} 个方法</p>
                </div>
                <n-button size="small" @click="clearResults">
                    <template #icon><Trash2 :size="14" /></template>
                    清空结果
                </n-button>
            </div>

            <n-collapse>
                <n-collapse-item v-for="m in methods" :key="m.name" :name="m.name">
                    <template #header>
                        <div class="method-header">
                            <span class="method-name">{{ m.name }}</span>
                            <n-tag :bordered="false" size="small" type="info">{{ m.paramCount }} 参数</n-tag>
                        </div>
                    </template>

                    <div class="method-body">
                        <div class="method-params" v-if="m.params.length > 0">
                            <div class="param-row" v-for="(p, i) in m.params" :key="i">
                                <span class="param-label">{{ p.name }}</span>
                                <n-input v-model:value="p.value" size="small" placeholder="参数值（JSON 或字符串）" />
                            </div>
                        </div>

                        <n-button size="small" type="primary" @click="callMethod(m)" style="margin-top: 8px;">
                            <template #icon><Play :size="14" /></template>
                            调用
                        </n-button>

                        <div class="method-result" v-if="results[m.name]">
                            <div class="result-header">
                                <n-tag :type="results[m.name].status === 'success' ? 'success' : 'error'" :bordered="false"
                                    size="small">
                                    {{ results[m.name].status === 'success' ? '成功' : '失败' }}
                                </n-tag>
                                <span class="result-time">{{ results[m.name].time }}ms</span>
                                <n-button text size="small" @click="copyResult(m.name)">
                                    <template #icon><Copy :size="14" /></template>
                                </n-button>
                            </div>
                            <pre class="result-json">{{ JSON.stringify(results[m.name].parsed, null, 2) }}</pre>
                        </div>
                    </div>
                </n-collapse-item>
            </n-collapse>
        </n-card>
    </div>
</template>

<style scoped lang="scss">
.page-debug {
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

.method-header {
    display: flex;
    align-items: center;
    gap: 8px;
}

.method-name {
    font-family: 'Cascadia Code', 'Fira Code', 'Consolas', monospace;
    font-size: 13px;
    font-weight: 600;
    color: var(--color-text);
}

.method-body {
    padding: 4px 0;
}

.method-params {
    display: flex;
    flex-direction: column;
    gap: 8px;
    margin-bottom: 8px;
}

.param-row {
    display: flex;
    align-items: center;
    gap: 8px;
}

.param-label {
    font-family: 'Cascadia Code', 'Fira Code', 'Consolas', monospace;
    font-size: 12px;
    color: var(--color-text-secondary);
    min-width: 50px;
}

.method-result {
    margin-top: 12px;
    background: var(--color-bg);
    border-radius: 8px;
    padding: 12px;
}

.result-header {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 8px;
}

.result-time {
    font-size: 12px;
    color: var(--color-text-secondary);
}

.result-json {
    font-family: 'Cascadia Code', 'Fira Code', 'Consolas', monospace;
    font-size: 12px;
    line-height: 1.5;
    color: var(--color-text);
    margin: 0;
    white-space: pre-wrap;
    word-break: break-all;
    max-height: 300px;
    overflow-y: auto;
}
</style>
