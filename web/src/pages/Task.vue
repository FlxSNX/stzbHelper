<script setup>
import { ref, onMounted } from "vue";
import { NDrawerContent,NDrawer,NInput,NFormItem,NList, NSpace, useMessage, NListItem, NTag, NThing, NButton, useDialog,NSelect,NDatePicker } from 'naive-ui'
import { ApiGetTeamGroup,ApiCreateTask } from '@/api'

const nmessage = useMessage()
const dialog = useDialog();
const teamUsers = ref([]);
const usersNum = ref(0);
const addtaskshow = ref(false);
const targetgroup = ref([])
const grouplist =ref([]);
const tasktime = ref(new Date().getTime());
const taskname = ref("")
const taskpos = ref()
const createing = ref(false)

const syncuser = () => {
    dialog.info({
        title: '信息',
        content: '请前往游戏中,点开同盟成员列表即可同步',
        positiveText: '确认',
        transformOrigin: "center",
        onPositiveClick: () => {

        }
    })
}

const createTask = () => {
    createing.value = true;
    console.log("taskname",taskname.value);
    console.log("tasktime",tasktime.value);
    console.log("targetgroup",targetgroup.value);
    console.log("taskpos",taskpos.value);

    ApiCreateTask({
        taskname:taskname.value,
        tasktime:tasktime.value,
        targetgroup:targetgroup.value,
        taskpos:taskpos.value,
    }).then(v => {
        if(v.status == 200){
            if(v.data.code == 200){
                nmessage.success(v.data.msg)
            }else{
                nmessage.error(v.data.msg)
            }
        }else{
            nmessage.error("创建出错")
        }
        createing.value = false;
    }).catch(e => {
        createing.value = false;
        nmessage.error(e)
    })
}

function getUserList() {
    teamUsers.value = [];
    usersNum.value = 0;
    ApiGetTeamUser().then(v => {
        if (v.status == 200) {
            let resp = v.data;
            let data = resp.data;
            console.log(data);
            teamUsers.value = data;
            usersNum.value = data.length;
        } else {
            console.log("请求错误...");
        }
    }).catch(e => {

    });
}

onMounted(() => {
    ApiGetTeamGroup().then(v => {
        if(v.status == 200){
            let resp = v.data;
            let data = resp.data;
            console.log(data);
            grouplist.value = [];
            data.forEach(e =>{
                grouplist.value.push({
                    label: e,
                    value: e
                });
            });
        }
    })
});

function formatTimestamp(timestamp) {
    // 将秒级时间戳转换为毫秒级
    const date = new Date(timestamp * 1000);

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
</script>

<template>
    <n-drawer v-model:show="addtaskshow" :min-width="512" :max-width="1024" :default-width="512" placement="right" :mask-closable="true" :resizable="true">
        <n-drawer-content title="新增任务" :native-scrollbar="false">
            <n-form-item label="任务名称">
                <n-input v-model:value="taskname" placeholder="例如：内黄LV5 或者你也可以随意填写" />
            </n-form-item>
            <n-form-item label="任务坐标">
                <n-input
                    pair
                    separator="，"
                    :placeholder="['X坐标','Y坐标']"
                    v-model:value="taskpos"
                    clearable
                />
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
                    <a class="button" @click="getUserList">
                        刷新
                    </a>
                    <a class="button" @click="addtaskshow = true">
                        新增任务
                    </a>
                    <a class="button">
                        任务数量:{{ usersNum }}
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

                        <n-list-item>
                            <n-thing content-style="margin-top: 10px;">
                                <template #header>
                                    内黄LV5 (361,793)
                                </template>
                                <n-space style="margin-bottom: 16px;">
                                    <n-tag type="success">已完成</n-tag>
                                </n-space>
                               
                                <p>目标分组&nbsp;：&nbsp;&nbsp;<n-tag :bordered="false" type="info" size="small">分组1</n-tag></p>
                                <p>目标人数&nbsp;：&nbsp;&nbsp;50</p>
                                <p>实到人数&nbsp;：&nbsp;&nbsp;0</p>
                                <p>任务时间&nbsp;：&nbsp;&nbsp;{{ formatTimestamp(1746174023) }}</p>
                                <p>同步时间&nbsp;：&nbsp;&nbsp;{{ formatTimestamp(1746174023) }}</p>
                                <n-space style="margin-top: 8px;">
                                    <n-button size="small">
                                        查看详情
                                    </n-button>
                                    <n-button type="info" size="small">
                                        开始考勤
                                    </n-button>
                                    <n-button type="info" size="small">
                                        标记为已完成
                                    </n-button>
                                    <n-button type="error" size="small">
                                        删除任务
                                    </n-button>
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