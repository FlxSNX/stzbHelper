<script setup>
import { ref, onMounted } from 'vue'
import { NCard, NButton, NStatistic, NSpace, NGrid, NGi, NAlert, NSpin, useMessage } from 'naive-ui'
import { EnableGetBattleReport, DisableGetBattleReport, GetTaskList, GetTeamUser, CheckUpdate, GetVersion } from '../../wailsjs/go/main/App'
import { BrowserOpenURL } from '../../wailsjs/runtime/runtime'
import { RefreshCw, Download } from 'lucide-vue-next'

const nmessage = useMessage()

const taskCount = ref(0)
const memberCount = ref(0)
const version = ref('')
const updateInfo = ref(null)
const checkingUpdate = ref(false)

const onCheckUpdate = () => {
    checkingUpdate.value = true
    updateInfo.value = null
    CheckUpdate().then(v => {
        let resp = JSON.parse(v)
        if (resp.code == 200) {
            updateInfo.value = resp.data
        } else {
            nmessage.error(resp.msg)
        }
    }).catch(e => {
        nmessage.error('检查更新失败: ' + e)
    }).finally(() => {
        checkingUpdate.value = false
    })
}

const openUpdateUrl = (url) => {
    BrowserOpenURL(url)
}

const onEnableGetBattleReport = () => {
    EnableGetBattleReport().then(v => {
        let data = JSON.parse(v)
        if (data.code == 200) {
            nmessage.success('开启成功')
        } else {
            nmessage.error(data.msg)
        }
    }).catch(e => {
        nmessage.error('开启获取战报详情失败:' + e)
    })
}

const onDisableGetBattleReport = () => {
    DisableGetBattleReport().then(v => {
        let data = JSON.parse(v)
        if (data.code == 200) {
            nmessage.success('关闭成功')
        } else {
            nmessage.error(data.msg)
        }
    }).catch(e => {
        nmessage.error('关闭获取战报详情失败:' + e)
    })
}

onMounted(() => {
    GetVersion().then(v => {
        let resp = JSON.parse(v)
        if (resp.code == 200) {
            version.value = resp.data
        }
    }).catch(() => {})

    GetTaskList().then(v => {
        let resp = JSON.parse(v)
        if (resp.code == 200) {
            taskCount.value = resp.data.length
        }
    }).catch(() => {})

    GetTeamUser().then(v => {
        let data = JSON.parse(v)
        if (data.data) {
            memberCount.value = data.data.length
        }
    }).catch(() => {})
})
</script>

<template>
    <div class="page-index">
        <div class="page-hero">
            <div class="page-hero-content">
                <h1 class="page-hero-title">率土之滨助手</h1>
                <p class="page-hero-desc">stzbHelper &middot; Version {{ version }}</p>
            </div>
        </div>

        <n-grid :cols="3" :x-gap="16" :y-gap="16" class="stat-grid">
            <n-gi>
                <n-card embedded size="small">
                    <n-statistic label="同盟成员" :value="memberCount" />
                </n-card>
            </n-gi>
            <n-gi>
                <n-card embedded size="small">
                    <n-statistic label="攻城任务" :value="taskCount" />
                </n-card>
            </n-gi>
            <n-gi>
                <n-card embedded size="small">
                    <n-statistic label="应用版本" :value="version" />
                </n-card>
            </n-gi>
        </n-grid>

        <n-alert v-if="updateInfo && updateInfo.hasUpdate" type="success" :show-icon="true"
            style="border-radius: 12px;">
            <div style="margin-bottom: 8px;">
                <strong>发现新版本：{{ updateInfo.latestVer }}</strong>（当前版本：{{ updateInfo.currentVer }}）
            </div>
            <div v-if="updateInfo.name" style="margin-bottom: 4px;">{{ updateInfo.name }}</div>
            <div v-if="updateInfo.body" style="font-size: 13px; color: #64748b; white-space: pre-wrap;">{{ updateInfo.body }}</div>
            <n-button size="small" type="primary" style="margin-top: 8px;" @click="openUpdateUrl(updateInfo.url)">
                <template #icon><Download :size="14" /></template>
                前往下载
            </n-button>
        </n-alert>

        <n-card class="control-card" title="控制面板" embedded>
            <div class="control-section">
                <div class="control-item">
                    <div class="control-item-info">
                        <div class="control-item-title">获取详细战报</div>
                        <div class="control-item-desc">用于查询队伍功能拉取战报使用，开启时无法获取攻城战报</div>
                    </div>
                    <n-space>
                        <n-button type="primary" @click="onEnableGetBattleReport">开启</n-button>
                        <n-button @click="onDisableGetBattleReport">关闭</n-button>
                    </n-space>
                </div>
                <div class="control-item">
                    <div class="control-item-info">
                        <div class="control-item-title">检查更新</div>
                        <div class="control-item-desc">
                            <span v-if="updateInfo && !updateInfo.hasUpdate">当前已是最新版本 ({{ updateInfo.currentVer }})</span>
                            <span v-else>检查是否有新版本可用</span>
                        </div>
                    </div>
                    <n-button @click="onCheckUpdate" :loading="checkingUpdate">
                        <template #icon><RefreshCw :size="14" /></template>
                        检查更新
                    </n-button>
                </div>
            </div>
        </n-card>
    </div>
</template>

<style scoped lang="scss">
.page-index {
    display: flex;
    flex-direction: column;
    gap: 20px;
}

.page-hero {
    background: var(--color-hero-bg);
    border-radius: 12px;
    padding: 32px;
    color: var(--color-hero-text);

    &-title {
        font-size: 24px;
        font-weight: 700;
        margin-bottom: 8px;
    }

    &-desc {
        font-size: 14px;
        opacity: 0.85;
    }
}

.stat-grid {
    margin-top: 0;
}

.control-card {
    border-radius: 12px;
}

.control-section {
    display: flex;
    flex-direction: column;
    gap: 16px;
}

.control-item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 16px;
    background: var(--color-surface-hover);
    border-radius: 8px;

    &-info {
        flex: 1;
    }

    &-title {
        font-size: 15px;
        font-weight: 600;
        color: var(--color-text);
        margin-bottom: 4px;
    }

    &-desc {
        font-size: 13px;
        color: var(--color-text-secondary);
    }
}
</style>
