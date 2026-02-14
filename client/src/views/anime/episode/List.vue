<template>
  <div class="w-full max-w-3xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
    <div class="space-y-4">
      <div
        v-for="episode in sortedEpisodes"
        :key="episode.id"
        class="group relative bg-gray-900/60 backdrop-blur-sm border border-gray-700/70 rounded-xl overflow-hidden hover:border-indigo-500/50 transition-all duration-300 hover:shadow-xl hover:shadow-indigo-500/10"
      >
        <div class="flex items-center gap-4 p-4 sm:p-5">
          <!-- Номер эпизода -->
          <div class="flex-shrink-0 w-14 h-14 rounded-lg bg-gradient-to-br from-indigo-600/80 to-purple-600/80 flex items-center justify-center text-white font-bold text-xl sm:text-2xl shadow-md">
            {{ episode.episodeNumber }}
          </div>

          <!-- Основная информация -->
          <div class="flex-1 min-w-0">
            <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-1 sm:gap-3">
              <h3 class="text-lg font-medium text-gray-100 truncate group-hover:text-indigo-400 transition-colors">
                {{ cleanEpisodeName(episode.name) }}
              </h3>

              <div class="flex items-center gap-3 flex-wrap">
                <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-gray-700/70 text-gray-300">
                  {{ episode.width }}
                </span>
                <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-950/60 text-blue-300 border border-blue-800/40">
                  {{ episode.translator }}
                </span>
              </div>
            </div>

            <!-- Дополнительно (если нужно показать полное имя или другие данные) -->
            <!-- <p class="mt-1 text-sm text-gray-400 line-clamp-1">{{ episode.name }}</p> -->
          </div>

          <!-- Кнопка скачивания -->
          <div class="flex-shrink-0">
            <a
              :href="episode.torrentUrl"
              target="_blank"
              rel="noopener noreferrer"
              class="inline-flex items-center gap-2 px-4 py-2.5 bg-indigo-600/90 hover:bg-indigo-600 active:bg-indigo-700 text-white text-sm font-medium rounded-lg transition-colors shadow-md hover:shadow-lg hover:shadow-indigo-500/20 focus:outline-none focus:ring-2 focus:ring-indigo-500/50"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
              </svg>
              Torrent
            </a>
          </div>
        </div>
      </div>
    </div>

    <p v-if="!sortedEpisodes?.length" class="text-center text-gray-500 py-12">
      Нет доступных эпизодов
    </p>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  data: {
    type: Array,
    default: () => []
  }
})

const sortedEpisodes = computed(() => {
  return [...props.data].sort((a, b) => b.episodeNumber - a.episodeNumber)
})

const cleanEpisodeName = (name) => {
  // Убираем повторяющееся название аниме + номер из поля name
  return name.replace(/^(.*?\s*-\s*\d{2,})\s*-\s*\d{2,}$/, '$1').trim()
}
</script>