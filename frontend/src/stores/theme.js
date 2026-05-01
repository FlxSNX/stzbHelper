import { defineStore } from 'pinia'
import { ref, watch } from 'vue'

export const useThemeStore = defineStore('theme', () => {
    const isDark = ref(false)

    const applyTheme = (dark) => {
        document.documentElement.setAttribute('data-theme', dark ? 'dark' : 'light')
    }

    const toggle = () => {
        isDark.value = !isDark.value
    }

    watch(isDark, (val) => {
        applyTheme(val)
    }, { immediate: true })

    return { isDark, toggle }
}, {
    persist: true,
})
