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

			<div v-if="anilistInfo" class="bg-white rounded-xl shadow-md border border-gray-200 overflow-hidden">
				<!-- Заголовок секции -->
				<div class="bg-gradient-to-r from-indigo-50 to-blue-50 px-6 py-4 border-b border-gray-200">
					<h5 class="text-lg font-semibold text-gray-800 flex items-center gap-2">
						Информация об аниме
					</h5>
				</div>

				<!-- Основной контент -->
				<div class="p-6 space-y-5 text-gray-700">
					<!-- Названия (часто их несколько — romaji, english, native) -->
					<div class="grid grid-cols-[auto,1fr] gap-x-4 gap-y-1">
						<dt class="font-medium text-gray-600 whitespace-nowrap">Название:</dt>
						<dd class="text-gray-900">
							<div v-for="(value, key) in anilistInfo.title" :key="key" class="mb-0.5">
								<span class="text-xs text-gray-500 uppercase mr-1.5">{{ key }}:</span>
								{{ value || '—' }}
							</div>
						</dd>
					</div>

					<!-- Остальные поля в grid-формате -->
					<dl class="grid grid-cols-[auto,1fr] gap-x-5 gap-y-3">
						<dt class="font-medium text-gray-600">Дата начала</dt>
						<dd class="text-gray-900">
							{{ formatDate(anilistInfo.startDate) || '—' }}
						</dd>

						<dt class="font-medium text-gray-600">Жанры</dt>
						<dd class="text-gray-900">
							{{ anilistInfo.genres?.join(', ') || '—' }}
						</dd>

						<dt class="font-medium text-gray-600">Средний рейтинг</dt>
						<dd class="text-gray-900 font-semibold">
							{{ anilistInfo.averageScore ? anilistInfo.averageScore + ' / 100' : '—' }}
						</dd>

						<dt class="font-medium text-gray-600">Формат</dt>
						<dd class="text-gray-900">{{ anilistInfo.format || '—' }}</dd>

						<dt class="font-medium text-gray-600">Эпизоды</dt>
						<dd class="text-gray-900">{{ anilistInfo.episodes || '—' }}</dd>

						<dt class="font-medium text-gray-600">Статус</dt>
						<dd class="text-gray-900">
							<span :class="getStatusColor(anilistInfo.status)">
							{{ formatStatus(anilistInfo.status) || '—' }}
							</span>
						</dd>

						<dt class="font-medium text-gray-600">Сезон</dt>
						<dd class="text-gray-900">{{ anilistInfo.season || '—' }}</dd>

						<dt class="font-medium text-gray-600 pt-2">Описание</dt>
						<dd class="text-gray-800 col-span-2 leading-relaxed">
							<div v-html="anilistInfo.description || '<em>Описание отсутствует</em>'" class="prose prose-sm max-w-none"></div>
						</dd>
					</dl>
				</div>
				</div>
		</div>
	</div>
</template>

<script setup>
	import { useRoute } from 'vue-router'
	import { onMounted, ref } from 'vue'
	import { useAnimeStore } from "../../stores/anime";

	const animeStore = useAnimeStore();
	const route = useRoute()
	const anilistInfo = ref({})

	const formatDate = (dateObj) => {
		if (!dateObj?.year) return null
		const d = dateObj.day?.toString().padStart(2, '0') || '??'
		const m = dateObj.month?.toString().padStart(2, '0') || '??'
		return `${d}.${m}.${dateObj.year}`
	}

	const formatStatus = (status) => {
		const map = {
			FINISHED: 'Завершено',
			RELEASING: 'Выходит',
			NOT_YET_RELEASED: 'Анонсировано',
			CANCELLED: 'Отменено',
			HIATUS: 'На паузе'
		}
		return map[status] || status
	}

	const getStatusColor = (status) => {
		const colors = {
			FINISHED: 'text-green-700 font-medium',
			RELEASING: 'text-blue-700 font-medium animate-pulse',
			NOT_YET_RELEASED: 'text-purple-700',
			CANCELLED: 'text-red-700 line-through',
			HIATUS: 'text-amber-700'
		}
		return colors[status] || 'text-gray-600'
	}

	onMounted(async () => {
		const id = route.params.id
		await animeStore.getItem(id)
		anilistInfo.value = JSON.parse(animeStore.item.anilistInfo)
	})
</script>