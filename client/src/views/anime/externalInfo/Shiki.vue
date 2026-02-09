<template>
  <div v-if="data" class="space-y-10">
	<div class="bg-gradient-to-r from-indigo-50 to-blue-50 px-6 py-4 border-b border-gray-200">
		<h5 class="text-lg font-semibold text-gray-800 flex items-center gap-2">
			Shikimori
		</h5>
	</div>
    <!-- HERO секция -->
    <div class="flex flex-col lg:flex-row gap-8 bg-white rounded-2xl shadow-xl overflow-hidden border border-gray-100">
      <!-- Постер -->
      <div class="lg:w-2/5 relative">
        <img
          :src="data.poster?.mainUrl || data.poster?.originalUrl"
          :alt="data.russian || data.name"
          class="w-full h-full object-cover aspect-[2/3] lg:aspect-auto"
        />
        <div class="absolute top-4 left-4 bg-black/70 text-white text-sm font-bold px-3 py-1 rounded-lg">
          {{ data.score }}
        </div>
      </div>

      <!-- Основная информация -->
      <div class="flex-1 p-6 lg:p-10 flex flex-col">
        <h1 class="text-4xl font-bold text-gray-900 leading-tight">
          {{ data.russian || data.name }}
        </h1>
        <p class="text-xl text-gray-600 mt-2">
          {{ data.english }}
        </p>

        <div class="flex flex-wrap gap-3 mt-6">
          <div class="bg-gray-100 px-4 py-1.5 rounded-full text-sm font-medium">
            {{ formatKind(data.kind) }}
          </div>
          <div class="bg-emerald-100 text-emerald-700 px-4 py-1.5 rounded-full text-sm font-medium">
            {{ formatStatus(data.status) || '—' }}
          </div>
          <div class="bg-amber-100 text-amber-700 px-4 py-1.5 rounded-full text-sm font-medium">
            {{ data.episodes }} эпизодов
          </div>
        </div>
      </div>
    </div>

    <!-- Описание -->
    <div class="bg-white rounded-2xl shadow p-8">
      <h2 class="text-2xl font-semibold mb-4">Описание</h2>
      <div
        v-html="data.descriptionHtml || data.description"
        class="prose prose-lg max-w-none text-gray-700 leading-relaxed"
      ></div>
    </div>

    <!-- Персоны (Staff) -->
    <div v-if="props.data.personRoles?.length" class="bg-white rounded-2xl shadow p-8">
      <h2 class="text-2xl font-semibold mb-6">Над аниме работали</h2>

      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <PersonCard
          v-for="role in mainStaff"
          :key="role.id"
          :person-role="role"
        />
      </div>

      <!-- Остальные роли (можно в аккордеон или отдельный блок) -->
      <details class="mt-10">
        <summary class="text-lg font-medium cursor-pointer text-gray-600 hover:text-gray-900">
          Показать весь состав ({{ props.data.personRoles?.length }} человек)
        </summary>
        <div class="mt-6 grid grid-cols-1 md:grid-cols-2 gap-4 text-sm">
          <div v-for="role in otherStaff" :key="role.id" class="flex gap-3">
            <div class="font-medium text-gray-500 w-48 flex-shrink-0">
              {{ role.rolesRu.join(', ') }}
            </div>
            <div class="text-gray-800">{{ role.person.name }}</div>
          </div>
        </div>
      </details>
    </div>

    <!-- Внешние ссылки -->
    <div v-if="props.data.externalLinks?.length" class="bg-white rounded-2xl shadow p-8">
      <h2 class="text-2xl font-semibold mb-5">Ссылки</h2>
      <div class="flex flex-wrap gap-3">
        <a
          v-for="link in props.data.externalLinks"
          :key="link.url"
          :href="link.url"
          target="_blank"
          class="flex items-center gap-2 px-5 py-3 bg-gray-50 hover:bg-gray-100 rounded-xl transition-colors text-sm"
        >
          <span class="text-blue-600">↗</span>
          {{ formatLinkName(link.kind) }}
        </a>
      </div>
    </div>
  </div>
</template>

<script setup>
	import { computed } from 'vue'
	import PersonCard from './shiki/PersonCard.vue'

	const props = defineProps({
		data: {
			type: Object,
			required: true
		}
	})

	// Основные роли (выделяем самые важные)
	const mainStaff = computed(() => {
		if (props.data.personRoles === undefined) {
			return []
		}

		const importantRoles = ['Режиссёр', 'Сценарий', 'Автор оригинала', 'Дизайн персонажей', 'Композитор']
		return props.data.personRoles
			.filter(role => (role.rolesRu || []).some(r => importantRoles.includes(r)))
			.slice(0, 9) // не больше 9 карточек
	})

	// Остальные
	const otherStaff = computed(() =>
		(props.data.personRoles || []).filter(role =>
			!mainStaff.value.includes(role)
		)
	)

	// Вспомогательные функции
	const formatKind = (kind) => ({
		tv: 'TV Сериал',
		movie: 'Фильм',
		ova: 'OVA',
		ona: 'ONA',
		special: 'Спешл'
	}[kind] || kind?.toUpperCase())

	const formatStatus = (status) => ({
		released: 'Вышел',
		ongoing: 'Выходит',
		anons: 'Анонс'
	}[status] || status)

	const formatLinkName = (kind) => ({
		official_site: 'Официальный сайт',
		wikipedia: 'Wikipedia',
		myanimelist: 'MyAnimeList',
		anime_news_network: 'Anime News Network',
		anidb: 'AniDB',
		kinopoisk: 'Кинопоиск'
	}[kind] || kind.replace('_', ' ').toUpperCase())
</script>