<script setup>
import { ref, onMounted } from "vue";
import {  useMessage, useDialog,NTable } from 'naive-ui'
import { ApiGetGroupWu } from '@/api'
import * as XLSX from 'xlsx';

const nmessage = useMessage()
const groupdata = ref([])


function getData() {
    groupdata.value = [];
    ApiGetGroupWu().then(v => {
        if (v.status == 200) {
            let resp = v.data;
            if (resp.code == 200) {
                let data = resp.data;
                groupdata.value = data;
            } else {
                nmessage.error(resp.msg);
            }
        } else {
            nmessage.error("获取分组武勋数据失败");
        }
    });
}

onMounted(() => {
    getData();
});
</script>

<template>
    <div class="bikamoeapp">
        <div class="bikamoeapp-content">
            <div class="bikamoeapp-title">
                <h2 style="margin-bottom: 4px;">攻城考勤助手</h2>
                <p>分组本周武勋</p>
                <p>更新武勋数据请同步成员数据</p>
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
                    <a class="button" @click="getData">
                        刷新
                    </a>
                </div>
                <div>
                   <n-table>
                        <thead>
                            <tr>
                                <th>分组名称</th>
                                <th>人数</th>
                                <th>总武勋</th>
                                <th>平均武勋</th>
                                <th>0武勋人数</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="u in groupdata">
                                <td>{{ u.group }}</td>
                                <td>{{ u.member_count }}</td>
                                <td>{{ u.total_wu }}</td>
                                <td>{{ u.average_wu }}</td>
                                <td>{{ u.zero_wu_count }}</td>
                            </tr>
                        </tbody>
                    </n-table>
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