<script setup>
import { ref, computed, h, watch } from 'vue'
import { NCard, NButton, NInput, NInputNumber, NEmpty, NSpin, NTag, NPagination, NDataTable, useMessage } from 'naive-ui'
import { GetTeamWinRate, GetTeamWinRateByTeam } from '../../wailsjs/go/main/App'
import { Search, Swords, Image, Table, Users, Layers } from 'lucide-vue-next'
import { herocfg, skillcfg } from '../cfg'

const heroMap = JSON.parse(herocfg)
const skillMap = JSON.parse(skillcfg)

const nmessage = useMessage()
const loading = ref(false)
const results = ref([])

const searchName = ref('')
const searchUnion = ref('')
const searchIdu = ref('')
const minLevel = ref(30)
const minHp = ref(20000)

const hasSearched = ref(false)
const useImageMode = ref(true)
const groupByPlayer = ref(true)
const page = ref(1)
const pageSize = ref(50)
const total = ref(0)

const doSearch = (newPage) => {
    if (typeof newPage === 'number') page.value = newPage
    else page.value = 1
    loading.value = true
    results.value = []
    hasSearched.value = true
    const apiFn = groupByPlayer.value ? GetTeamWinRate : GetTeamWinRateByTeam
    apiFn(searchName.value, searchUnion.value, searchIdu.value, page.value, pageSize.value, minLevel.value, minHp.value).then(v => {
        let resp = JSON.parse(v)
        if (resp.code == 200) {
            results.value = resp.data.list || []
            total.value = resp.data.total || 0
        } else {
            nmessage.error(resp.msg)
        }
    }).catch(e => {
        nmessage.error('查询失败: ' + e)
    }).finally(() => {
        loading.value = false
    })
}

const resolveHeroId = (id) => {
    if (!id) return id
    const num = Number(id)
    return num >= 130000 ? num - 30000 : num
}

const getHeroIconId = (id) => {
    if (!id) return id
    const hero = heroMap[String(resolveHeroId(id))]
    return hero ? hero.iconId : id
}

const getHeroName = (id) => {
    if (!id) return ''
    const hero = heroMap[String(resolveHeroId(id))]
    return hero ? hero.name : `未知(${id})`
}

const getHeroCountry = (id) => {
    if (!id) return ''
    const hero = heroMap[String(resolveHeroId(id))]
    return hero ? hero.country : ''
}

const getHeroType = (id) => {
    if (!id) return ''
    const hero = heroMap[String(resolveHeroId(id))]
    return hero ? hero.type : ''
}

const getHeroQuality = (id) => {
    if (!id) return 5
    const hero = heroMap[String(resolveHeroId(id))]
    return hero ? hero.quality : 5
}

const getSkillName = (id) => {
    if (!id) return ''
    const skill = skillMap[String(id)]
    return skill ? skill.name : `未知(${id})`
}

const getSkillQuality = (id) => {
    if (!id) return ''
    const skill = skillMap[String(id)]
    return skill ? skill.zfQuality : ''
}

const getSkillType = (id) => {
    if (!id) return ''
    const skill = skillMap[String(id)]
    return skill ? skill.type : ''
}

const parseSkillInfo = (str, role) => {
    if (!str) return []
    const groups = String(str).split(';').filter(s => s.trim() !== '')
    const parsed = groups.map(g => {
        const parts = g.split(',')
        return {
            index: parseInt(parts[0]),
            skills: [
                { id: parts[1], level: parseInt(parts[2]) },
                { id: parts[3], level: parseInt(parts[4]) },
                { id: parts[5], level: parseInt(parts[6]) },
            ]
        }
    })
    let filtered = role === 'defend'
        ? parsed.filter(g => g.index >= 4 && g.index <= 6)
        : parsed.filter(g => g.index >= 1 && g.index <= 3)
    if (role === 'defend') filtered.reverse()
    return filtered
}

const formatTime = (ts) => {
    if (!ts) return ''
    const d = new Date(ts * 1000)
    const pad = (n) => String(n).padStart(2, '0')
    return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}

const winRateColor = (rate) => {
    if (rate >= 70) return '#22c55e'
    if (rate >= 50) return '#3b82f6'
    if (rate >= 30) return '#f59e0b'
    return '#ef4444'
}

const lossRateColor = (rate) => {
    if (rate >= 70) return '#ef4444'
    if (rate >= 50) return '#f59e0b'
    if (rate >= 30) return '#3b82f6'
    return '#22c55e'
}

const groupedResults = computed(() => {
    const map = {}
    results.value.forEach(r => {
        if (!map[r.player_name]) map[r.player_name] = []
        map[r.player_name].push(r)
    })
    return map
})

const tableData = computed(() => results.value)

watch(groupByPlayer, () => {
    page.value = 1
    results.value = []
    if (hasSearched.value) doSearch()
})

const playerColumns = [
    {
        title: '玩家',
        key: 'player_name',
        width: 120,
        ellipsis: { tooltip: true },
    },
    {
        title: '队伍',
        key: 'heroes',
        minWidth: 200,
        render(row) {
            return [1, 2, 3].map(i => `${getHeroName(row[`hero${i}_id`])}(${getHeroType(row[`hero${i}_id`])})`).join(' / ')
        }
    },
    {
        title: '战法',
        key: 'skills',
        width: 200,
        render(row) {
            const groups = parseSkillInfo(row.all_skill_info, row.role)
            const lines = groups.map(g => {
                const names = g.skills.slice(1).filter(s => s.id && s.id !== '0').map(s => getSkillName(s.id))
                return names.join('/')
            }).filter(Boolean)
            return h('div', { style: { whiteSpace: 'pre-line' } }, lines.join('\n'))
        }
    },
    {
        title: '红度',
        key: 'total_star',
        width: 60,
        align: 'center',
        sorter: (a, b) => a.total_star - b.total_star,
    },
    {
        title: '总场次',
        key: 'total_battles',
        width: 80,
        align: 'center',
        sorter: (a, b) => a.total_battles - b.total_battles,
    },
    {
        title: '胜',
        key: 'win_count',
        width: 60,
        align: 'center',
        sorter: (a, b) => a.win_count - b.win_count,
    },
    {
        title: '负',
        key: 'loss_count',
        width: 60,
        align: 'center',
        sorter: (a, b) => a.loss_count - b.loss_count,
    },
    {
        title: '平',
        key: 'draw_count',
        width: 60,
        align: 'center',
        sorter: (a, b) => a.draw_count - b.draw_count,
    },
    {
        title: '胜率',
        key: 'win_rate',
        width: 80,
        align: 'center',
        sorter: (a, b) => a.win_rate - b.win_rate,
        defaultSortOrder: 'descend',
        render(row) {
            const rate = row.win_rate || 0
            return h('span', { style: { color: winRateColor(rate), fontWeight: '700' } }, rate + '%')
        }
    },
    {
        title: '负率',
        key: 'loss_rate',
        width: 80,
        align: 'center',
        render(row) {
            const rate = row.total_battles > 0 ? Math.round(row.loss_count / row.total_battles * 1000) / 10 : 0
            return h('span', { style: { color: lossRateColor(rate), fontWeight: '700' } }, rate + '%')
        }
    },
    {
        title: '平局率',
        key: 'draw_rate',
        width: 80,
        align: 'center',
        render(row) {
            const rate = row.total_battles > 0 ? Math.round(row.draw_count / row.total_battles * 1000) / 10 : 0
            return h('span', { style: { fontWeight: '700' } }, rate + '%')
        }
    },
    {
        title: '最近战斗',
        key: 'last_time',
        width: 140,
        render(row) {
            return formatTime(row.last_time)
        }
    },
    {
        title: '队伍ID',
        key: 'idu',
        width: 100,
        ellipsis: { tooltip: true },
    },
]

const teamColumns = [
    {
        title: '队伍',
        key: 'heroes',
        minWidth: 200,
        render(row) {
            return [1, 2, 3].map(i => getHeroName(row[`hero${i}_id`])).join(' / ')
        }
    },
    {
        title: '玩家',
        key: 'players',
        width: 160,
        ellipsis: { tooltip: true },
    },
    {
        title: '战法',
        key: 'skills',
        width: 200,
        render(row) {
            const groups = parseSkillInfo(row.all_skill_info, row.role)
            const lines = groups.map(g => {
                const names = g.skills.slice(1).filter(s => s.id && s.id !== '0').map(s => getSkillName(s.id))
                return names.join('/')
            }).filter(Boolean)
            return h('div', { style: { whiteSpace: 'pre-line' } }, lines.join('\n'))
        }
    },
    {
        title: '总场次',
        key: 'total_battles',
        width: 80,
        align: 'center',
        sorter: (a, b) => a.total_battles - b.total_battles,
    },
    {
        title: '胜',
        key: 'win_count',
        width: 60,
        align: 'center',
        sorter: (a, b) => a.win_count - b.win_count,
    },
    {
        title: '负',
        key: 'loss_count',
        width: 60,
        align: 'center',
        sorter: (a, b) => a.loss_count - b.loss_count,
    },
    {
        title: '平',
        key: 'draw_count',
        width: 60,
        align: 'center',
        sorter: (a, b) => a.draw_count - b.draw_count,
    },
    {
        title: '胜率',
        key: 'win_rate',
        width: 80,
        align: 'center',
        sorter: (a, b) => a.win_rate - b.win_rate,
        defaultSortOrder: 'descend',
        render(row) {
            const rate = row.win_rate || 0
            return h('span', { style: { color: winRateColor(rate), fontWeight: '700' } }, rate + '%')
        }
    },
    {
        title: '负率',
        key: 'loss_rate',
        width: 80,
        align: 'center',
        render(row) {
            const rate = row.total_battles > 0 ? Math.round(row.loss_count / row.total_battles * 1000) / 10 : 0
            return h('span', { style: { color: lossRateColor(rate), fontWeight: '700' } }, rate + '%')
        }
    },
    {
        title: '平局率',
        key: 'draw_rate',
        width: 80,
        align: 'center',
        render(row) {
            const rate = row.total_battles > 0 ? Math.round(row.draw_count / row.total_battles * 1000) / 10 : 0
            return h('span', { style: { fontWeight: '700' } }, rate + '%')
        }
    },
    {
        title: '最近战斗',
        key: 'last_time',
        width: 140,
        render(row) {
            return formatTime(row.last_time)
        }
    },
]

const currentColumns = computed(() => groupByPlayer.value ? playerColumns : teamColumns)
</script>

<template>
    <div class="page-team-winrate">
        <n-card class="page-card" embedded>
            <div class="page-header">
                <div class="page-header-info">
                    <h2 class="page-title">队伍胜率</h2>
                    <p class="page-desc">统计队伍组合的战斗胜率（攻方胜: 1/2/3/4/10/18/19，攻方负: 0，平局: 6/7/8/13）</p>
                </div>
            </div>

            <div class="search-bar">
                <n-input v-model:value="searchName" placeholder="玩家名称" clearable @keyup.enter="doSearch" />
                <n-input v-model:value="searchUnion" placeholder="同盟名称" clearable @keyup.enter="doSearch" />
                <n-input v-model:value="searchIdu" placeholder="队伍 ID" clearable @keyup.enter="doSearch" />
                <n-button type="primary" @click="doSearch()" :loading="loading">
                    <template #icon><Search :size="16" /></template>
                    查询
                </n-button>
                <n-button quaternary :type="groupByPlayer ? 'primary' : 'default'" @click="groupByPlayer = !groupByPlayer">
                    <template #icon><Users :size="16" /></template>
                    {{ groupByPlayer ? '按玩家' : '按队伍' }}
                </n-button>
                <n-button quaternary :type="useImageMode ? 'primary' : 'default'" @click="useImageMode = !useImageMode">
                    <template #icon><Image :size="16" /></template>
                    {{ useImageMode ? '图片' : '表格' }}
                </n-button>
            </div>
            <div class="filter-bar">
                <div class="filter-item">
                    <span class="filter-label">最低等级</span>
                    <n-input-number v-model:value="minLevel" :min="1" :max="50" style="width: 100px;" />
                </div>
                <div class="filter-item">
                    <span class="filter-label">最低兵力</span>
                    <n-input-number v-model:value="minHp" :min="0" :max="99999" :step="1000" style="width: 130px;" />
                </div>
            </div>

            <div class="pagination-wrap" v-if="total > pageSize">
                <n-pagination
                    v-model:page="page"
                    :page-size="pageSize"
                    :item-count="total"
                    :on-update:page="(p) => doSearch(p)"
                />
            </div>

            <div class="result-area" v-if="loading">
                <div class="loading-wrap">
                    <n-spin size="medium" />
                    <span>查询中...</span>
                </div>
            </div>

            <div class="result-area" v-else-if="hasSearched && results.length === 0">
                <n-empty description="未找到队伍数据" style="padding: 60px 0;" />
            </div>

            <div class="result-area" v-else-if="results.length > 0">
                <div class="result-summary">
                    共 <strong>{{ total }}</strong> {{ groupByPlayer ? '条记录' : '支队伍' }}
                </div>

                <!-- 表格模式 -->
                <template v-if="!useImageMode">
                    <n-data-table
                        :columns="currentColumns"
                        :data="tableData"
                        :bordered="false"
                        size="small"
                        :scroll-x="groupByPlayer ? 1320 : 1220"
                    />
                </template>

                <!-- 大图模式 按玩家 -->
                <template v-else-if="groupByPlayer">
                    <div class="player-section" v-for="(teams, playerName) in groupedResults" :key="playerName">
                        <div class="player-name">
                            <Swords :size="16" />
                            {{ playerName }}
                        </div>

                        <div class="team-card" v-for="team in teams" :key="playerName + team.hero1_id + team.hero2_id + team.hero3_id">
                            <div class="team-header">
                                <span class="team-idu">{{ team.player_name }} · ID {{ team.idu }}</span>
                                <span class="team-time">{{ formatTime(team.last_time) }}</span>
                            </div>

                            <div class="hero-row hero-row--big">
                                <div class="hero-big" v-for="i in 3" :key="i">
                                    <div class="hero-big-img">
                                        <img v-if="team[`hero${i}_id`]"
                                            :src="`https://cbg-stzb.res.netease.com/game_res/cards/cut/card_medium_${getHeroIconId(team[`hero${i}_id`])}.jpg`"
                                            @error="$event.target.style.display='none'" />
                                        <div class="hero-placeholder" v-else>?</div>
                                        <div class="hero-big-stars">
                                            <img v-for="s in team[`hero${i}_star`]" :key="'r'+s" class="hero-big-star-img" :src="`data:image/png;base64,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`" />
                                            <img v-for="s in (getHeroQuality(team[`hero${i}_id`]) - team[`hero${i}_star`])" :key="'n'+s" class="hero-big-star-img" :src="`data:image/png;base64,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`" />
                                            </div>
                                        </div>
                                    <div class="hero-big-info">
                                        <div class="hero-big-header">
                                            <span class="hero-big-name">{{ getHeroName(team[`hero${i}_id`]) }}</span>
                                            <span class="hero-big-meta">
                                                <span v-if="getHeroCountry(team[`hero${i}_id`])">{{ getHeroCountry(team[`hero${i}_id`]) }}</span>
                                                <span v-if="getHeroType(team[`hero${i}_id`])">·{{ getHeroType(team[`hero${i}_id`]) }}</span>
                                                <span class="hero-big-star">·{{ team[`hero${i}_star`] }}红</span>
                                            </span>
                                        </div>
                                        <div class="hero-big-skills" v-if="team.all_skill_info">
                                            <div class="hero-big-skill" v-for="(skill, si) in (parseSkillInfo(team.all_skill_info, team.role)[i - 1]?.skills || [])" :key="si">
                                                <template v-if="skill && skill.id && skill.id !== '0'">
                                                    <n-tag v-if="getSkillQuality(skill.id)" size="tiny" :bordered="false" :type="getSkillQuality(skill.id) === 'S' ? 'warning' : getSkillQuality(skill.id) === 'A' ? 'info' : 'default'">{{ getSkillQuality(skill.id) }}</n-tag>
                                                    <n-tag v-if="getSkillType(skill.id)" size="tiny" :bordered="false">{{ getSkillType(skill.id) }}</n-tag>
                                                    <span class="hero-big-skill-name">{{ getSkillName(skill.id) }}</span>
                                                </template>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </div>

                            <div class="stats-bar">
                                <div class="stat-item">
                                    <span class="stat-label">总场次</span>
                                    <span class="stat-value">{{ team.total_battles }}</span>
                                </div>
                                <div class="stat-item stat-win">
                                    <span class="stat-label">胜</span>
                                    <span class="stat-value">{{ team.win_count }}</span>
                                </div>
                                <div class="stat-item stat-loss">
                                    <span class="stat-label">负</span>
                                    <span class="stat-value">{{ team.loss_count }}</span>
                                </div>
                                <div class="stat-item stat-draw">
                                    <span class="stat-label">平</span>
                                    <span class="stat-value">{{ team.draw_count }}</span>
                                </div>
                                <div class="stat-item">
                                    <span class="stat-label">胜率</span>
                                    <span class="stat-value" :style="{ color: winRateColor(team.win_rate || 0), fontWeight: 700 }">{{ (team.win_rate || 0) }}%</span>
                                </div>
                                <div class="stat-item">
                                    <span class="stat-label">负率</span>
                                    <span class="stat-value" :style="{ color: lossRateColor(team.total_battles > 0 ? (team.loss_count / team.total_battles * 100) : 0), fontWeight: 700 }">{{ team.total_battles > 0 ? (team.loss_count / team.total_battles * 100).toFixed(1) : 0 }}%</span>
                                </div>
                                <div class="stat-item">
                                    <span class="stat-label">平局率</span>
                                    <span class="stat-value" :style="{ fontWeight: 700 }">{{ team.total_battles > 0 ? (team.draw_count / team.total_battles * 100).toFixed(1) : 0 }}%</span>
                                </div>
                            </div>
                        </div>
                    </div>
                </template>

                <!-- 大图模式 按队伍组合 -->
                <template v-else>
                    <div class="team-card" v-for="(team, key) in results" :key="key">
                        <div class="team-header">
                            <span class="team-idu" v-if="team.players">玩家: {{ team.players }}</span>
                            <span class="team-time">{{ formatTime(team.last_time) }}</span>
                        </div>

                        <div class="hero-row hero-row--big">
                            <div class="hero-big" v-for="i in 3" :key="i">
                                <div class="hero-big-img">
                                    <img v-if="team[`hero${i}_id`]"
                                        :src="`https://cbg-stzb.res.netease.com/game_res/cards/cut/card_medium_${getHeroIconId(team[`hero${i}_id`])}.jpg`"
                                        @error="$event.target.style.display='none'" />
                                    <div class="hero-placeholder" v-else>?</div>
                                </div>
                                <div class="hero-big-info">
                                    <div class="hero-big-header">
                                        <span class="hero-big-name">{{ getHeroName(team[`hero${i}_id`]) }}</span>
                                        <span class="hero-big-meta">
                                            <span v-if="getHeroCountry(team[`hero${i}_id`])">{{ getHeroCountry(team[`hero${i}_id`]) }}</span>
                                            <span v-if="getHeroType(team[`hero${i}_id`])">·{{ getHeroType(team[`hero${i}_id`]) }}</span>
                                        </span>
                                    </div>
                                    <div class="hero-big-skills" v-if="team.all_skill_info">
                                        <div class="hero-big-skill" v-for="(skill, si) in (parseSkillInfo(team.all_skill_info, team.role)[i - 1]?.skills || [])" :key="si">
                                            <template v-if="skill && skill.id && skill.id !== '0'">
                                                <n-tag v-if="getSkillQuality(skill.id)" size="tiny" :bordered="false" :type="getSkillQuality(skill.id) === 'S' ? 'warning' : getSkillQuality(skill.id) === 'A' ? 'info' : 'default'">{{ getSkillQuality(skill.id) }}</n-tag>
                                                <n-tag v-if="getSkillType(skill.id)" size="tiny" :bordered="false">{{ getSkillType(skill.id) }}</n-tag>
                                                <span class="hero-big-skill-name">{{ getSkillName(skill.id) }}</span>
                                            </template>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>

                        <div class="stats-bar">
                            <div class="stat-item">
                                <span class="stat-label">总场次</span>
                                <span class="stat-value">{{ team.total_battles }}</span>
                            </div>
                            <div class="stat-item stat-win">
                                <span class="stat-label">胜</span>
                                <span class="stat-value">{{ team.win_count }}</span>
                            </div>
                            <div class="stat-item stat-loss">
                                <span class="stat-label">负</span>
                                <span class="stat-value">{{ team.loss_count }}</span>
                            </div>
                            <div class="stat-item stat-draw">
                                <span class="stat-label">平</span>
                                <span class="stat-value">{{ team.draw_count }}</span>
                            </div>
                            <div class="stat-item">
                                <span class="stat-label">胜率</span>
                                <span class="stat-value" :style="{ color: winRateColor(team.win_rate || 0), fontWeight: 700 }">{{ (team.win_rate || 0) }}%</span>
                            </div>
                            <div class="stat-item">
                                <span class="stat-label">负率</span>
                                <span class="stat-value" :style="{ color: lossRateColor(team.total_battles > 0 ? (team.loss_count / team.total_battles * 100) : 0), fontWeight: 700 }">{{ team.total_battles > 0 ? (team.loss_count / team.total_battles * 100).toFixed(1) : 0 }}%</span>
                            </div>
                            <div class="stat-item">
                                <span class="stat-label">平局率</span>
                                <span class="stat-value" :style="{ fontWeight: 700 }">{{ team.total_battles > 0 ? (team.draw_count / team.total_battles * 100).toFixed(1) : 0 }}%</span>
                            </div>
                        </div>
                    </div>
                </template>

            </div>
        </n-card>
    </div>
</template>

<style scoped lang="scss">
.page-team-winrate {
    display: flex;
    flex-direction: column;
}

.page-card {
    border-radius: 12px;
    overflow: hidden;
}

.page-header {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    margin-bottom: 20px;
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

.search-bar {
    display: flex;
    gap: 12px;
    margin-bottom: 20px;
    flex-wrap: wrap;
}

.search-bar .n-input {
    flex: 1;
    min-width: 160px;
}

.filter-bar {
    display: flex;
    gap: 20px;
    margin-bottom: 20px;
    flex-wrap: wrap;
}

.filter-item {
    display: flex;
    align-items: center;
    gap: 8px;
}

.filter-label {
    font-size: 13px;
    color: var(--color-text-secondary);
    white-space: nowrap;
}

.result-area {
    min-height: 200px;
    overflow-x: auto;
}

.loading-wrap {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 12px;
    padding: 60px 0;
    color: var(--color-text-secondary);
    font-size: 14px;
}

.result-summary {
    font-size: 13px;
    color: var(--color-text-secondary);
    margin-bottom: 16px;
}

.pagination-wrap {
    display: flex;
    justify-content: center;
    margin-top: 20px;
    padding: 16px 0;
}

.player-section {
    margin-bottom: 24px;
}

.player-name {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 16px;
    font-weight: 700;
    color: var(--color-text);
    margin-bottom: 12px;
    padding-bottom: 8px;
    border-bottom: 2px solid var(--color-border);
}

.team-card {
    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: 10px;
    padding: 16px;
    margin-bottom: 12px;
    transition: box-shadow 0.2s;

    &:hover {
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.06);
    }
}

.team-header {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 12px;
    flex-wrap: wrap;
}

.team-idu {
    display: flex;
    align-items: center;
    gap: 4px;
    font-size: 13px;
    color: var(--color-text-secondary);
}

.team-time {
    margin-left: auto;
    font-size: 12px;
    color: var(--color-text-secondary);
}

.hero-row {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 16px;
    margin-bottom: 16px;
}

.hero-big {
    display: flex;
    gap: 0;
    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: 10px;
    overflow: hidden;

    &-img {
        position: relative;
        width: 160px;
        flex-shrink: 0;

        img {
            width: 100%;
            height: 100%;
            object-fit: cover;
            display: block;
        }
    }

    &-stars {
        position: absolute;
        top: 6px;
        right: 6px;
        display: flex;
        gap: 1px;
        padding: 3px 4px;
    }

    &-star-img {
        width: 18px !important;
        height: 18px !important;
    }

    &-info {
        flex: 1;
        min-width: 0;
        padding: 10px 12px;
        display: flex;
        flex-direction: column;
        gap: 6px;
    }

    &-header {
        display: flex;
        flex-direction: column;
        gap: 2px;
    }

    &-name {
        font-size: 14px;
        font-weight: 700;
        color: var(--color-text);
    }

    &-meta {
        display: flex;
        align-items: center;
        gap: 2px;
        font-size: 11px;
        color: var(--color-text-secondary);
        flex-wrap: wrap;
    }

    &-star {
        color: #f59e0b;
    }

    &-skills {
        display: flex;
        flex-direction: column;
        gap: 4px;
        margin-top: 4px;
    }

    &-skill {
        display: flex;
        align-items: center;
        gap: 4px;
        font-size: 11px;
    }

    &-skill-name {
        color: var(--color-text);
        font-size: 11px;
    }
}

.hero-placeholder {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 18px;
    color: var(--color-text-secondary);
}

.stats-bar {
    display: flex;
    align-items: center;
    gap: 24px;
    padding: 12px 16px;
    background: var(--color-bg);
    border-radius: 8px;
    flex-wrap: wrap;
}

.stat-item {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 2px;
}

.stat-label {
    font-size: 12px;
    color: var(--color-text-secondary);
}

.stat-value {
    font-size: 16px;
    font-weight: 600;
    color: var(--color-text);
}

.stat-win .stat-value {
    color: #22c55e;
}

.stat-loss .stat-value {
    color: #ef4444;
}

.stat-draw .stat-value {
    color: #f59e0b;
}
</style>
