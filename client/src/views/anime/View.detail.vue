<template>
	<div class="bg-white shadow rounded-lg p-6">
		<h2 class="text-2xl font-semibold mb-4">{{ animeStore.item?.name || "name" }}</h2>

		<div v-if="animeStore.loading" class="text-center py-8">
			Загрузка...
		</div>

		<div v-else-if="animeStore.item === null" class="text-center py-8 text-gray-500">
			Не найдено
		</div>

		<div v-else class="space-y-4">
			<p class="text-gray-700">
				{{ animeStore.item?.description || 'Описание отсутствует' }}
			</p>

			<Anilist v-if="anilistInfo" :data="anilistInfo"/>
			<Shiki v-if="shikiInfo" :data="shikiInfo"/>
		</div>
	</div>
</template>

<script setup>
	import { useRoute } from 'vue-router'
	import { onMounted, ref } from 'vue'
	import { useAnimeStore } from "../../stores/anime";
	import Anilist from "./externalInfo/Anilist.vue";
	import Shiki from "./externalInfo/Shiki.vue";

	const animeStore = useAnimeStore();
	const route = useRoute()
	const anilistInfo = ref({})
	const shikiInfo = ref({})

	onMounted(async () => {
		const id = route.params.id
		await animeStore.getItem(id)
		anilistInfo.value = JSON.parse(animeStore.item.anilistInfo || false)
		shikiInfo.value = JSON.parse(animeStore.item.shikiInfo || false)
	})
</script>