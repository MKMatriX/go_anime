<template>
  <div class="anime-card relative overflow-hidden rounded-lg border bg-white shadow-sm hover:shadow-md transition-shadow">
    <!-- Основная область карточки — кликабельная, ведёт в деталку -->
    <RouterLink
      :to="{ name: 'anime.detail', params: { id: anime.id } }"
      class="block p-4"
    >
      <h3 class="text-xl font-bold mb-2 line-clamp-2">{{ anime.name }}</h3>
      <p class="text-gray-700 line-clamp-3 mb-3">
        {{ anime.description || 'Описание отсутствует' }}
      </p>
    </RouterLink>

    <!-- Кнопки действий — не перехватывают клик по всей карточке -->
    <div class="absolute top-3 right-3 flex gap-2">
      <button
        @click.stop="editAnime"
        class="flex items-center justify-center w-9 h-9 rounded-full bg-blue-100 text-blue-700 hover:bg-blue-200 transition-colors"
        title="Редактировать"
      >
        ✏️
      </button>

      <button
        @click.stop="confirmDelete"
        class="flex items-center justify-center w-9 h-9 rounded-full bg-red-100 text-red-700 hover:bg-red-200 transition-colors"
        title="Удалить"
      >
        🗑
      </button>
    </div>
  </div>
</template>


<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useAnimeStore } from '../../stores/anime'

const props = defineProps<{
  anime: {
    id: number | string
    name: string
    description?: string
  }
}>()

const router = useRouter()
const animeStore = useAnimeStore()

const editAnime = () => {
  router.push({
    name: 'anime.edit',
    params: { id: props.anime.id }
  })
}


const confirmDelete = () => {
  if (!confirm(`Удалить аниме "${props.anime.name}"?`)) {
    return
  }

  animeStore.delete(props.anime.id)
    .then(() => {
      console.log('Аниме удалено')
    })
    .catch(err => {
      alert('Не удалось удалить: ' + err.message)
    })
}

</script>

<style scoped>
.anime-card {
  position: relative;
  background-color: #f9fafb;
}

.anime-card:hover .anime-card__actions {
  opacity: 1;
}
</style>