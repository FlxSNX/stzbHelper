<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { NCard, NButton, NEmpty, NSpin, NModal, NInput, useMessage } from 'naive-ui'
import { GetDbList, SelectDb, CreateDb } from '../../wailsjs/go/main/App'
import { Database, ChevronRight, RefreshCw, Plus } from 'lucide-vue-next'

const router = useRouter()
const nmessage = useMessage()
const dbList = ref([])
const loading = ref(false)
const selecting = ref(false)
const selectedName = ref('')

const showCreateModal = ref(false)
const newDbName = ref('')
const creating = ref(false)

const loadDbList = () => {
    loading.value = true
    GetDbList().then(v => {
        let resp = JSON.parse(v)
        if (resp.code == 200) {
            dbList.value = resp.data || []
        } else {
            nmessage.error(resp.msg)
        }
    }).catch(e => {
        nmessage.error('获取数据库列表失败: ' + e)
    }).finally(() => {
        loading.value = false
    })
}

const handleSelect = (name) => {
    selecting.value = true
    selectedName.value = name
    SelectDb(name).then(v => {
        let resp = JSON.parse(v)
        if (resp.code == 200) {
            sessionStorage.setItem('selectedDb', name)
            router.push('/')
        } else {
            nmessage.error(resp.msg)
        }
    }).catch(e => {
        nmessage.error('连接数据库失败: ' + e)
    }).finally(() => {
        selecting.value = false
        selectedName.value = ''
    })
}

const handleCreate = () => {
    if (!newDbName.value.trim()) {
        nmessage.warning('请输入数据库名称')
        return
    }
    creating.value = true
    CreateDb(newDbName.value.trim()).then(v => {
        let resp = JSON.parse(v)
        if (resp.code == 200) {
            nmessage.success(resp.msg)
            sessionStorage.setItem('selectedDb', newDbName.value.trim())
            router.push('/')
        } else {
            nmessage.error(resp.msg)
        }
    }).catch(e => {
        nmessage.error('创建数据库失败: ' + e)
    }).finally(() => {
        creating.value = false
    })
}

onMounted(() => {
    loadDbList()
})
</script>

<template>
    <n-modal v-model:show="showCreateModal" preset="card" title="创建数据库" size="small"
        style="width: 400px" :bordered="false" :mask-closable="false" to="body">
        <n-input v-model:value="newDbName" placeholder="请输入数据库名称" @keyup.enter="handleCreate" />
        <template #footer>
            <div style="display: flex; justify-content: flex-end; gap: 8px;">
                <n-button strong secondary type="primary" :loading="creating" @click="handleCreate">
                    创建并进入
                </n-button>
                <n-button strong secondary @click="showCreateModal = false">
                    取消
                </n-button>
            </div>
        </template>
    </n-modal>

    <div class="select-db-page">
        <div class="select-db-container">
            <div class="select-db-header">
                <div class="select-db-icon">
                    <Database :size="28" />
                </div>
                <h1 class="select-db-title">选择数据库</h1>
                <p class="select-db-desc">请选择一个数据库文件以继续使用</p>
            </div>

            <div class="select-db-actions">
                <n-button size="small" @click="loadDbList" :loading="loading">
                    <template #icon><RefreshCw :size="14" /></template>
                    重新扫描
                </n-button>
                <n-button size="small" type="primary" @click="showCreateModal = true">
                    <template #icon><Plus :size="14" /></template>
                    创建数据库
                </n-button>
            </div>

            <n-card class="select-db-card" embedded :bordered="false">
                <div class="db-list-loading" v-if="loading">
                    <n-spin size="medium" />
                    <span>正在扫描数据库文件...</span>
                </div>

                <div class="db-list" v-else-if="dbList.length > 0">
                    <div class="db-item" v-for="db in dbList" :key="db"
                        :class="{ 'db-item--selecting': selecting && selectedName === db }"
                        @click="!selecting && handleSelect(db)">
                        <div class="db-item-left">
                            <div class="db-item-icon">
                                <Database :size="18" />
                            </div>
                            <span class="db-item-name">{{ db }}</span>
                        </div>
                        <div class="db-item-right">
                            <n-spin :size="14" v-if="selecting && selectedName === db" />
                            <ChevronRight :size="16" v-else />
                        </div>
                    </div>
                </div>

                <n-empty v-else description="未找到数据库文件" style="padding: 32px 0;" />
            </n-card>
        </div>
    </div>
</template>

<style scoped lang="scss">
.select-db-page {
    display: flex;
    align-items: flex-start;
    justify-content: center;
    min-height: 100%;
    padding-top: 80px;
    background: var(--color-bg);
}

.select-db-container {
    width: 100%;
    max-width: 480px;
}

.select-db-header {
    text-align: center;
    margin-bottom: 24px;
}

.select-db-icon {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 56px;
    height: 56px;
    border-radius: 14px;
    background: linear-gradient(135deg, #404040, #171717);
    color: #fff;
    margin-bottom: 16px;
}

.select-db-title {
    font-size: 22px;
    font-weight: 700;
    color: var(--color-text);
    margin: 0 0 6px;
}

.select-db-desc {
    font-size: 14px;
    color: var(--color-text-secondary);
    margin: 0;
}

.select-db-card {
    border-radius: 12px;
}

.db-list-loading {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 12px;
    padding: 40px 0;
    color: var(--color-text-secondary);
    font-size: 14px;
}

.db-list {
    display: flex;
    flex-direction: column;
    gap: 4px;
}

.db-item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 14px 16px;
    border-radius: 8px;
    cursor: pointer;
    transition: background-color 0.15s, transform 0.15s;

    &:hover {
        background: var(--color-primary-light);
    }

    &:active {
        transform: scale(0.99);
    }

    &--selecting {
        opacity: 0.6;
        pointer-events: none;
    }
}

.db-item-left {
    display: flex;
    align-items: center;
    gap: 12px;
}

.db-item-icon {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 36px;
    height: 36px;
    border-radius: 8px;
    background: var(--color-primary-light);
    color: var(--color-accent);
}

.db-item-name {
    font-size: 14px;
    font-weight: 500;
    color: var(--color-text);
}

.db-item-right {
    display: flex;
    align-items: center;
    color: var(--color-text-secondary);
}

.select-db-actions {
    display: flex;
    justify-content: center;
    gap: 12px;
    margin-bottom: 16px;
}
</style>
