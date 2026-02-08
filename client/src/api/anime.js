import { handleResponse } from "./common";
import api from './axios'

export async function listRequest() {
	return await api.get("/anime")
}

export async function itemRequest(id) {
	return await api.get("/anime/" + id)
}