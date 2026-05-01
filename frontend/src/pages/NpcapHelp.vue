<script setup>
import { NCard, NButton, NAlert, NSpace } from 'naive-ui'
import { Download, ExternalLink, ShieldCheck } from 'lucide-vue-next'
import { BrowserOpenURL } from '../../wailsjs/runtime/runtime'

const openDownload = () => {
    BrowserOpenURL('https://npcap.com/#download')
}

const openNpcapGuide = () => {
    BrowserOpenURL('https://npcap.com/guide/npcap-users-guide.html')
}
</script>

<template>
    <div class="npcap-page">
        <div class="npcap-container">
            <div class="npcap-header">
                <div class="npcap-icon">
                    <ShieldCheck :size="28" />
                </div>
                <h1 class="npcap-title">需要安装 Npcap</h1>
                <p class="npcap-desc">本软件依赖 Npcap 进行网络数据包捕获，请先安装后再使用</p>
            </div>

            <n-card class="npcap-card" embedded :bordered="false">
                <n-alert type="warning" :show-icon="true" style="margin-bottom: 20px;">
                    Npcap 是 Windows 平台的网络数据包捕获库，本软件通过它监听游戏网络流量以实现自动化功能。未安装将无法正常使用。
                </n-alert>

                <div class="npcap-steps">
                    <h3 class="steps-title">安装步骤</h3>
                    <div class="step-item">
                        <span class="step-num">1</span>
                        <div class="step-content">
                            <span class="step-text">访问 Npcap 官网下载安装包</span>
                        </div>
                    </div>
                    <div class="step-item">
                        <span class="step-num">2</span>
                        <div class="step-content">
                            <span class="step-text">运行安装程序，按默认选项安装即可</span>
                        </div>
                    </div>
                    <div class="step-item">
                        <span class="step-num">3</span>
                        <div class="step-content">
                            <span class="step-text">安装完成后重新启动本软件</span>
                        </div>
                    </div>
                </div>
            </n-card>

            <div class="npcap-actions">
                <n-button type="primary" @click="openDownload">
                    <template #icon><Download :size="16" /></template>
                    下载 Npcap
                </n-button>
                <n-button @click="openNpcapGuide">
                    <template #icon><ExternalLink :size="16" /></template>
                    查看安装指南
                </n-button>
            </div>
        </div>
    </div>
</template>

<style scoped lang="scss">
.npcap-page {
    display: flex;
    align-items: flex-start;
    justify-content: center;
    min-height: 100%;
    padding-top: 80px;
    background: var(--color-bg);
}

.npcap-container {
    width: 100%;
    max-width: 520px;
}

.npcap-header {
    text-align: center;
    margin-bottom: 24px;
}

.npcap-icon {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 56px;
    height: 56px;
    border-radius: 14px;
    background: linear-gradient(135deg, #f59e0b, #ef4444);
    color: #fff;
    margin-bottom: 16px;
}

.npcap-title {
    font-size: 22px;
    font-weight: 700;
    color: var(--color-text);
    margin: 0 0 6px;
}

.npcap-desc {
    font-size: 14px;
    color: var(--color-text-secondary);
    margin: 0;
}

.npcap-card {
    border-radius: 12px;
}

.npcap-steps {
    display: flex;
    flex-direction: column;
    gap: 12px;
}

.steps-title {
    font-size: 15px;
    font-weight: 600;
    color: var(--color-text);
    margin: 0 0 4px;
}

.step-item {
    display: flex;
    align-items: center;
    gap: 12px;
}

.step-num {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 28px;
    height: 28px;
    border-radius: 50%;
    background: var(--color-primary-light);
    color: var(--color-accent);
    font-size: 13px;
    font-weight: 600;
    flex-shrink: 0;
}

.step-text {
    font-size: 14px;
    color: var(--color-text);
}

.npcap-actions {
    display: flex;
    justify-content: center;
    gap: 12px;
    margin-top: 20px;
}
</style>
