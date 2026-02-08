<template>
	<div class="bg-white shadow rounded-lg p-6">
		<h2 class="text-2xl font-semibold mb-4">
			Изменить аниме {{ animeStore.item?.name || "%animeName%" }}
		</h2>

		<form @submit.prevent="onSubmit" class="space-y-4">
			<div>
				<label class="block text-sm font-medium text-gray-700 mb-1">
					Название
				</label>
				<input
					v-model="form.name"
					type="text"
					class="w-full border rounded-md px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
				/>
				<p v-if="errors.name" class="text-xs text-red-600 mt-1">
					{{ errors.name }}
				</p>
			</div>

			<div>
				<label class="block text-sm font-medium text-gray-700 mb-1">
					Описание
				</label>
				<input
					v-model="form.description"
					type="text"
					class="w-full border rounded-md px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
				/>
				<p v-if="errors.description" class="text-xs text-red-600 mt-1">
					{{ errors.description }}
				</p>
			</div>

			<p v-if="animeStore.error" class="text-sm text-red-600">
				{{ animeStore.error }}
			</p>

			<button
				type="submit"
				:disabled="animeStore.loading"
				class="w-full bg-green-600 hover:bg-green-700 text-white text-sm font-medium py-2 px-4 rounded-md disabled:opacity-60"
			>
				<span v-if="animeStore.loading">Сохраняем...</span>
				<span v-else>Сохранить</span>
			</button>
		</form>
	</div>
</template>

<script setup lang="ts">
import { reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAnimeStore } from '../../stores/anime'

const props = defineProps<{
	id?: string | number
}>()

const animeStore = useAnimeStore()
const router = useRouter()


onMounted(async () => {
	let anime = await animeStore.getItem(props.id)
	form.name = anime.name
	form.description = anime.description
})

const form = reactive({
  name: '',
  description: '',
})

const errors = reactive({
  name: '',
  description: '',
})

const validate = () => {
  errors.name = ''
  errors.description = ''

  if (!form.name.trim()) {
    errors.name = 'Name is required'
  }

  return !errors.name
}

const onSubmit = async () => {
  if (!validate()) return
  try {
    let anime = await animeStore.edit(props.id, {
      name: form.name,
      description: form.description,
    })

	router.push({ name: 'anime.detail', params: {id: anime.id} })
  } catch {
    // error already set in auth.error
  }
}
</script>