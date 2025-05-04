import { createApp } from 'vue'
import { createPinia } from 'pinia';
import './bikamoeapp.scss'
import App from './App.vue'
import router from './routes';
import piniaPluginPersistedstate from 'pinia-plugin-persistedstate';

// createApp(App).mount('#app')

const pinia = createPinia();
pinia.use(piniaPluginPersistedstate);
createApp(App).use(router).use(pinia).mount('#app');
