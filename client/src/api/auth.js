import { BASE_URL, handleResponse } from "./common";
import api from './axios'

export async function loginRequest({ login, password }) {
	return await api.post("/users/login", {login, password})
}

export async function registerRequest({ login, password, confirmPassword }) {
	return await api.post("/users/login", { login, password, confirmPassword })
}

export async function getProfileRequest() {
	return await api.get('/profile')
}

export async function refreshRequest(refreshToken) {
	return await api.post('/refresh', { refreshToken })
}