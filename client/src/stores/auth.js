import { defineStore } from 'pinia'
import { loginRequest, getProfileRequest, registerRequest } from '../api/authApi'

export const useAuthStore = defineStore('auth', {
	state: () => ({
		accessToken: localStorage.getItem('access_token') || null,
		refreshToken: localStorage.getItem('refresh_token') || null,
		user: JSON.parse(localStorage.getItem('user') || 'null'),
		loading: false,
		error: null,
	}),
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

				localStorage.setItem('access_token', this.accessToken)
				localStorage.setItem('refresh_token', this.refreshToken)
				localStorage.setItem('user', JSON.stringify(this.user))
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
				const data = await getProfileRequest(this.accessToken)
				this.user = data.user || data
				localStorage.setItem('user', JSON.stringify(this.user))
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

			localStorage.removeItem('access_token')
			localStorage.removeItem('refresh_token')
			localStorage.removeItem('user')
		},
	},
})