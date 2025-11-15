import axios, { AxiosInstance } from "axios";
import { Backend } from "../types";

class ApiClient {
    private instances: Record<Backend, AxiosInstance> = {} as Record<Backend, AxiosInstance>;
    private currentBackend: Backend = 'fastapi';

    constructor() {
        this.initializeInstances();
    }

    private initializeInstances() {
        const backends: Backend[] = ['django', 'fastapi', 'java', 'golang'];
        const urls: Record<Backend, string> = {
            django: 'http://localhost:8000/api',
            fastapi: 'http://localhost:8001/api',
            java: 'http://localhost:8002/api',
            golang: 'http://localhost:8003/api',
        };

        backends.forEach((backend) => {
            const instance = axios.create({
                baseURL: urls[backend],
                headers: { 'Content-Type': 'application/json' },
            });
            instance.interceptors.request.use((config) => {
                const token = localStorage.getItem('token');
                if (token) {
                    config.headers['Authorization'] = `Bearer ${token}`;
                }
                return config;
            });

            instance.interceptors.response.use(
                (response) => response,
                (error) => {
                    if (error.response?.status === 401) {
                        localStorage.removeItem('token');
                        window.location.href = '/login';
                    }
                    return Promise.reject(error);
                }
            );

            this.instances[backend] = instance;
        })
    }

    setBackend(backend: Backend) {
        this.currentBackend = backend;
        localStorage.setItem('selectedBackend', backend);
    }

    getBackend(): Backend {
        return this.currentBackend;
    }

    getInstance() {
        return this.instances[this.currentBackend];
    }

    get<T = any>(url: string, config?: any) {
        return this.getInstance().get<T>(url, config);
    }

    post<T = any>(url: string, data?: any, config?: any) {
        return this.getInstance().post<T>(url, data, config);
    }

    put<T = any>(url: string, data?: any, config?: any) {
        return this.getInstance().put<T>(url, data, config);
    }

    delete<T = any>(url: string, config?: any) {
        return this.getInstance().delete<T>(url, config);
    }
}

export const apiClient = new ApiClient();