import { defineStore } from 'pinia'
import { loginRequest, getProfileRequest, registerRequest } from '../api/auth'

export const useUserStore = defineStore('auth', {
	state: () => ({
		accessToken: null,
		refreshToken: null,
		user: null,
		loading: false,
		error: null,
	}),
	persist: true,
	getters: {
		isAuthenticated: (state) => !!state.accessToken,
	},
	actions: {
		async login({ login, password }) {
			this.loading = true
			this.error = null
			try {
				const data = await loginRequest({ login, password })
				// { success, message, data: { access_token, refresh_token, user } }
				if (!data.success) {
					throw new Error(data.message || 'Login failed')
				}

				this.accessToken = data.data.access_token
				this.refreshToken = data.data.refresh_token
				this.user = data.data.user
			} catch (err) {
				this.error = err.message || 'Login error'
				throw err
			} finally {
				this.loading = false
			}
		},

		async register({ login, password, confirmPassword }) {
			this.loading = true
			this.error = null
			try {
				const data = await registerRequest({ login, password, confirmPassword })
				// You can adjust according to your real backend response
				if (data.success === false) {
					throw new Error(data.message || 'Register failed')
				}
				return data
			} catch (err) {
				this.error = err.message || 'Register error'
				throw err
			} finally {
				this.loading = false
			}
		},

		async fetchProfile() {
			if (!this.accessToken) return

			this.loading = true
			this.error = null
			try {
				const data = await getProfileRequest()
				this.user = data.user || data
			} catch (err) {
				this.error = err.message || 'Fetch profile error'
				throw err
			} finally {
				this.loading = false
			}
		},

		logout() {
			this.accessToken = null
			this.refreshToken = null
			this.user = null
			this.error = null
		},
	},
})