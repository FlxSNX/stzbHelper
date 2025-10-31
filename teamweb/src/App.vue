<script setup>
import {NMessageProvider} from 'naive-ui'
</script>

<template>
	<n-message-provider>
		<router-view v-slot="{ Component, route }">
			<!-- <transition> -->
			<keep-alive include="Index">
				<component :is="Component" :key="route.path"></component>
			</keep-alive>
			<!-- </transition> -->
		</router-view>
	</n-message-provider>
</template>

<style scoped></style>
