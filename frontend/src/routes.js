import { createRouter, createWebHashHistory } from 'vue-router';
import Index from './pages/Index.vue';
import TeamUser from './pages/TeamUser.vue';
import Task from './pages/Task.vue';
import GroupWu from './pages/GroupWu.vue';
import SelectDb from './pages/SelectDb.vue';
import Logs from './pages/Logs.vue';
import NpcapHelp from './pages/NpcapHelp.vue';
import Debug from './pages/Debug.vue';
import TeamQuery from './pages/TeamQuery.vue';
import Book from './pages/Book.vue';
// import BattleCall from './pages/BattleCall.vue';
import TeamWinRate from './pages/TeamWinRate.vue';
import { CheckNpcap } from '../wailsjs/go/main/App';

const routes = [
    {
        path: '/',
        name: 'index',
        component: Index,
        meta: { title: '控制面板' }
    },
    {
        path: '/teamuser',
        name: 'teamuser',
        component: TeamUser,
        meta: { title: '同盟成员' }
    },
    {
        path: '/task',
        name: 'task',
        component: Task,
        meta: { title: '攻城任务' }
    },
    {
        path: '/groupWu',
        name: 'groupWu',
        component: GroupWu,
        meta: { title: '分组武勋' }
    },
    {
        path: '/select-db',
        name: 'selectDb',
        component: SelectDb,
        meta: { title: '选择数据库' }
    },
    {
        path: '/logs',
        name: 'logs',
        component: Logs,
        meta: { title: '运行日志' }
    },
    {
        path: '/npcap-help',
        name: 'npcapHelp',
        component: NpcapHelp,
        meta: { title: '安装 Npcap' }
    },
    {
        path: '/debug',
        name: 'debug',
        component: Debug,
        meta: { title: 'API 调试' }
    },
    {
        path: '/teamquery',
        name: 'teamquery',
        component: TeamQuery,
        meta: { title: '队伍查询' }
    },
    {
        path: '/teamwinrate',
        name: 'teamwinrate',
        component: TeamWinRate,
        meta: { title: '队伍胜率' }
    },
    {
        path: '/book',
        name: 'book',
        component: Book,
        meta: { title: '主公簿' }
    },
    // {
    //     path: '/battlecall',
    //     name: 'battlecall',
    //     component: BattleCall,
    //     meta: { title: '战役叫阵' }
    // }
]

const router = createRouter({
    history: createWebHashHistory(),
    routes,
})

let npcapChecked = false
let npcapInstalled = false

router.beforeEach(async (to, from, next) => {
    if (to.name === 'npcapHelp') {
        next()
        return
    }

    if (!npcapChecked) {
        try {
            const resp = JSON.parse(await CheckNpcap())
            npcapInstalled = resp.code == 200 && resp.data.installed
        } catch {
            npcapInstalled = false
        }
        npcapChecked = true
    }

    if (!npcapInstalled) {
        next({ name: 'npcapHelp' })
        return
    }

    const selectedDb = sessionStorage.getItem('selectedDb')
    if (to.name !== 'selectDb' && !selectedDb) {
        next({ name: 'selectDb' })
    } else {
        next()
    }
})

export default router;
