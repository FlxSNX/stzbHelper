import {createRouter,createWebHashHistory} from 'vue-router';
import Index from './pages/Index.vue';
import TeamUser from './pages/TeamUser.vue';
import Task from './pages/Task.vue';
import GroupWu from './pages/GroupWu.vue';

const routes = [
    { 
        path: '/', 
        component: Index ,
        keepalive: true
    },
    { 
        path: '/teamuser', 
        component: TeamUser 
    },
    { 
        path: '/task', 
        component: Task 
    },
    {
        path: '/groupWu',
        component: GroupWu
    }
]

const router = createRouter({
    history: createWebHashHistory(),
    routes,
})

export default router;