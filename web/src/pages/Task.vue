<script setup>
import { ref, onMounted } from "vue";
import { NDrawerContent, NDrawer, NInput, NFormItem, NList, NSpace, useMessage, NListItem, NTag, NThing, NButton, useDialog, NSelect, NDatePicker, NPopconfirm, NModal, NCard,NSpin,NTable } from 'naive-ui'
import { ApiGetTeamGroup, ApiCreateTask, ApiGetTaskList, ApiDelTask,ApiEnableGetReport,ApiGetReportNumByTaskId,ApiStatisticsReport,ApiGetTask,ApiDelTaskReport } from '@/api'
import { getRect } from "naive-ui/es/affix/src/utils";
import * as XLSX from 'xlsx';

const nmessage = useMessage()
const dialog = useDialog();
const addtaskshow = ref(false);
const targetgroup = ref([])
const grouplist = ref([]);
const tasktime = ref(new Date().getTime());
const taskname = ref("")
const taskpos = ref()
const createing = ref(false)
const tasks = ref([]);
const taskNum = ref(0);

const createTask = () => {
    createing.value = true;
    console.log("taskname", taskname.value);
    console.log("tasktime", tasktime.value);
    console.log("targetgroup", targetgroup.value);
    console.log("taskpos", taskpos.value);

    ApiCreateTask({
        taskname: taskname.value,
        tasktime: tasktime.value,
        targetgroup: targetgroup.value,
        taskpos: taskpos.value,
    }).then(v => {
        if (v.status == 200) {
            if (v.data.code == 200) {
                nmessage.success(v.data.msg)
                taskname.value = "";
                targetgroup.value = [];
                taskpos.value = [];
                getTaskList();
            } else {
                nmessage.error(v.data.msg)
            }
        } else {
            nmessage.error("创建出错")
        }
        createing.value = false;
    }).catch(e => {
        createing.value = false;
        nmessage.error(e)
    })
}

const delTask = (id) => {
    // nmessage.info(id);
    ApiDelTask(id).then(v => {
        if (v.status == 200) {
            if (v.data.code == 200) {
                nmessage.success(v.data.msg);
                getTaskList();
            } else {
                nmessage.error(v.data.msg);
            }
        } else {
            nmessage.error("任务删除失败" + v.status)
        }
    });
}

const delTaskReport = (id) => {
    ApiDelTaskReport(id).then(v => {
        if (v.status == 200) {
            if (v.data.code == 200) {
                nmessage.success(v.data.msg);
                getTaskList();
            } else {
                nmessage.error(v.data.msg);
            }
        } else {
            nmessage.error("任务删除失败" + v.status)
        }
    });
}

function getTaskList() {
    tasks.value = [];
    taskNum.value = 0;
    ApiGetTaskList().then(v => {
        if (v.status == 200) {
            let resp = v.data;
            if (resp.code == 200) {
                let data = resp.data;
                tasks.value = data;
                taskNum.value = data.length;
            } else {
                nmessage.error(resp.msg);
            }
        } else {
            nmessage.error("获取任务列表失败");
        }
    });
}

onMounted(() => {
    ApiGetTeamGroup().then(v => {
        if (v.status == 200) {
            let resp = v.data;
            let data = resp.data;
            console.log(data);
            grouplist.value = [];
            data.forEach(e => {
                grouplist.value.push({
                    label: e,
                    value: e
                });
            });
        }
    });
    getTaskList();
});

function formatTimestamp(timestamp) {
    // 将秒级时间戳转换为毫秒级
    const date = new Date(timestamp);

    // 获取年、月、日、时、分、秒
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0'); // 月份从0开始，所以要加1
    const day = String(date.getDate()).padStart(2, '0');
    const hours = String(date.getHours()).padStart(2, '0');
    const minutes = String(date.getMinutes()).padStart(2, '0');
    const seconds = String(date.getSeconds()).padStart(2, '0');

    // 拼接成目标格式
    return `${year}/${month}/${day} ${hours}:${minutes}:${seconds}`;
}

function splitwid(num) {
    // 将数字转换为字符串
    const numStr = num.toString();

    // 获取后四位
    const lastFour = numStr.slice(-4);

    // 获取前面的部分
    const firstPart = numStr.slice(0, -4);

    // 将后四位转换为数字，并去掉前导零
    const lastFourNumber = parseInt(lastFour, 10);

    // 返回结果
    return `${firstPart},${lastFourNumber}`
}

const showModal = ref(false);
const getReporting = ref(false);
const reportNum = ref(0);
const getReportNumTimer = ref(null);
const inStatistics = ref(false);
const curtaskid = ref(0);

const enableGetReport = (id,pos) => {
    showModal.value = true;
    ApiEnableGetReport({
        pos
    })
    getReporting.value = true;
    reportNum.value = 0;
    curtaskid.value = id;
    getReportNumTimer.value = setInterval(() => {
        ApiGetReportNumByTaskId(id).then(v => {
            if(v.status == 200 && v.data.code == 200){
                reportNum.value = v.data.data.count;
            }
        });
    },1000);
}

const statistics = () =>{
    clearInterval(getReportNumTimer.value);
    getReporting.value = false;
    inStatistics.value = true;
    ApiStatisticsReport(curtaskid.value).then(v => {
        if(v.data.code == 200){
            nmessage.success(v.data.msg);
            showModal.value = false;
            curtaskid.value = 0;
            getTaskList()
        }else{
            nmessage.error(v.data.msg);
        }
        inStatistics.value = false;
    }).catch(e => {
        inStatistics.value = false;
        nmessage.error("统计考勤数据失败:"+e);
    });
}

const showModal2 = ref(false);
const taskDetail = ref({})
const getTask = (id) => {
    taskDetail.value = {};
    showModal2.value = true;
    ApiGetTask(id).then(v => {
        if(v.data.code == 200){
            taskDetail.value = v.data.data;
        }else{
            nmessage.error(v.data.msg);
        }
    }).catch(e => {
        nmessage.error("获取考勤数据失败:"+e);
    });
}

const exportExcel = () => {
	let data = [];
	data.push([
		"名字",
		"分组",
		"主力",
		"拆迁",
		"主力次数",
		"拆迁次数",
	]);

    Object.values(taskDetail.value.user_list).forEach(v => {
        data.push([
            v.name,
            v.group,
            v.atk_team_num,
            v.dis_team_num,
            v.atk_num,
            v.dis_num
        ]);
    });
	
	// 创建工作表
	const ws = XLSX.utils.aoa_to_sheet(data);

	// 创建工作簿
	const wb = XLSX.utils.book_new();
	XLSX.utils.book_append_sheet(wb, ws, 'Sheet1'); // 工作表名称

	// 生成 Excel 文件并下载
	XLSX.writeFile(wb, `${taskDetail.value.name}考勤表.xlsx`); // 文件名
};
</script>

<template>
    <n-drawer v-model:show="addtaskshow" :min-width="512" :max-width="1024" :default-width="512" placement="right" :mask-closable="false" :resizable="true">
        <n-drawer-content title="新增任务" :native-scrollbar="false">
            <n-form-item label="任务名称">
                <n-input v-model:value="taskname" placeholder="例如：内黄LV5 或者你也可以随意填写" />
            </n-form-item>
            <n-form-item label="任务坐标">
                <n-input pair separator="，" :placeholder="['X坐标', 'Y坐标']" v-model:value="taskpos" clearable />
            </n-form-item>
            <n-form-item label="任务时间">
                <n-date-picker v-model:value="tasktime" type="datetime" />
            </n-form-item>
            <n-form-item label="目标分组">
                <n-select v-model:value="targetgroup" multiple :options="grouplist" placeholder="" />
            </n-form-item>

            <n-space>
                <n-button strong secondary type="primary" :loading="createing" @click="createTask">
                    添加
                </n-button>
                <n-button strong secondary type="error" @click="addtaskshow = false">
                    关闭
                </n-button>
            </n-space>
        </n-drawer-content>
    </n-drawer>
    <n-modal v-model:show="showModal" class="custom-card" preset="card" title="攻城考勤" size="huge" :bordered="false" style="width:600px" :mask-closable="false">
        <div>
            <!-- <n-spin size="large" /> -->
            <p>请前往游戏中,到攻城任务坐标位置查看同盟战报,并勾选守城军士(否则获取不了拆迁战报)。然后一直往下滑直到没有战报为止</p>
            <p>系统只会获取攻城任务时间之后1小时内的战报(暂未实现时间限制)</p>
            <h2>已获取 <span style="color: #2080f0;">{{ reportNum }}</span> 封战报</h2>
        </div>
        <template #footer>
            <n-space>
                <n-button strong secondary type="info" :loading="true" v-if="getReporting">
                    获取战报中
                </n-button>
                <n-button strong secondary type="success" @click="statistics" :loading="inStatistics">
                    {{ inStatistics ? "统计考勤数据中" : "已获取完战报,开始统计考勤数据" }}
                </n-button>
            </n-space>
        </template>
    </n-modal>
    <n-modal v-model:show="showModal2" class="custom-card" preset="card" title="考勤详情" size="huge" :bordered="false" style="width:1024px"  :mask-closable="false">
        <div>
            <!-- <n-spin size="large" /> -->
            <n-button strong secondary type="info" style="margin-bottom: 8px;" @click="exportExcel">
                导出为表格
            </n-button>
            <n-table>
                <thead>
                    <tr>
                        <th>名称</th>
                        <th>分组</th>
                        <th>主力</th>
                        <th>拆迁</th>
                        <th>主力次数</th>
                        <th>拆迁次数</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="u in taskDetail.user_list">
                        <td>{{ u.name }}</td>
                        <td>{{ u.group }}</td>
                        <td>{{ u.atk_team_num }}</td>
                        <td>{{ u.dis_team_num }}</td>
                        <td>{{ u.atk_num }}</td>
                        <td>{{ u.dis_num }}</td>
                    </tr>
                </tbody>
            </n-table>
        </div>
        <!-- <template #footer>
            <n-space>
                <n-button strong secondary type="info" :loading="true" v-if="getReporting">
                    获取战报中
                </n-button>
                <n-button strong secondary type="success" @click="statistics" :loading="inStatistics">
                    {{ inStatistics ? "统计考勤数据中" : "已获取完战报,开始统计考勤数据" }}
                </n-button>
            </n-space>
        </template> -->
    </n-modal>
    <div class="bikamoeapp">
        <div class="bikamoeapp-content">
            <div class="bikamoeapp-title">
                <h2 style="margin-bottom: 4px;">攻城考勤助手</h2>
                <p>攻城任务列表</p>
            </div>
            <!-- <div class="bikamoeapp-list"> -->
            <div>
                <div style="width: 100%;
                    height: 48px;
                    border-bottom: 1px solid rgba(228, 228, 231, 0.6);
                    display: flex;
                    align-items: center;
                    padding: 0 8px;
                    box-sizing: border-box;">
                    <router-link class="button" to="/">返回</router-link>
                    <a class="button" @click="getTaskList">
                        刷新
                    </a>
                    <a class="button" @click="addtaskshow = true">
                        新增任务
                    </a>
                    <a class="button">
                        任务数量:{{ taskNum }}
                    </a>
                </div>
                <div>
                    <n-list hoverable clickable>
                        <!-- <n-list-item v-for="user in teamUsers">
                            <n-thing content-style="margin-top: 10px;">
                                <template #header>
                                    {{ user.name }}
                                    <n-tag :bordered="false" type="info" size="small">
                                        {{ user.group }}
                                    </n-tag>
                                </template>
                                <p>ID：{{ user.id }}</p>
                                <p>势力：{{ user.power }}</p>
                                <p>武勋：{{ user.wu }}</p>
                                <p>总贡献：{{ user.contribute_total }}</p>
                                <p>周贡献：{{ user.contribute_week }}</p>
                                <p>位置：({{ splitwid(user.pos) }})</p>
                                <p>进盟时间：{{ formatTimestamp(user.join_time) }}</p>
                                </n-thing>
                            </n-list-item> -->

                        <n-list-item v-for="task in tasks">
                            <n-thing content-style="margin-top: 10px;">
                                <template #header>
                                    {{ task.name }} ({{ splitwid(task.pos) }})
                                </template>
                                <!-- <n-space style="margin-bottom: 16px;">
                                    <n-tag type="success" v-if="task.status == 1">已完成</n-tag>
                                    <n-tag type="info" v-if="task.status === 0">待考勤</n-tag>
                                </n-space> -->

                                <p>目标分组&nbsp;：&nbsp;
                                    <n-tag :bordered="false" type="info" size="small" v-for="g in task.target"
                                        style="margin-right: 8px;">{{ g }}</n-tag>
                                </p>
                                <p>目标人数&nbsp;：&nbsp;&nbsp;{{ task.target_user_num }}</p>
                                <p>实到人数&nbsp;：&nbsp;&nbsp;{{ task.complete_user_num }}</p>
                                <p>任务时间&nbsp;：&nbsp;&nbsp;{{ formatTimestamp(task.time) }}</p>
                                <!-- <p>同步时间&nbsp;：&nbsp;&nbsp;{{ formatTimestamp(1746174023) }}</p> -->
                                <n-space style="margin-top: 8px;">
                                    <n-button size="small" @click="getTask(task.id)">
                                        考勤详情
                                    </n-button>
                                    <n-button type="info" size="small" @click="enableGetReport(task.id,task.pos)">
                                        开始考勤
                                    </n-button>
                                    <n-popconfirm @positive-click="delTaskReport(task.id)" :show-icon="false">
                                        <template #trigger>
                                            <n-button type="error" size="small">
                                                清理战报
                                            </n-button>
                                        </template>
                                        确认清理战报吗? 数据删除后无法恢复。<br>清理战报可以减少统计考勤的耗时
                                    </n-popconfirm>
                                    <n-popconfirm @positive-click="delTask(task.id)" :show-icon="false">
                                        <template #trigger>
                                            <n-button type="error" size="small">删除任务</n-button>
                                        </template>
                                        确认删除该任务吗?
                                    </n-popconfirm>
                                </n-space>
                            </n-thing>
                        </n-list-item>
                    </n-list>
                </div>
            </div>

            <!-- </div> -->
        </div>
    </div>
</template>

<style scoped>
a.button {
    border-bottom: 1px solid rgb(228 228 231 / 60%);
    margin-right: 8px;
    font-size: 14px;
}
</style>