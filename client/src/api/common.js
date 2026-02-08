export const BASE_URL = 'http://10.101.39.1/api/v1'

export async function handleResponse(response) {
	const data = response.data

	if (response.statusText !== "OK") {
		const message =
			(data && data.message) ||
			(typeof data === 'string' ? data : 'Request failed')
		throw new Error(message)
	}

	return data
}