import { useEffect } from 'react';
import { useAuthStore } from '@/src/store/authStore';

export const useAuth = () => {
    const store = useAuthStore();

    useEffect(() => {
        store.loadFromStorage();
    }, [store]);

    return store;
};