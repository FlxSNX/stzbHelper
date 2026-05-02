<script setup>
import { ref, watch, onMounted, onUnmounted } from 'vue'
import { NCard, NGrid, NGi, NStatistic, NSwitch, NEmpty, NCollapse, NCollapseItem, useMessage } from 'naive-ui'
import { EnableBookData, DisableBookData } from '../../wailsjs/go/main/App'
import { EventsOn, EventsOff } from '../../wailsjs/runtime/runtime'
import { herocfg } from '../cfg'
import { splitwid } from '@/utils/format'

const heroMap = JSON.parse(herocfg)
const nmessage = useMessage()
const enabled = ref(sessionStorage.getItem('bookDataEnabled') === 'true')
const bookData = ref(null)

const getHeroIconId = (id) => {
    if (!id) return id
    const num = Number(id)
    const resolved = num >= 130000 ? num - 30000 : num
    const hero = heroMap[String(resolved)]
    return hero ? hero.iconId : resolved
}

const getHeroName = (id) => {
    if (!id) return ''
    const num = Number(id)
    const resolved = num >= 130000 ? num - 30000 : num
    const hero = heroMap[String(resolved)]
    return hero ? hero.name : `未知(${id})`
}

const formatNum = (val) => {
    const n = Number(val)
    if (isNaN(n)) return val
    if (n >= 10000) return (n / 10000).toFixed(2) + '万'
    return n
}

const formatPos = (pos) => {
    if (!pos) return '-'
    return splitwid(pos)
}

let unsubscribe = null

const handleSwitch = (val) => {
    sessionStorage.setItem('bookDataEnabled', val)
    if (val) {
        EnableBookData().then(v => {
            let resp = JSON.parse(v)
            if (resp.code == 200) {
                nmessage.success('已开启主公簿数据推送')
            }
        })
    } else {
        DisableBookData().then(v => {
            let resp = JSON.parse(v)
            if (resp.code == 200) {
                nmessage.success('已关闭主公簿数据推送')
            }
        })
    }
}

onMounted(() => {
    unsubscribe = EventsOn('bookData', (data) => {
        bookData.value = data
    })
    if (enabled.value) {
        EnableBookData()
    }
})

onUnmounted(() => {
    if (unsubscribe) EventsOff('bookData')
    if (enabled.value) DisableBookData()
})
</script>

<template>
    <div class="page-book">
        <n-card class="page-card" embedded>
            <div class="page-header">
                <div class="page-header-info">
                    <h2 class="page-title">主公簿</h2>
                    <p class="page-desc">开启开关后，在游戏中打开主公簿即可自动获取数据</p>
                </div>
                <div class="switch-wrap">
                    <span class="switch-label">{{ enabled ? '接收中' : '未开启' }}</span>
                    <n-switch v-model:value="enabled" @update:value="handleSwitch" />
                </div>
            </div>
        </n-card>

        <template v-if="bookData">
            <n-card class="page-card hero-card" embedded>
                <div class="hero-info">
                    <div class="hero-main">
                        <h2 class="hero-name">{{ bookData.role_name || '-' }}</h2>
                        <div class="hero-meta">
                            <span v-if="bookData.server">服务器: {{ bookData.server }}</span>
                            <span v-if="bookData.alliance_name">同盟: {{ bookData.alliance_name }}</span>
                            <span v-if="bookData.group_name">分组: {{ bookData.group_name }}</span>
                            <span v-if="bookData.ip">IP: {{ bookData.ip }}</span>
                            <span v-if="bookData.location">地区: {{ bookData.location }}</span>
                        </div>
                    </div>
                    <div class="hero-stats">
                        <n-statistic label="势力值" :value="bookData.power || 0" />
                        <n-statistic label="主城坐标" :value="formatPos(bookData.main_city_pos)" />
                    </div>
                </div>
            </n-card>

            <n-grid :cols="3" :x-gap="16" :y-gap="16" class="stat-grid">
                <n-gi>
                    <n-card embedded size="small">
                        <n-statistic label="登录天数" :value="bookData.login_days || 0" />
                    </n-card>
                </n-gi>
                <n-gi>
                    <n-card embedded size="small">
                        <n-statistic label="最高赛季灭敌" :value="formatNum(bookData.max_season_kills)" />
                    </n-card>
                </n-gi>
                <n-gi>
                    <n-card embedded size="small">
                        <n-statistic label="历史最高武勋" :value="formatNum(bookData.max_merit)" />
                    </n-card>
                </n-gi>
                <n-gi>
                    <n-card embedded size="small">
                        <n-statistic label="历史最高势力" :value="bookData.max_power || 0" />
                    </n-card>
                </n-gi>
                <n-gi>
                    <n-card embedded size="small">
                        <n-statistic label="赛季参与数" :value="bookData.season_count || 0" />
                    </n-card>
                </n-gi>
                <n-gi>
                    <n-card embedded size="small">
                        <n-statistic label="最高赛季攻城" :value="bookData.max_season_siege || 0" />
                    </n-card>
                </n-gi>
                <n-gi>
                    <n-card embedded size="small">
                        <n-statistic label="最高赛季拆除" :value="bookData.max_season_demolish || 0" />
                    </n-card>
                </n-gi>
                <n-gi>
                    <n-card embedded size="small">
                        <n-statistic label="访客数" :value="bookData.visitors || 0" />
                    </n-card>
                </n-gi>
                <n-gi>
                    <n-card embedded size="small">
                        <n-statistic label="点赞数" :value="bookData.likes || 0" />
                    </n-card>
                </n-gi>
            </n-grid>

            <n-card v-if="bookData.best_team_hero_id" class="page-card best-team-card" embedded>
                <div class="best-team-header">
                    <span class="best-team-title">最高赛季灭敌队伍</span>
                </div>
                <div class="best-team-content">
                    <div class="best-team-hero">
                        <div class="best-team-hero-img">
                            <img :src="`https://g0.gph.netease.com/ngsocial/community/stzb/cn/cards/cut/card_small_${getHeroIconId(bookData.best_team_hero_id)}.jpg?gameid=g10`"
                                @error="$event.target.style.display='none'" />
                        </div>
                        <div class="best-team-hero-info">
                            <span class="best-team-hero-name">{{ bookData.best_team_hero_name || getHeroName(bookData.best_team_hero_id) }}</span>
                            <span class="best-team-hero-role">大营</span>
                        </div>
                    </div>
                    <n-statistic label="灭敌数量" :value="formatNum(bookData.best_team_kills)" />
                </div>
            </n-card>

            <n-card class="page-card" embedded>
                <n-collapse>
                    <n-collapse-item title="原始数据（调试）" name="raw">
                        <pre class="raw-data">{{ JSON.stringify(bookData.raw, null, 2) }}</pre>
                    </n-collapse-item>
                </n-collapse>
            </n-card>
        </template>

        <n-empty v-else description="暂无数据，请开启开关后在游戏中打开主公簿" style="padding: 80px 0;" />
    </div>
</template>

<style scoped lang="scss">
.page-book {
    display: flex;
    flex-direction: column;
    gap: 20px;
}

.page-card {
    border-radius: 12px;
}

.page-header {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
}

.page-header-info {
    flex: 1;
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

.switch-wrap {
    display: flex;
    align-items: center;
    gap: 10px;
}

.switch-label {
    font-size: 13px;
    color: var(--color-text-secondary);
}

.hero-card {
    background: var(--color-hero-bg);
    color: var(--color-hero-text);
}

.hero-info {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 24px;
}

.hero-name {
    font-size: 22px;
    font-weight: 700;
    margin-bottom: 8px;
}

.hero-meta {
    display: flex;
    gap: 16px;
    font-size: 14px;
    opacity: 0.85;
}

.hero-stats {
    display: flex;
    gap: 32px;

    :deep(.n-statistic .n-statistic__label) {
        color: rgba(255, 255, 255, 0.85);
    }

    :deep(.n-statistic .n-statistic-value__content) {
        color: #ffffff;
    }

    :deep(.n-statistic .n-statistic-value__suffix) {
        color: rgba(255, 255, 255, 0.85);
    }
}

.stat-grid {
    margin-top: 0;
}

.best-team-card {
    border-radius: 12px;
}

.best-team-header {
    margin-bottom: 16px;
}

.best-team-title {
    font-size: 16px;
    font-weight: 600;
    color: var(--color-text);
}

.best-team-content {
    display: flex;
    align-items: center;
    justify-content: space-between;
}

.best-team-hero {
    display: flex;
    align-items: center;
    gap: 12px;
}

.best-team-hero-img {
    width: 56px;
    height: 56px;
    border-radius: 8px;
    overflow: hidden;
    background: var(--color-surface-hover);
    flex-shrink: 0;

    img {
        width: 100%;
        height: 100%;
        object-fit: cover;
    }
}

.best-team-hero-info {
    display: flex;
    flex-direction: column;
    gap: 2px;
}

.best-team-hero-name {
    font-size: 15px;
    font-weight: 600;
    color: var(--color-text);
}

.best-team-hero-role {
    font-size: 12px;
    color: var(--color-text-secondary);
}

.raw-data {
    font-size: 12px;
    line-height: 1.5;
    color: var(--color-text-secondary);
    background: var(--color-bg);
    padding: 12px;
    border-radius: 6px;
    overflow-x: auto;
    max-height: 400px;
    white-space: pre-wrap;
    word-break: break-all;
}
</style>
