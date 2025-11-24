import axios from "axios";
import { API_BASE_URL } from "./api";

const api = axios.create({
	baseURL: API_BASE_URL,
	headers: {
		"Content-Type": "application/json",
	},
});

// Digest API functions
export const digestApi = {
	triggerIntegratedDigest: async (): Promise<void> => {
		await api.post("/digest/trigger/integrated");
	},

	triggerIntegratedDigestForUser: async (userId: number): Promise<void> => {
		await api.post(`/digest/trigger/integrated/user/${userId}`);
	},

	triggerDailyDigest: async (): Promise<void> => {
		await api.post("/digest/trigger");
	},

	triggerDailyDigestForUser: async (userId: number): Promise<void> => {
		await api.post(`/digest/trigger/user/${userId}`);
	},
};
