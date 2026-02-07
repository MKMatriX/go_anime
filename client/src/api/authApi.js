const BASE_URL = 'http://10.101.39.1/api/v1'

async function handleResponse(response) {
	const contentType = response.headers.get('content-type') || ''
	const data = contentType.includes('application/json')
		? await response.json()
		: await response.text()

	if (!response.ok) {
		const message =
			(data && data.message) ||
			(typeof data === 'string' ? data : 'Request failed')
		throw new Error(message)
	}

	return data
}

export async function loginRequest({ login, password }) {
	const res = await fetch(`${BASE_URL}/users/login`, {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json',
		},
		body: JSON.stringify({ login, password }),
	})
	return handleResponse(res)
}

export async function registerRequest({ login, password, confirmPassword }) {
	const res = await fetch(`${BASE_URL}/users`, {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json',
		},
		body: JSON.stringify({ login, password, confirmPassword }),
	})
	return handleResponse(res)
}

export async function getProfileRequest(accessToken) {
	const res = await fetch(`${BASE_URL}/profile`, {
		method: 'GET',
		headers: {
			'Authorization': `Bearer ${accessToken}`,
		},
	})
	return handleResponse(res)
}