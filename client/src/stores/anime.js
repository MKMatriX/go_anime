import { defineStore } from 'pinia'
import { listRequest, itemRequest } from '../api/anime'

export const useAnimeStore = defineStore('anime', {
	state: () => ({
		loading: false,
		items: null,
		item: null
	}),
	persist: true,
	actions: {
		async list() {
			try {
				this.loading = true
				const data = await listRequest()

				if (!data.success) {
					throw new Error(data.message || 'Anime list request failed')
				}

				this.items = data.data
			} catch (err) {
				this.error = err.message || 'Anime list error'
				throw err
			} finally {
				this.loading = false
			}
		},
		async item(id) {
			try {
				this.loading = true
				const data = await itemRequest(id)

				if (!data.success) {
					throw new Error(data.message || 'Anime item request failed')
				}

				this.item = data.data
			} catch (err) {
				this.error = err.message || 'Anime item error'
				throw err
			} finally {
				this.loading = false
			}
			return this.item
		},
	},
})