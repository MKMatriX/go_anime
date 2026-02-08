<template>
	<div class="bg-white shadow rounded-lg p-6">
		<h2 class="text-2xl font-semibold mb-4">{{ animeStore.item.name || "name" }}</h2>

		<div v-if="animeStore.loading" class="text-center py-8">
			Загрузка...
		</div>

		<div v-else-if="animeStore.item === null" class="text-center py-8 text-gray-500">
			Не найдено
		</div>

		<div v-else class="space-y-4">
			<p class="text-gray-700">
				{{ animeStore.item.description || 'Описание отсутствует' }}
			</p>
		</div>
	</div>
</template>

<script setup>
	import { useRoute } from 'vue-router'
	import { onMounted } from 'vue'
	import { useAnimeStore } from "../../stores/anime";

	const animeStore = useAnimeStore();
	const route = useRoute()

	onMounted(async () => {
		const id = route.params.id
		await animeStore.item(id)
	})
</script>