<script setup>
import { ref,onMounted } from "vue";
import { NList, NSpace, useMessage, NListItem, NTag, NThing, NButton,useDialog  } from 'naive-ui'
import { ApiGetTeamUser } from '@/api'

const dialog = useDialog();
const teamUsers = ref([]);
const usersNum = ref(0);

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

function getUserList(){
    teamUsers.value = [];
    usersNum.value = 0;
    ApiGetTeamUser().then(v => {
        if(v.status == 200){
            let resp = v.data;
            let data = resp.data;
            console.log(data);
            teamUsers.value = data;
            usersNum.value = data.length;
        }else{
            console.log("请求错误...");
        }
    }).catch(e => {

    });
}

onMounted(() => {
    getUserList()
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
    <div class="bikamoeapp">
        <div class="bikamoeapp-content">
            <div class="bikamoeapp-title">
                <h2 style="margin-bottom: 4px;">攻城考勤助手</h2>
                <p>同盟成员列表</p>
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
                    <a class="button" @click="getUserList">
                        刷新
                    </a>
                    <a class="button" @click="syncuser">
                        同步成员
                    </a>
                    <a class="button">
                        成员数量:{{ usersNum }}
                    </a>
                </div>
                <div>
                    <n-list hoverable clickable>
                        <n-list-item v-for="user in teamUsers">
                            <n-thing content-style="margin-top: 10px;">
                                <template #header>
                                    {{ user.name }}
                                    <n-tag :bordered="false" type="info" size="small">
                                        {{ user.group }}
                                    </n-tag>
                                </template>
                                <!-- <template #description> -->
                                    <!-- <n-space size="small" style="margin-top: 4px"> -->
                                        <p>ID：{{ user.id }}</p>
                                        <p>势力：{{ user.power }}</p>
                                        <p>武勋：{{ user.wu }}</p>
                                        <p>总贡献：{{ user.contribute_total }}</p>
                                        <p>周贡献：{{ user.contribute_week }}</p>
                                        <p>位置：({{ splitwid(user.pos) }})</p>
                                        <p>进盟时间：{{ formatTimestamp(user.join_time) }}</p>
                                    <!-- </n-space> -->
                                <!-- </template> -->
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