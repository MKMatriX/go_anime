// src/api/axios.ts  или src/services/api.ts
import axios from 'axios'
import { useUserStore } from '../stores/auth' // ← твой стор
import { BASE_URL } from "./common";

const api = axios.create({
	baseURL: BASE_URL,
	timeout: 10000,
	headers: {
		'Content-Type': 'application/json',
	},
})

// ← Самое важное — request interceptor
api.interceptors.request.use(
	(config) => {
		const userStore = useUserStore()

		if (userStore.accessToken) {
			config.headers.Authorization = `Bearer ${userStore.accessToken}`
		}

		return config
	},
	(error) => Promise.reject(error)
)

api.interceptors.response.use(
	response => response.data,   // сразу возвращаем .data для 2xx

	error => {
		const data = error.response?.data
		if (error.response && error.response.statusText !== 'OK') {
			const message =
				(data && data.message) ||
				(typeof data === 'string' ? data : 'Request failed')
			return Promise.reject(new Error(message))
		}
		return Promise.reject(error) // другие ошибки пробрасываем как есть
	}
)

// (опционально) response interceptor — обработка 401 и refresh токена
api.interceptors.response.use(
	(response) => response,
	async (error) => {
		const originalRequest = error.config
		const userStore = useUserStore()

		// Если 401 и это не запрос на refresh / login
		if (error.response?.status === 401 && !originalRequest._retry) {
			if (originalRequest.url?.includes('/refresh')) {
				userStore.logout()
				return Promise.reject(error)
			}

			originalRequest._retry = true

			try {
				// TODO: implement refreshRequest
				const { data } = await refreshRequest(userStore.refreshToken)
				userStore.accessToken = data.access_token
				userStore.refreshToken = data.refresh_token // если обновляется

				localStorage.setItem('access_token', data.access_token)
				// можно и refresh_token обновить, если бэк его присылает новый

				// Повторяем оригинальный запрос уже с новым токеном
				originalRequest.headers.Authorization = `Bearer ${data.access_token}`
				return api(originalRequest)
			} catch (refreshError) {
				userStore.logout()
				// можно сделать router.push('/login')
				return Promise.reject(refreshError)
			}
		}

		return Promise.reject(error)
	}
)

export default api