<script setup>
import { ref, computed, onMounted } from 'vue'
import { NCard, NButton, NSpace, NTag, NInput, NEmpty, NGrid, NGi, useMessage, useDialog } from 'naive-ui'
import { GetTeamUser } from '../../wailsjs/go/main/App'
import { formatTimestamp, splitwid } from '@/utils/format'
import * as XLSX from 'xlsx'
import { Search, RefreshCw, Download, UserPlus } from 'lucide-vue-next'

const dialog = useDialog()
const nmessage = useMessage()
const teamUsers = ref([])
const usersNum = ref(0)
const searchText = ref('')
const loading = ref(false)

const filteredUsers = computed(() => {
    if (!searchText.value) return teamUsers.value
    const keyword = searchText.value.toLowerCase()
    return teamUsers.value.filter(u =>
        u.name.toLowerCase().includes(keyword) ||
        u.group.toLowerCase().includes(keyword)
    )
})

const syncuser = () => {
    dialog.info({
        title: '信息',
        content: '请前往游戏中，点开同盟成员列表即可同步',
        positiveText: '确认',
        transformOrigin: 'center',
    })
}

function getUserList() {
    loading.value = true
    teamUsers.value = []
    usersNum.value = 0
    GetTeamUser().then(v => {
        let data = JSON.parse(v)
        teamUsers.value = data.data
        usersNum.value = data.data.length
    }).catch(() => {}).finally(() => {
        loading.value = false
    })
}

const exportExcel = () => {
    let data = []
    data.push(['名字', '分组', '势力', '本周武勋', '总贡献', '周贡献', '位置', '进盟时间'])
    Object.values(teamUsers.value).forEach(v => {
        data.push([
            v.name, v.group, v.power, v.wu,
            v.contribute_total, v.contribute_week,
            splitwid(v.pos), formatTimestamp(v.join_time),
        ])
    })
    const ws = XLSX.utils.aoa_to_sheet(data)
    const wb = XLSX.utils.book_new()
    XLSX.utils.book_append_sheet(wb, ws, 'Sheet1')
    XLSX.writeFile(wb, `${formatTimestamp(parseInt(new Date().getTime() / 1000))}同盟成员表.xlsx`)
}

onMounted(() => {
    getUserList()
})
</script>

<template>
    <div class="page-teamuser">
        <n-card class="page-card" embedded>
            <div class="page-header">
                <div class="page-header-info">
                    <h2 class="page-title">同盟成员</h2>
                    <p class="page-desc">成员数量：{{ usersNum }}</p>
                </div>
                <n-space>
                    <n-button @click="getUserList" :loading="loading">
                        <template #icon><RefreshCw :size="16" /></template>
                        刷新
                    </n-button>
                    <n-button @click="syncuser">
                        <template #icon><UserPlus :size="16" /></template>
                        同步成员
                    </n-button>
                    <n-button type="primary" @click="exportExcel">
                        <template #icon><Download :size="16" /></template>
                        导出表格
                    </n-button>
                </n-space>
            </div>

            <n-input
                v-model:value="searchText"
                placeholder="搜索成员名字或分组..."
                clearable
                class="search-input"
            >
                <template #prefix>
                    <Search :size="16" />
                </template>
            </n-input>

            <div class="member-list" v-if="filteredUsers.length > 0">
                <div class="member-card" v-for="user in filteredUsers" :key="user.id">
                    <div class="member-header">
                        <span class="member-name">{{ user.name }}</span>
                        <n-tag :bordered="false" type="info" size="small">{{ user.group }}</n-tag>
                    </div>
                    <n-grid :cols="2" :x-gap="12" :y-gap="8" class="member-stats">
                        <n-gi>
                            <div class="stat-item">
                                <span class="stat-label">ID</span>
                                <span class="stat-value">{{ user.id }}</span>
                            </div>
                        </n-gi>
                        <n-gi>
                            <div class="stat-item">
                                <span class="stat-label">势力</span>
                                <span class="stat-value">{{ user.power }}</span>
                            </div>
                        </n-gi>
                        <n-gi>
                            <div class="stat-item">
                                <span class="stat-label">周武勋</span>
                                <span class="stat-value highlight">{{ user.wu }}</span>
                            </div>
                        </n-gi>
                        <n-gi>
                            <div class="stat-item">
                                <span class="stat-label">总贡献</span>
                                <span class="stat-value">{{ user.contribute_total }}</span>
                            </div>
                        </n-gi>
                        <n-gi>
                            <div class="stat-item">
                                <span class="stat-label">周贡献</span>
                                <span class="stat-value">{{ user.contribute_week }}</span>
                            </div>
                        </n-gi>
                        <n-gi>
                            <div class="stat-item">
                                <span class="stat-label">位置</span>
                                <span class="stat-value">{{ splitwid(user.pos) }}</span>
                            </div>
                        </n-gi>
                        <n-gi :span="2">
                            <div class="stat-item">
                                <span class="stat-label">进盟时间</span>
                                <span class="stat-value">{{ formatTimestamp(user.join_time) }}</span>
                            </div>
                        </n-gi>
                    </n-grid>
                </div>
            </div>
            <n-empty v-else description="暂无成员数据" style="padding: 60px 0;" />
        </n-card>
    </div>
</template>

<style scoped lang="scss">
.page-teamuser {
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
    margin-bottom: 20px;
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

.search-input {
    margin-bottom: 20px;
}

.member-list {
    display: flex;
    flex-direction: column;
    gap: 12px;
}

.member-card {
    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: 10px;
    padding: 16px 20px;
    transition: box-shadow 0.2s, transform 0.2s;

    &:hover {
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
        transform: translateY(-1px);
    }
}

.member-header {
    display: flex;
    align-items: center;
    gap: 10px;
    margin-bottom: 12px;
}

.member-name {
    font-size: 16px;
    font-weight: 600;
    color: var(--color-text);
}

.member-stats {
    margin-top: 4px;
}

.stat-item {
    display: flex;
    flex-direction: column;
    gap: 2px;
}

.stat-label {
    font-size: 12px;
    color: var(--color-text-secondary);
}

.stat-value {
    font-size: 14px;
    color: var(--color-text);
    font-weight: 500;

    &.highlight {
        color: var(--color-accent);
    }
}
</style>
