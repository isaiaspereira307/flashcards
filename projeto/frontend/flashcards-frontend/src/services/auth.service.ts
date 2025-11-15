import { apiClient } from './api';
import { LoginRequest, RegisterRequest, AuthResponse } from '@/src/types';

export const authService = {
    async register(email: string, password: string): Promise<AuthResponse> {
        const response = await apiClient.post<AuthResponse>('/auth/register', {
            email,
            password,
        });
        return response.data;
    },

    async login(email: string, password: string): Promise<AuthResponse> {
        const response = await apiClient.post<AuthResponse>('/auth/login', {
            email,
            password,
        });
        return response.data;
    },
    
    async logout(): Promise<void> {
        await apiClient.post('/auth/logout');
    },

    async getMe(): Promise<AuthResponse> {
        const response = await apiClient.get<AuthResponse>('/auth/me');
        return response.data;
    }
}