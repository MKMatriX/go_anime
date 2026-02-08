<template>
	<div class="bg-white shadow rounded-lg p-6">
		<h2 class="text-2xl font-semibold mb-4">Anime list</h2>

		<div v-if="anime.loading" class="text-center py-8">
			Загрузка...
		</div>

		<div v-else-if="anime.items === null || anime.items.length === 0" class="text-center py-8 text-gray-500">
			Список пуст
		</div>

		<div v-else class="space-y-4">
			<AnimeItem
				v-for="animeItem in anime.items"
				:key="animeItem.id || animeItem.name"
				:anime="animeItem"
			/>
		</div>
	</div>
</template>

<script setup>
	import { onMounted } from 'vue'
	import { useAnimeStore } from "../../stores/anime";
	import AnimeItem from './View.list.vue'

	const anime = useAnimeStore();

	onMounted(() => {
		anime.list()
	})
</script>