import {createRouter,createWebHashHistory} from 'vue-router';
import Index from './pages/Index.vue';
import Team from './pages/Team.vue';

const routes = [
    { 
        path: '/', 
        component: Index ,
        keepalive: true
    },
    { 
        path: '/team', 
        component: Team 
    }
]

const router = createRouter({
    history: createWebHashHistory(),
    routes,
})

export default router;