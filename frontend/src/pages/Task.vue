<script setup>
import { ref, onMounted, computed, watch } from 'vue'
import {
    NCard, NButton, NSpace, NTag, NEmpty,
    NInput, NFormItem, NSelect, NDatePicker, NPopconfirm, NModal,
    NDataTable, NStatistic, NSpin,
    useMessage, useDialog
} from 'naive-ui'
import {
    GetTeamGroup, CreateTask, GetTaskList, DeleteTask,
    EnableGetReport, DisableGetReport, GetReportNumByTaskId, StatisticsReport,
    GetTask, DeleteTaskReport
} from '../../wailsjs/go/main/App'
import { formatTimestampMs, splitwid } from '@/utils/format'
import * as XLSX from 'xlsx'
import { Plus, RefreshCw, Eye, Play, Trash2, Eraser } from 'lucide-vue-next'

const nmessage = useMessage()
const addtaskshow = ref(false)
const targetgroup = ref([])
const grouplist = ref([])
const tasktime = ref(new Date().getTime())
const taskname = ref('')
const taskpos = ref()
const createing = ref(false)
const tasks = ref([])
const taskNum = ref(0)
const loading = ref(false)

const createTask = () => {
    createing.value = true
    CreateTask(taskname.value, tasktime.value, targetgroup.value, taskpos.value).then(v => {
        let resp = JSON.parse(v)
        if (resp.code == 200) {
            nmessage.success(resp.msg)
            taskname.value = ''
            targetgroup.value = []
            taskpos.value = []
            getTaskList()
        } else {
            nmessage.error(resp.msg)
        }
        createing.value = false
    }).catch(e => {
        createing.value = false
        nmessage.error(e)
    })
}

const delTask = (id) => {
    DeleteTask(id).then(v => {
        let resp = JSON.parse(v)
        if (resp.code == 200) {
            nmessage.success(resp.msg)
            getTaskList()
        } else {
            nmessage.error(resp.msg)
        }
    })
}

const delTaskReport = (id) => {
    DeleteTaskReport(id).then(v => {
        let resp = JSON.parse(v)
        if (resp.code == 200) {
            nmessage.success(resp.msg)
            getTaskList()
        } else {
            nmessage.error(resp.msg)
        }
    })
}

function getTaskList() {
    loading.value = true
    tasks.value = []
    taskNum.value = 0
    GetTaskList().then(v => {
        let resp = JSON.parse(v)
        if (resp.code == 200) {
            tasks.value = resp.data
            taskNum.value = resp.data.length
        } else {
            nmessage.error(resp.msg)
        }
    }).finally(() => {
        loading.value = false
    })
}

onMounted(() => {
    GetTeamGroup().then(v => {
        let resp = JSON.parse(v)
        if (resp.code == 200) {
            let data = resp.data
            grouplist.value = []
            data.forEach(e => {
                grouplist.value.push({ label: e, value: e })
            })
        }
    })
    getTaskList()
})

const showModal = ref(false)
const getReporting = ref(false)
const reportNum = ref(0)
const getReportNumTimer = ref(null)
const inStatistics = ref(false)
const curtaskid = ref(0)

const enableGetReport = (id, pos) => {
    showModal.value = true
    EnableGetReport(pos)
    getReporting.value = true
    reportNum.value = 0
    curtaskid.value = id
    getReportNumTimer.value = setInterval(() => {
        GetReportNumByTaskId(id).then(v => {
            let resp = JSON.parse(v)
            if (resp.code == 200) {
                reportNum.value = resp.data.count
            }
        })
    }, 1000)
}

const stopReport = () => {
    clearInterval(getReportNumTimer.value)
    getReportNumTimer.value = null
    getReporting.value = false
    inStatistics.value = false
    DisableGetReport()
}

watch(showModal, (val) => {
    if (!val) stopReport()
})

const statistics = () => {
    clearInterval(getReportNumTimer.value)
    getReportNumTimer.value = null
    getReporting.value = false
    inStatistics.value = true
    StatisticsReport(curtaskid.value).then(v => {
        let resp = JSON.parse(v)
        if (resp.code == 200) {
            nmessage.success(resp.msg)
            curtaskid.value = 0
            getTaskList()
        } else {
            nmessage.error(resp.msg)
        }
        inStatistics.value = false
        showModal.value = false
    }).catch(e => {
        inStatistics.value = false
        nmessage.error('统计考勤数据失败:' + e)
    })
}

const showModal2 = ref(false)
const taskDetail = ref({})
const getTaskDetail = (id) => {
    taskDetail.value = {}
    showModal2.value = true
    GetTask(id).then(v => {
        let resp = JSON.parse(v)
        if (resp.code == 200) {
            taskDetail.value = resp.data
        } else {
            nmessage.error(resp.msg)
        }
    }).catch(e => {
        nmessage.error('获取考勤数据失败:' + e)
    })
}

const exportExcel = () => {
    let data = []
    data.push(['名字', '分组', '主力', '拆迁', '主力次数', '拆迁次数'])
    Object.values(taskDetail.value.user_list).forEach(v => {
        data.push([v.name, v.group, v.atk_team_num, v.dis_team_num, v.atk_num, v.dis_num])
    })
    const ws = XLSX.utils.aoa_to_sheet(data)
    const wb = XLSX.utils.book_new()
    XLSX.utils.book_append_sheet(wb, ws, 'Sheet1')
    XLSX.writeFile(wb, `${taskDetail.value.name}考勤表.xlsx`)
}

const detailColumns = [
    { title: '名称', key: 'name', sorter: 'default', defaultSortOrder: false },
    { title: '分组', key: 'group', sorter: 'default', defaultSortOrder: false },
    { title: '主力', key: 'atk_team_num', sorter: (a, b) => a.atk_team_num - b.atk_team_num, defaultSortOrder: false },
    { title: '拆迁', key: 'dis_team_num', sorter: (a, b) => a.dis_team_num - b.dis_team_num, defaultSortOrder: false },
    { title: '主力次数', key: 'atk_num', sorter: (a, b) => a.atk_num - b.atk_num, defaultSortOrder: 'descend' },
    { title: '拆迁次数', key: 'dis_num', sorter: (a, b) => a.dis_num - b.dis_num, defaultSortOrder: false },
]

const detailData = computed(() => {
    if (!taskDetail.value.user_list) return []
    return Object.values(taskDetail.value.user_list)
})
</script>

<template>
    <n-modal v-model:show="addtaskshow" preset="card" title="新增任务" size="huge" :bordered="false"
        style="width: 520px" :mask-closable="false" to="body">
        <div class="modal-form">
            <n-form-item label="任务名称">
                <n-input v-model:value="taskname" placeholder="例如：内黄LV5 或者你也可以随意填写" />
            </n-form-item>
            <n-form-item label="任务坐标">
                <n-input pair separator="，" :placeholder="['X坐标', 'Y坐标']" v-model:value="taskpos" clearable />
            </n-form-item>
            <n-form-item label="任务时间">
                <n-date-picker v-model:value="tasktime" type="datetime" style="width: 100%;" />
            </n-form-item>
            <n-form-item label="目标分组">
                <n-select v-model:value="targetgroup" multiple :options="grouplist" placeholder="请选择分组" />
            </n-form-item>
        </div>
        <template #footer>
            <n-space justify="end">
                <n-button strong secondary type="primary" :loading="createing" @click="createTask">
                    添加
                </n-button>
                <n-button strong secondary type="error" @click="addtaskshow = false">
                    关闭
                </n-button>
            </n-space>
        </template>
    </n-modal>

    <n-modal v-model:show="showModal" preset="card" title="攻城考勤" size="huge" :bordered="false"
        style="width: 600px" :mask-closable="false" to="body">
        <div class="report-modal">
            <p class="report-tip">请前往游戏中，到攻城任务坐标位置查看同盟战报，并勾选守城军士（否则获取不了拆迁战报）。然后一直往下滑直到没有战报为止</p>
            <div class="report-counter">
                <n-statistic label="已获取战报" :value="reportNum">
                    <template #suffix>
                        <span style="font-size: 14px; color: #64748b;">封</span>
                    </template>
                </n-statistic>
            </div>
        </div>
        <template #footer>
            <n-space>
                <n-button strong secondary type="info" :loading="true" v-if="getReporting">
                    获取战报中
                </n-button>
                <n-button strong secondary type="success" @click="statistics" :loading="inStatistics">
                    {{ inStatistics ? '统计考勤数据中' : '已获取完战报，开始统计考勤数据' }}
                </n-button>
                <n-button strong secondary @click="showModal = false">
                    取消
                </n-button>
            </n-space>
        </template>
    </n-modal>

    <n-modal v-model:show="showModal2" preset="card" title="考勤详情" size="huge" :bordered="false"
        style="width: 1024px" :mask-closable="false" to="body">
        <div class="detail-modal">
            <n-button type="primary" style="margin-bottom: 16px;" @click="exportExcel">
                导出为表格
            </n-button>
            <n-data-table :columns="detailColumns" :data="detailData" :bordered="true" :single-line="false"
                :max-height="500" />
        </div>
    </n-modal>

    <div class="page-task">
        <n-card class="page-card" embedded>
            <div class="page-header">
                <div class="page-header-info">
                    <h2 class="page-title">攻城任务</h2>
                    <p class="page-desc">任务数量：{{ taskNum }}</p>
                </div>
                <n-space>
                    <n-button @click="getTaskList" :loading="loading">
                        <template #icon><RefreshCw :size="16" /></template>
                        刷新
                    </n-button>
                    <n-button type="primary" @click="addtaskshow = true">
                        <template #icon><Plus :size="16" /></template>
                        新增任务
                    </n-button>
                </n-space>
            </div>

            <div class="task-list" v-if="tasks.length > 0">
                <div class="task-card" v-for="task in tasks" :key="task.id">
                    <div class="task-header">
                        <div class="task-title-row">
                            <span class="task-name">{{ task.name }}</span>
                            <span class="task-coords">({{ splitwid(task.pos) }})</span>
                        </div>
                    </div>

                    <div class="task-stats">
                        <div class="task-stat-item">
                            <span class="task-stat-label">目标分组</span>
                            <div class="task-stat-tags">
                                <n-tag :bordered="false" type="info" size="small" v-for="g in task.target" :key="g">
                                    {{ g }}
                                </n-tag>
                            </div>
                        </div>
                        <div class="task-stat-item">
                            <span class="task-stat-label">目标人数</span>
                            <span class="task-stat-value">{{ task.target_user_num }}</span>
                        </div>
                        <div class="task-stat-item">
                            <span class="task-stat-label">实到人数</span>
                            <span class="task-stat-value highlight">{{ task.complete_user_num }}</span>
                        </div>
                        <div class="task-stat-item">
                            <span class="task-stat-label">任务时间</span>
                            <span class="task-stat-value">{{ formatTimestampMs(task.time) }}</span>
                        </div>
                    </div>

                    <div class="task-actions">
                        <n-button size="small" @click="getTaskDetail(task.id)">
                            <template #icon><Eye :size="14" /></template>
                            考勤详情
                        </n-button>
                        <n-button type="info" size="small" @click="enableGetReport(task.id, task.pos)">
                            <template #icon><Play :size="14" /></template>
                            开始考勤
                        </n-button>
                        <n-popconfirm @positive-click="delTaskReport(task.id)" :show-icon="false">
                            <template #trigger>
                                <n-button type="warning" size="small">
                                    <template #icon><Eraser :size="14" /></template>
                                    清理战报
                                </n-button>
                            </template>
                            确认清理战报吗？数据删除后无法恢复。<br>清理战报可以减少统计考勤的耗时
                        </n-popconfirm>
                        <n-popconfirm @positive-click="delTask(task.id)" :show-icon="false">
                            <template #trigger>
                                <n-button type="error" size="small">
                                    <template #icon><Trash2 :size="14" /></template>
                                    删除任务
                                </n-button>
                            </template>
                            确认删除该任务吗？
                        </n-popconfirm>
                    </div>
                </div>
            </div>
            <n-empty v-else description="暂无攻城任务" style="padding: 60px 0;" />
        </n-card>
    </div>
</template>

<style scoped lang="scss">
.page-task {
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

.task-list {
    display: flex;
    flex-direction: column;
    gap: 12px;
}

.task-card {
    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: 10px;
    padding: 20px;
    transition: box-shadow 0.2s, transform 0.2s;

    &:hover {
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
        transform: translateY(-1px);
    }
}

.task-header {
    margin-bottom: 16px;
}

.task-title-row {
    display: flex;
    align-items: baseline;
    gap: 8px;
}

.task-name {
    font-size: 16px;
    font-weight: 600;
    color: var(--color-text);
}

.task-coords {
    font-size: 13px;
    color: var(--color-text-secondary);
}

.task-stats {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 12px 24px;
    margin-bottom: 16px;
    padding: 16px;
    background: var(--color-bg);
    border-radius: 8px;
}

.task-stat-item {
    display: flex;
    flex-direction: column;
    gap: 4px;
}

.task-stat-label {
    font-size: 12px;
    color: var(--color-text-secondary);
}

.task-stat-value {
    font-size: 14px;
    color: var(--color-text);
    font-weight: 500;

    &.highlight {
        color: var(--color-accent);
    }
}

.task-stat-tags {
    display: flex;
    flex-wrap: wrap;
    gap: 6px;
}

.task-actions {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
}

.modal-form {
    display: flex;
    flex-direction: column;
    gap: 8px;
}

.report-modal {
    text-align: center;
}

.report-tip {
    font-size: 14px;
    color: var(--color-text-secondary);
    margin-bottom: 24px;
    line-height: 1.6;
}

.report-counter {
    display: flex;
    justify-content: center;
    padding: 20px 0;
}

.detail-modal {
    overflow: auto;

    :deep(.n-data-table-sorter) {
        opacity: 0;
        color: var(--color-primary);
        transition: opacity 0.15s;
    }

    :deep(th:hover .n-data-table-sorter) {
        opacity: 1;
    }

    :deep(.n-data-table-th--sorting .n-data-table-sorter) {
        opacity: 1;
        color: var(--color-primary);
    }
}
</style>
