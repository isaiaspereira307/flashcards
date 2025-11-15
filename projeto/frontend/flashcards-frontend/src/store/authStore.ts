import { create } from 'zustand';
import { User } from '@/src/types';
import { authService } from '@/src/services/auth.service';

interface AuthStore {
    user: User | null;
    token: string | null;
    isAuthenticated: boolean;
    isLoading: boolean;

    login: (email: string, password: string) => Promise<void>;
    register: (email: string, password: string) => Promise<void>;
    logout: () => void;
    loadFromStorage: () => void;
    setUser: (user: User) => void;
}

export const useAuthStore = create<AuthStore>((set) => ({
    user: null,
    token: null,
    isAuthenticated: false,
    isLoading: false,

    login: async (email: string, password: string) => {
        set({ isLoading: true });
        try {
            const response = await authService.login(email, password);
            set({
                user: response.data.user,
                token: response.data.token,
                isAuthenticated: true,
                isLoading: false,
            });
            localStorage.setItem('token', response.data.token);
            localStorage.setItem('user', JSON.stringify(response.data.user));
        } catch (error) {
            set({ isLoading: false });
            throw error;
        }
    },

    register: async (email: string, password: string) => {
        set({ isLoading: true });
        try {
            const response = await authService.register(email, password);
            set({
                user: response.data.user,
                token: response.data.token,
                isAuthenticated: true,
                isLoading: false,
            });
            localStorage.setItem('token', response.data.token);
            localStorage.setItem('user', JSON.stringify(response.data.user));
        } catch (error) {
            set({ isLoading: false });
            throw error;
        }
    },

    logout: () => {
        set({ user: null, token: null, isAuthenticated: false });
        localStorage.removeItem('token');
        localStorage.removeItem('user');
    },

    setUser: (user: User) => {
        set({ user, isAuthenticated: true });
        localStorage.setItem('user', JSON.stringify(user));
    },

    loadFromStorage: () => {
        const token = localStorage.getItem('token');
        const userStr = localStorage.getItem('user');
        if (token && userStr) {
            try {
                set({
                    token,
                    user: JSON.parse(userStr),
                    isAuthenticated: true,
                });
            } catch {
                localStorage.removeItem('token');
                localStorage.removeItem('user');
            }
        }
    },
}));