<script setup>
import { Star, Link, Server, KeyRound, Plus, BotMessageSquare,Users,ClipboardList,Swords,UserRoundSearch } from 'lucide-vue-next';
import { ref } from "vue";
import { useMessage,NDivider } from 'naive-ui'
import { ApiEnableGetBattleReport,ApiDisableGetBattleReport } from '../api';
const isLogin = ref(true);
const nmessage = useMessage()
const EnableGetBattleReport = () =>{
    ApiEnableGetBattleReport().then(v => {
        if(v.data.code == 200){
            nmessage.success("开启成功")
        }else{
            nmessage.error(v.data.msg);
        }
    }).catch(e => {
        nmessage.error("开启获取战报详情失败:"+e);
    });
}

const DisableGetBattleReport = () =>{
    ApiDisableGetBattleReport().then(v => {
        if(v.data.code == 200){
            nmessage.success("关闭成功")
        }else{
            nmessage.error(v.data.msg);
        }
    }).catch(e => {
        nmessage.error("关闭获取战报详情失败:"+e);
    });
}

</script>

<template>
    <div class="bikamoeapp">
        <div class="bikamoeapp-content">
            <div class="bikamoeapp-title">
                <h2 style="margin-bottom: 4px;">率土之滨助手</h2>
                <p>Version: 0.0.3</p>
            </div>
            <div style="margin: 16px auto;padding: 16px;">
                <div style="font-size: 18px;font-weight: 600;padding-left: 4px;">控制面板</div>
                <div style="border: 1px solid rgba(228, 228, 231, 0.6); padding: 16px;display: flex;flex-wrap: wrap;flex-direction: column;">
                    <div style="padding: 8px;margin: 8px;flex-shrink: 0;">
                        <div style="font-size: 18px;font-weight: 600;">获取详细战报</div>
                        <span style="font-size: 13px;">用于查询队伍功能拉取战报使用,开启时无法获取攻城战报</span>
                        <div style="margin: 16px 0;">
                            <a class="button" style="margin-right: 8px;" @click="EnableGetBattleReport">开启</a>
                            <a class="button" @click="DisableGetBattleReport">关闭</a>
                        </div>
                    </div>
                </div>
            </div>
            <n-divider />
            <div class="bikamoeapp-list" v-if="isLogin">
                <div class="bikamoeapp-appitem-warp">
                    <div class="item">
                        <div class="icon">
                            <Users :size="64" />
                        </div>
                        <div class="bottom">
                            <div class="name">同盟成员</div>
                            <router-link class="button" to="/teamuser">进入</router-link>
                        </div>
                    </div>
                </div>
                <div class="bikamoeapp-appitem-warp">
                    <div class="item">
                        <div class="icon">
                            <ClipboardList :size="64" />
                        </div>
                        <div class="bottom">
                            <div class="name">攻城任务</div>
                            <router-link class="button" to="/task">进入</router-link>
                        </div>
                    </div>
                </div>
                <div class="bikamoeapp-appitem-warp">
                    <div class="item">
                        <div class="icon">
                            <Swords :size="64" />
                        </div>
                        <div class="bottom">
                            <div class="name">分组武勋</div>
                            <router-link class="button" to="/groupWu">进入</router-link>
                        </div>
                    </div>
                </div>
                <div class="bikamoeapp-appitem-warp">
                    <div class="item">
                        <div class="icon">
                            <UserRoundSearch :size="64" />
                        </div>
                        <div class="bottom">
                            <div class="name">队伍查询</div>
                            <a class="button" href="/data.html#/team" target="_blank">进入</a>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped></style>