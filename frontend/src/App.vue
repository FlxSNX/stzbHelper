<script setup>
import { h, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { NMessageProvider, NDialogProvider, NConfigProvider, NLayout, NLayoutSider, NLayoutContent, NMenu, NIcon, NButton } from 'naive-ui'
import { zhCN, dateZhCN, darkTheme } from 'naive-ui'
import { Home, Users, ClipboardList, Swords, UserRoundSearch, ScrollText, Bug, Moon, Sun, BookOpen, MessageSquare, BarChart3 } from 'lucide-vue-next'
import { useThemeStore } from './stores/theme'

import TitleBar from './components/TitleBar.vue'

const route = useRoute()
const router = useRouter()
const themeStore = useThemeStore()

const activeKey = computed(() => {
    const path = route.path
    if (path === '/') return 'index'
    return path.replace('/', '')
})

const isFullscreenPage = computed(() => route.path === '/select-db' || route.path === '/npcap-help')

function renderIcon(icon) {
    return () => h(NIcon, null, { default: () => h(icon, { size: 18 }) })
}

const menuOptions = [
    {
        label: '控制面板',
        key: 'index',
        icon: renderIcon(Home)
    },
    {
        label: '同盟成员',
        key: 'teamuser',
        icon: renderIcon(Users)
    },
    {
        label: '攻城任务',
        key: 'task',
        icon: renderIcon(ClipboardList)
    },
    {
        label: '分组武勋',
        key: 'groupWu',
        icon: renderIcon(Swords)
    },
    {
        label: '队伍查询',
        key: 'teamquery',
        icon: renderIcon(UserRoundSearch)
    },
    {
        label: '队伍胜率',
        key: 'teamwinrate',
        icon: renderIcon(BarChart3)
    },
    {
        label: '主公簿',
        key: 'book',
        icon: renderIcon(BookOpen)
    },
    // {
    //     label: '战役叫阵',
    //     key: 'battlecall',
    //     icon: renderIcon(MessageSquare)
    // },
    {
        label: '运行日志',
        key: 'logs',
        icon: renderIcon(ScrollText)
    },
    {
        label: 'API 调试',
        key: 'debug',
        icon: renderIcon(Bug)
    }
]

const handleMenuUpdate = (key) => {
    router.push(key === 'index' ? '/' : `/${key}`)
}

const themeOverrides = computed(() => ({
    common: {
        primaryColor: themeStore.isDark ? '#e5e5e5' : '#3b82f6',
        primaryColorHover: themeStore.isDark ? '#d4d4d4' : '#2563eb',
        primaryColorPressed: themeStore.isDark ? '#a3a3a3' : '#1d4ed8',
        borderRadius: '8px',
        borderRadiusSmall: '6px',
        fontFamily: 'geist-sans, ui-sans-serif, system-ui, sans-serif',
    },
    Card: {
        borderRadius: '12px',
        paddingMedium: '20px',
    },
    Table: {
        borderRadius: '12px',
        thColor: themeStore.isDark ? '#262626' : '#f8f9fb',
        tdColorHover: themeStore.isDark ? '#262626' : '#f1f3f5',
    },
    Button: {
        borderRadiusMedium: '8px',
        borderRadiusSmall: '6px',
    },
    Menu: {
        borderRadius: '8px',
        itemColorActive: themeStore.isDark ? 'rgba(255, 255, 255, 0.06)' : '#eff6ff',
        itemColorActiveHover: themeStore.isDark ? 'rgba(255, 255, 255, 0.06)' : '#eff6ff',
        itemTextColorActive: themeStore.isDark ? '#ffffff' : '#3b82f6',
        itemTextColorActiveHover: themeStore.isDark ? '#ffffff' : '#3b82f6',
        itemIconColorActive: themeStore.isDark ? '#ffffff' : '#3b82f6',
        itemIconColorActiveHover: themeStore.isDark ? '#ffffff' : '#3b82f6',
    },
}))
</script>

<template>
    <n-config-provider :locale="zhCN" :date-locale="dateZhCN" :theme="themeStore.isDark ? darkTheme : undefined" :theme-overrides="themeOverrides">
        <n-dialog-provider>
            <n-message-provider>
                <div class="app-shell">
                    <TitleBar />
                    <n-layout v-if="!isFullscreenPage" has-sider class="app-layout">
                        <n-layout-sider
                            bordered
                            :width="220"
                            :native-scrollbar="false"
                            content-style="display: flex; flex-direction: column; height: 100%;"
                        >
                            <n-menu
                                :value="activeKey"
                                :options="menuOptions"
                                @update:value="handleMenuUpdate"
                            />
                            <div class="sidebar-bottom">
                                <n-button quaternary circle size="small" @click="themeStore.toggle" :title="themeStore.isDark ? '切换浅色模式' : '切换深色模式'">
                                    <template #icon>
                                        <Sun v-if="themeStore.isDark" :size="16" />
                                        <Moon v-else :size="16" />
                                    </template>
                                </n-button>
                            </div>
                        </n-layout-sider>
                        <n-layout>
                            <n-layout-content
                                :native-scrollbar="false"
                                class="main-content"
                            >
                                <router-view v-slot="{ Component, route: r }">
                                    <keep-alive include="Index">
                                        <component :is="Component" :key="r.path" />
                                    </keep-alive>
                                </router-view>
                            </n-layout-content>
                        </n-layout>
                    </n-layout>
                    <n-layout v-else class="app-layout">
                        <n-layout-content
                            :native-scrollbar="false"
                            class="main-content"
                        >
                            <router-view />
                        </n-layout-content>
                    </n-layout>
                </div>
            </n-message-provider>
        </n-dialog-provider>
    </n-config-provider>
</template>

<style scoped>
.app-shell {
    display: flex;
    flex-direction: column;
    height: 100vh;
    width: 100vw;
    overflow: hidden;
}

.app-layout {
    flex: 1;
    min-height: 0;
}

.sidebar-bottom {
    margin-top: auto;
    padding: 12px 16px;
    display: flex;
    align-items: center;
    justify-content: flex-end;
}

.main-content {
    padding: 24px;
    background: var(--color-bg);
    min-height: 100%;
}

:deep(.n-layout-scroll-container) {
    background: var(--color-bg);
}
</style>
