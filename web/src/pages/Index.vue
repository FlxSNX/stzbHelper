<script setup>
import { Star, Link, Server, KeyRound, Plus, BotMessageSquare,Users,ClipboardList,Swords } from 'lucide-vue-next';
import { ref } from "vue";
const isLogin = ref(true);
</script>

<template>
    <div class="bikamoeapp">
        <div class="bikamoeapp-content">
            <div class="bikamoeapp-title">
                <h2 style="margin-bottom: 4px;">攻城考勤助手</h2>
                <p>Version: 0.0.2</p>
            </div>
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
            </div>
        </div>
    </div>
</template>

<style scoped></style>