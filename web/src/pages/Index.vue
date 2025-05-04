<script setup>
import { Star, Link, Server, KeyRound, Plus, BotMessageSquare,Users,ClipboardList } from 'lucide-vue-next';
import { ref } from "vue";
const isLogin = ref(true);
</script>

<template>
    <div class="bikamoeapp">
        <div class="bikamoeapp-content">
            <div class="bikamoeapp-title">
                <h2 style="margin-bottom: 4px;">攻城考勤助手</h2>
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
            </div>

            <div class="bikamoeapp-login" v-else>
                <!-- <img src="https://cdn.fastly.steamstatic.com/steamcommunity/public/images/items/1406990/915b1b4a05133186525a956d7ca5c142a3c3c9f3.webm" alt=""> -->
                <div class="login-box">
                    <h2 style="line-height: 48px;">Login</h2>
                    <p style="font-size: 13px;color: #71717a;">Enter your username below to login to your account.</p>
                    <div class="input-box">
                        <label>UserName</label>
                        <input type="text" />
                    </div>
                    <div class="input-box">
                        <label>PassWord</label>
                        <input type="password" />
                    </div>
                    <a class="button">Login</a>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped></style>