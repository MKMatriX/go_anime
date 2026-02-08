<template>
  <RouterLink
    :to="{ name: 'anime.detail', params: { id: anime.id } }"
    class="block hover:shadow-lg transition-shadow rounded-lg overflow-hidden bg-white border"
  >
    <div class="anime-item p-4 mb-4">
      <h3 class="text-xl font-bold mb-2">{{ anime.name }}</h3>
      <p class="text-gray-700">
        {{ anime.description || 'Описание отсутствует' }}
      </p>
    </div>
  </RouterLink>
</template>


<script setup lang="ts">
defineProps<{
  anime: {
    id: number | string
    name: string
    description?: string
  }
}>()
</script>

<style scoped>
.anime-item {
  background-color: #f9fafb;
}
</style>