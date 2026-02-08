import { handleResponse } from "./common";
import api from './axios'

export async function listRequest() {
	return await api.get("/anime")
}

export async function itemRequest(id) {
	return await api.get("/anime/" + id)
}

export async function addRequest(params) {
	return await api.post("/anime", params)
}

export async function deleteRequest(id) {
	return await api.delete("/anime/" + id)
}

export async function editRequest(id, params) {
	return await api.put("/anime/" + id, params)
}