<template>
	<div v-if="data" class="bg-white rounded-xl shadow-md border border-gray-200">
		<div class="bg-gradient-to-r from-indigo-50 to-blue-50 px-6 py-4 border-b border-gray-200">
			<h5 class="text-lg font-semibold text-gray-800 flex items-center gap-2">
				Anilist
			</h5>
		</div>

		<div class="flex flex-col lg:flex-row gap-8 bg-white shadow-xl overflow-hidden border border-gray-100">
			<div class="lg:w-2/5 relative">
				<img
					v-if="data.coverImage?.large"
					:src="data.coverImage.large"
					class="w-full h-full object-cover aspect-[2/3] lg:aspect-auto "
				/>
				<div class="absolute top-4 left-4 bg-black/70 text-white text-sm font-bold px-3 py-1 rounded-lg">
					{{ data.averageScore }}
				</div>
			</div>

			 <!-- Основная информация -->
			<div class="flex-1 p-6 lg:p-10 flex flex-col">
				<div class="grid grid-cols-[auto,1fr] gap-x-4 gap-y-1">
					<dt class="font-medium text-gray-600 whitespace-nowrap">Название:</dt>
					<dd class="text-gray-900">
						<div v-for="(value, key) in data.title" :key="key" class="mb-0.5">
							<span class="text-xs text-gray-500 uppercase mr-1.5">{{ key }}:</span>
							{{ value || '—' }}
						</div>
					</dd>
				</div>

				<div class="flex flex-wrap gap-3 mt-6">
					<div class="bg-gray-100 px-4 py-1.5 rounded-full text-sm font-medium">
						{{ data.format || '—' }}
					</div>
					<div
						class="bg-emerald-100 text-emerald-700 px-4 py-1.5 rounded-full text-sm font-medium"
						:class="getStatusColor(data.status)"
					>
						{{ formatStatus(data.status) || '—' }}
					</div>
					<div class="bg-amber-100 text-amber-700 px-4 py-1.5 rounded-full text-sm font-medium">
						{{ data.episodes || '—' }} эпизодов
					</div>
				</div>
				<dl class="grid grid-cols-[auto,1fr] gap-x-5 gap-y-3">
					<dt class="font-medium text-gray-600">Дата начала</dt>
					<dd class="text-gray-900">
						{{ formatDate(data.startDate) || '—' }}
					</dd>

					<dt class="font-medium text-gray-600">Жанры</dt>
					<dd class="text-gray-900">
						{{ data.genres?.join(', ') || '—' }}
					</dd>

					<dt class="font-medium text-gray-600">Статус</dt>
					<dd class="text-gray-900">
						<span :class="getStatusColor(data.status)">
							{{ formatStatus(data.status) || '—' }}
						</span>
					</dd>

					<dt class="font-medium text-gray-600">Сезон</dt>
					<dd class="text-gray-900">{{ formatSeason(data.season) || '—' }}</dd>

					<dt class="font-medium text-gray-600 pt-2">Описание</dt>
					<dd class="text-gray-800 col-span-2 leading-relaxed">
						<div v-html="data.description || '<em>Описание отсутствует</em>'" class="prose prose-sm max-w-none"></div>
					</dd>
				</dl>
			</div>

		</div>


	</div>
</template>

<script setup>
	const props = defineProps(['data'])

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

	const formatSeason = (status) => {
		const map = {
			WINTER: 'Зимний',
			SPRING: 'Весенний',
			SUMMER: 'Летний',
			AUTUMN: 'Осенний',
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
</script>