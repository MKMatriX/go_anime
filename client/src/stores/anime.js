import { defineStore } from 'pinia'
import { listRequest, itemRequest, addRequest, deleteRequest, editRequest } from '../api/anime'

export const useAnimeStore = defineStore('anime', {
	state: () => ({
		loading: false,
		error: null,
		items: null,
		item: null
	}),
	persist: true,
	actions: {
		async getList() {
			this.loading = true
			this.error = null
			try {
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
		async getItem(id) {
			this.loading = true
			this.error = null
			try {
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
		async add(params) {
			this.loading = true
			this.error = null
			try {
				const data = await addRequest(params)
				this.items.push(data.data)
				return data.data
			} catch (error) {
				this.error = err.message || 'Anime add error'
				throw err
			} finally {
				this.loading = false
			}
		},
		async edit(id, params) {
			this.loading = true
			this.error = null
			try {
				const data = await editRequest(id, params)
				let item = this.items.find((item) => item.id == id)
				let index = this.items.indexOf(item)
				this.items[index] = data.data
				return data.data
			} catch (error) {
				this.error = err.message || 'Anime add error'
				throw err
			} finally {
				this.loading = false
			}
		},
		async delete(id) {
			this.loading = true
			this.error = null
			try {
				const data = await deleteRequest(id)
				this.items = this.items.filter((item) => item.id != id)
			} catch (error) {
				this.error = err.message || 'Anime add error'
				throw err
			} finally {
				this.loading = false
			}
		},
	},
})