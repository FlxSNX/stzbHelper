<script setup>
import {NMessageProvider,NDialogProvider} from 'naive-ui'
import { NConfigProvider } from 'naive-ui'
import { zhCN, dateZhCN } from 'naive-ui'
</script>

<template>
	<n-config-provider :locale="zhCN" :date-locale="dateZhCN">
	<n-dialog-provider>
	<n-message-provider>
		<router-view v-slot="{ Component, route }">
			<!-- <transition> -->
			<keep-alive include="Index">
				<component :is="Component" :key="route.path"></component>
			</keep-alive>
			<!-- </transition> -->
		</router-view>
	</n-message-provider>
	</n-dialog-provider>
	</n-config-provider>
</template>

<style scoped></style>
