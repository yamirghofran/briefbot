import axios from "axios";
import type { User, CreateUserRequest, UpdateUserRequest } from "@/types";
import { API_BASE_URL } from "./api";

const api = axios.create({
	baseURL: API_BASE_URL,
	headers: {
		"Content-Type": "application/json",
	},
});

// User API functions
export const userApi = {
	createUser: async (data: CreateUserRequest): Promise<User> => {
		const response = await api.post<User>("/users", data);
		return response.data;
	},

	getUser: async (id: number): Promise<User> => {
		const response = await api.get<User>(`/users/${id}`);
		return response.data;
	},

	getUserByEmail: async (email: string): Promise<User> => {
		const response = await api.get<User>(`/users/email/${email}`);
		return response.data;
	},

	listUsers: async (): Promise<User[]> => {
		const response = await api.get<User[]>("/users");
		return response.data;
	},

	updateUser: async (id: number, data: UpdateUserRequest): Promise<User> => {
		const response = await api.put<User>(`/users/${id}`, data);
		return response.data;
	},

	deleteUser: async (id: number): Promise<void> => {
		await api.delete(`/users/${id}`);
	},
};
