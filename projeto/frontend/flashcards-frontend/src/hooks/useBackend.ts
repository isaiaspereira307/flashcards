import { useState, useEffect } from 'react';
import { Backend } from '@/src/types';
import { apiClient } from '@/src/services/api';

export const useBackend = () => {
    const [backend, setBackendState] = useState<Backend | null>('fastapi');

    useEffect(() => {
        const stored = localStorage.getItem('selectedBackend') as Backend | null;
        if (stored) {
            setBackendState(stored);
            apiClient.setBackend(stored);
        }
    }, []);
    
    return {
        backend,
        setBackend: (b: Backend) => {
            setBackendState(b);
            apiClient.setBackend(b);
        },
    };
};