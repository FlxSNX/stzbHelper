<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { WindowMinimise, WindowToggleMaximise, WindowIsMaximised, Quit } from '../../wailsjs/runtime/runtime'
import { Swords, Minus, Square, X, Copy } from 'lucide-vue-next'
import { useThemeStore } from '../stores/theme'

const themeStore = useThemeStore()

const isMaximised = ref(false)

const checkMaximised = async () => {
    try {
        isMaximised.value = await WindowIsMaximised()
    } catch {
        isMaximised.value = false
    }
}

const handleMinimise = () => {
    WindowMinimise()
}

const handleToggleMaximise = () => {
    WindowToggleMaximise()
    setTimeout(checkMaximised, 100)
}

const handleClose = () => {
    Quit()
}

const handleDblClick = () => {
    handleToggleMaximise()
}

onMounted(() => {
    window.addEventListener('resize', checkMaximised)
    checkMaximised()
})

onUnmounted(() => {
    window.removeEventListener('resize', checkMaximised)
})
</script>

<template>
    <div class="titlebar">
        <div class="titlebar-left" @dblclick="handleDblClick">
            <div class="titlebar-icon" :class="{ 'titlebar-icon--dark': themeStore.isDark }">
                <Swords :size="16" />
            </div>
            <span class="titlebar-text">率土之滨助手</span>
        </div>
        <div class="titlebar-controls">
            <button class="titlebar-btn" @click="handleMinimise" title="最小化">
                <Minus :size="14" />
            </button>
            <button class="titlebar-btn" @click="handleToggleMaximise" :title="isMaximised ? '还原' : '最大化'">
                <Copy v-if="isMaximised" :size="13" />
                <Square v-else :size="12" />
            </button>
            <button class="titlebar-btn titlebar-btn--close" @click="handleClose" title="关闭">
                <X :size="14" />
            </button>
        </div>
    </div>
</template>

<style scoped lang="scss">
.titlebar {
    display: flex;
    align-items: center;
    justify-content: space-between;
    height: 36px;
    background: var(--color-surface);
    border-bottom: 1px solid var(--color-border);
    user-select: none;
    flex-shrink: 0;
    --wails-draggable: drag;
}

.titlebar-left {
    display: flex;
    align-items: center;
    gap: 10px;
    padding-left: 14px;
    pointer-events: none;
}

.titlebar-icon {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 24px;
    height: 24px;
    border-radius: 6px;
    background: linear-gradient(135deg, #404040, #171717);
    color: #fff;
    flex-shrink: 0;

    &--dark {
        background: linear-gradient(135deg, #e5e5e5, #a3a3a3);
        color: #0a0a0a;
    }
}

.titlebar-text {
    font-size: 13px;
    font-weight: 600;
    color: var(--color-text);
    white-space: nowrap;
}

.titlebar-controls {
    display: flex;
    align-items: center;
    height: 100%;
    --wails-draggable: no-drag;
}

.titlebar-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 46px;
    height: 100%;
    border: none;
    background: transparent;
    color: var(--color-text-secondary);
    cursor: pointer;
    transition: background-color 0.15s, color 0.15s;
    outline: none;
    --wails-draggable: no-drag;

    &:hover {
        background-color: var(--color-surface-hover);
        color: var(--color-text);
    }

    &--close:hover {
        background-color: #ef4444;
        color: #ffffff;
    }
}
</style>
