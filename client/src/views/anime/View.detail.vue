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


			<div v-if="hasAnilist || hasShiki" class="mt-6">
				<div class="border-b border-gray-200">
					<nav class="-mb-px flex space-x-8" aria-label="Tabs">
						<button
							v-if="hasAnilist"
							@click="activeTab = 'anilist'"
							:class="[
								activeTab === 'anilist'
								? 'border-indigo-500 text-indigo-600'
								: 'border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700',
								'whitespace-nowrap border-b-2 py-4 px-1 text-sm font-medium'
							]"
						>
							Anilist
						</button>

						<button
							v-if="hasShiki"
							@click="activeTab = 'shiki'"
							:class="[
								activeTab === 'shiki'
								? 'border-indigo-500 text-indigo-600'
								: 'border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700',
								'whitespace-nowrap border-b-2 py-4 px-1 text-sm font-medium'
							]"
						>
							Shikimori
						</button>

						<button
							v-if="hasEpisodes"
							@click="activeTab = 'episodes'"
							:class="[
								activeTab === 'episodes'
								? 'border-indigo-500 text-indigo-600'
								: 'border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700',
								'whitespace-nowrap border-b-2 py-4 px-1 text-sm font-medium'
							]"
						>
							Серии
						</button>
					</nav>
				</div>

				<div class="mt-6">
					<Anilist v-if="activeTab === 'anilist' && anilistInfo" :data="anilistInfo" />
					<Shiki v-if="activeTab === 'shiki' && shikiInfo" :data="shikiInfo" />
					<Episodes v-if="activeTab === 'episodes' && episodes" :data="episodes" />
				</div>
			</div>

			<!-- Если нет ни одного источника -->
			<div v-else class="text-center py-6 text-gray-500 italic">
				Дополнительная информация отсутствует
			</div>
		</div>
	</div>
</template>

<script setup>
	import { useRoute } from 'vue-router'
	import { onMounted, ref, computed } from 'vue'
	import { useAnimeStore } from "../../stores/anime";
	import Anilist from "./externalInfo/Anilist.vue";
	import Shiki from "./externalInfo/Shiki.vue";
	import Episodes from "./episode/List.vue";

	const animeStore = useAnimeStore();
	const route = useRoute()

	const anilistInfo = ref(null)
	const shikiInfo = ref(null)
	const episodes = ref(null)
	const activeTab = ref('anilist')

	onMounted(async () => {
		const id = route.params.id
		await animeStore.getItem(id)
		try {
			anilistInfo.value = animeStore.item?.anilistInfo
			? JSON.parse(animeStore.item.anilistInfo)
			: null
		} catch (e) {
			console.warn('Anilist JSON parse error', e)
			anilistInfo.value = null
		}

		try {
			shikiInfo.value = animeStore.item?.shikiInfo
			? JSON.parse(animeStore.item.shikiInfo)
			: null
		} catch (e) {
			console.warn('Shiki JSON parse error', e)
			shikiInfo.value = null
		}

		try {
			episodes.value = animeStore.item?.episodes
			? animeStore.item?.episodes
			: null
		} catch (e) {
			episodes.value = null
		}


		// Умная инициализация активной вкладки
		if (!anilistInfo.value && shikiInfo.value) {
			activeTab.value = 'shiki'
		}
	})

	const hasAnilist = computed(() => !!anilistInfo.value)
	const hasShiki   = computed(() => !!shikiInfo.value)
	const hasEpisodes   = computed(() => !!episodes.value && episodes.value.length > 0)

</script>