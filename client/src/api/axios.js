// src/api/axios.ts  или src/services/api.ts
import axios from 'axios'
import { useUserStore } from '../stores/auth' // ← твой стор
import { BASE_URL } from "./common";
import { refreshRequest } from "./auth";

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

		if (error.response?.status === 401) {
			return Promise.reject(error)
		}

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
	(response) => {return response},
	async (error) => {
		console.log("Axios error occurred:", {
			message: error.message,
			name: error.name,           // обычно "AxiosError"
			code: error.code,           // например ECONNABORTED, ERR_NETWORK
			status: error.response?.status,
			statusText: error.response?.statusText,
			data: error.response?.data,     // ← здесь чаще всего тело ошибки от сервера
			headers: error.response?.headers,
			config: {
				url: error.config?.url,
				method: error.config?.method,
				headers: error.config?.headers,
			},
			request: error.request ? "request exists" : "no request",
		});


		const originalRequest = error.config
		const userStore = useUserStore()

		if (error.response?.status === 401 && !originalRequest._retry) {
			console.log("making refresh request");

			if (originalRequest.url?.includes('/refresh')) {
				userStore.logout()
				return Promise.reject(error)
			}

			originalRequest._retry = true

			try {
				const { data } = await refreshRequest(userStore.refreshToken)
				userStore.accessToken = data.access_token
				userStore.refreshToken = data.refresh_token

				originalRequest.headers.Authorization = `Bearer ${data.access_token}`
				return api(originalRequest)
			} catch (refreshError) {
				userStore.logout()
				return Promise.reject(refreshError)
			}
		}

		return Promise.reject(error)
	}
)

export default api