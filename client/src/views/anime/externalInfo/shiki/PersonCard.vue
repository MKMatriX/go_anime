<!-- PersonCard.vue -->
<template>
  <div
    class="group bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden hover:shadow-md transition-all duration-200 hover:border-gray-300"
  >
    <!-- Фото + имя -->
    <div class="flex items-center gap-4 p-5">
      <div class="relative flex-shrink-0">
        <img
          v-if="personRole.person.poster?.mainUrl || personRole.person.poster?.originalUrl"
          :src="personRole.person.poster?.mainUrl || personRole.person.poster?.originalUrl"
          :alt="personRole.person.name"
          class="w-16 h-16 rounded-full object-cover border-2 border-gray-100 shadow-sm"
        />
        <div
          v-else
          class="w-16 h-16 rounded-full bg-gray-100 flex items-center justify-center text-gray-400 text-xl font-medium"
        >
          {{ personRole.person.name?.charAt(0) || '?' }}
        </div>
      </div>

      <div class="flex-1 min-w-0">
        <h4 class="text-base font-semibold text-gray-900 truncate group-hover:text-indigo-600 transition-colors">
          {{ personRole.person.name }}
        </h4>
        <p class="text-sm text-gray-600 mt-0.5 line-clamp-2">
          {{ rolesRu }}
        </p>
      </div>
    </div>

    <!-- Английские роли (опционально, мелким шрифтом) -->
    <div v-if="showEnglishRoles" class="px-5 pb-4 text-xs text-gray-500 italic">
      {{ rolesEn }}
    </div>
  </div>
</template>

<script setup>
	import { computed } from 'vue'

	const props = defineProps({
		personRole: {
			type: Object,
			required: true
		},
		// Показывать английские роли (можно отключить для компактности)
		showEnglishRoles: {
			type: Boolean,
			default: false
		}
	})

	const rolesRu = computed(() => {
		return props.personRole.rolesRu.join(', ')
	})

	const rolesEn = computed(() => {
		return props.personRole.rolesEn.join(', ')
	})
</script>