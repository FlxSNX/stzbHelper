<script setup>
import { ref, computed, onMounted } from 'vue'
import { NCard, NButton, NSpace, NTable, NStatistic, NGrid, NGi, NEmpty, useMessage } from 'naive-ui'
import { GetGroupWu } from '../../wailsjs/go/main/App'
import { RefreshCw } from 'lucide-vue-next'

const nmessage = useMessage()
const groupdata = ref([])
const loading = ref(false)
const sortKey = ref('total_wu')
const sortOrder = ref('desc')

const toggleSort = (key) => {
    if (sortKey.value === key) {
        sortOrder.value = sortOrder.value === 'asc' ? 'desc' : 'asc'
    } else {
        sortKey.value = key
        sortOrder.value = 'desc'
    }
}

const sortedData = computed(() => {
    if (!sortKey.value) return groupdata.value
    return [...groupdata.value].sort((a, b) => {
        let va = a[sortKey.value], vb = b[sortKey.value]
        if (typeof va === 'string') {
            return sortOrder.value === 'asc' ? va.localeCompare(vb) : vb.localeCompare(va)
        }
        return sortOrder.value === 'asc' ? va - vb : vb - va
    })
})

const totalMembers = computed(() => groupdata.value.reduce((sum, g) => sum + g.member_count, 0))
const totalWu = computed(() => groupdata.value.reduce((sum, g) => sum + g.total_wu, 0))
const avgWu = computed(() => groupdata.value.length > 0 ? Math.round(totalWu.value / groupdata.value.length) : 0)

function formatWu(val) {
    const n = Math.floor(val)
    if (n >= 10000) {
        return (n / 10000).toFixed(2) + '万'
    }
    return n
}

function getData() {
    loading.value = true
    groupdata.value = []
    GetGroupWu().then(v => {
        let resp = JSON.parse(v)
        if (resp.code == 200) {
            groupdata.value = resp.data
        } else {
            nmessage.error(resp.msg)
        }
    }).catch(() => {
        nmessage.error('获取分组武勋数据失败')
    }).finally(() => {
        loading.value = false
    })
}

onMounted(() => {
    getData()
})
</script>

<template>
    <div class="page-groupwu">
        <n-grid :cols="3" :x-gap="16" :y-gap="16" class="stat-grid">
            <n-gi>
                <n-card embedded size="small">
                    <n-statistic label="总人数" :value="totalMembers" />
                </n-card>
            </n-gi>
            <n-gi>
                <n-card embedded size="small">
                    <n-statistic label="总武勋" :value="formatWu(totalWu)" />
                </n-card>
            </n-gi>
            <n-gi>
                <n-card embedded size="small">
                    <n-statistic label="平均武勋" :value="formatWu(avgWu)" />
                </n-card>
            </n-gi>
        </n-grid>

        <n-card class="page-card" embedded>
            <div class="page-header">
                <div class="page-header-info">
                    <h2 class="page-title">分组武勋</h2>
                    <p class="page-desc">更新武勋数据请同步成员数据</p>
                </div>
                <n-space>
                    <n-button @click="getData" :loading="loading">
                        <template #icon><RefreshCw :size="16" /></template>
                        刷新
                    </n-button>
                </n-space>
            </div>

            <n-table v-if="groupdata.length > 0" :bordered="true" :single-line="false" class="styled-table">
                <thead>
                    <tr>
                        <th class="sortable-th" @click="toggleSort('group')">分组名称 <span class="sort-icon" v-if="sortKey==='group'">{{ sortOrder === 'asc' ? '↑' : '↓' }}</span></th>
                        <th class="sortable-th" @click="toggleSort('member_count')">人数 <span class="sort-icon" v-if="sortKey==='member_count'">{{ sortOrder === 'asc' ? '↑' : '↓' }}</span></th>
                        <th class="sortable-th" @click="toggleSort('total_wu')">总武勋 <span class="sort-icon" v-if="sortKey==='total_wu'">{{ sortOrder === 'asc' ? '↑' : '↓' }}</span></th>
                        <th class="sortable-th" @click="toggleSort('average_wu')">平均武勋 <span class="sort-icon" v-if="sortKey==='average_wu'">{{ sortOrder === 'asc' ? '↑' : '↓' }}</span></th>
                        <th class="sortable-th" @click="toggleSort('zero_wu_count')">0武勋人数 <span class="sort-icon" v-if="sortKey==='zero_wu_count'">{{ sortOrder === 'asc' ? '↑' : '↓' }}</span></th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="u in sortedData" :key="u.group">
                        <td>
                            <span class="group-name">{{ u.group }}</span>
                        </td>
                        <td>{{ u.member_count }}</td>
                        <td>{{ formatWu(u.total_wu) }}</td>
                        <td>{{ formatWu(u.average_wu) }}</td>
                        <td>
                            <span :class="{ 'zero-warn': u.zero_wu_count > 0 }">{{ u.zero_wu_count }}</span>
                        </td>
                    </tr>
                </tbody>
            </n-table>
            <n-empty v-else description="暂无分组武勋数据" style="padding: 60px 0;" />
        </n-card>
    </div>
</template>

<style scoped lang="scss">
.page-groupwu {
    display: flex;
    flex-direction: column;
    gap: 20px;
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

.styled-table {
    border-radius: 8px;
    overflow: hidden;

    thead {
        th {
            background: var(--color-bg);
            font-weight: 600;
            color: var(--color-text-secondary);
            font-size: 13px;
            padding: 12px 16px;
        }

        .sortable-th {
            cursor: pointer;
            user-select: none;
            transition: color 0.15s;

            &:hover {
                color: var(--color-text);
            }
        }

        .sort-icon {
            color: var(--color-primary);
            font-size: 12px;
        }
    }

    tbody {
        td {
            padding: 12px 16px;
            font-size: 14px;
        }

        tr:nth-child(even) {
            background: var(--color-bg);
        }

        tr:hover td {
            background: var(--color-surface-hover);
        }
    }
}

.group-name {
    font-weight: 600;
    color: var(--color-text);
}

.zero-warn {
    color: #ef4444;
    font-weight: 600;
}
</style>
