import { apiClient } from './api';
import { Collection, CreateCollectionRequest, UpdateCollectionRequest, Flashcard, CreateFlashcardRequest, UpdateFlashcardRequest, GenerationStats, ApiResponse } from '@/src/types';

export const collectionsService = {
    async list(): Promise<Collection[]> {
        const response = await apiClient.get<ApiResponse<{collections: Collection[]}>>('/collections');
        return response.data.data?.collections || [];
    },

    async create(data: CreateCollectionRequest): Promise<Collection> {
        const response = await apiClient.post<ApiResponse<Collection>>('/collections', data);
        return response.data.data || response.data.data as unknown as Collection;
    },

    async getById(id: string): Promise<Collection> {
        const response = await apiClient.get<ApiResponse<Collection>>(`/collections/${id}`);
        return response.data.data || response.data.data as unknown as Collection;
    },

    async update(id: string, data: UpdateCollectionRequest): Promise<Collection> {
        const response = await apiClient.put<ApiResponse<Collection>>(`/collections/${id}`, data);
        return response.data.data || response.data.data as unknown as Collection;
    },

    async delete(id: string): Promise<void> {
        await apiClient.delete<ApiResponse<void>>(`/collections/${id}`);
    },
};

export const flashcardsService = {
    async listByCollection(collectionId: string): Promise<Flashcard[]> {
        const response = await apiClient.get<ApiResponse<{flashcards: Flashcard[]}>>(`/collections/${collectionId}/flashcards`);
        return response.data.data?.flashcards || [];
    },

    async create(collectionId: string, data: CreateFlashcardRequest): Promise<Flashcard> {
        const response = await apiClient.post<ApiResponse<Flashcard>>(`/collections/${collectionId}/flashcards`, data);
        return response.data.data || response.data.data as unknown as Flashcard;
    },

    async getById(id: string): Promise<Flashcard> {
        const response = await apiClient.get<ApiResponse<Flashcard>>(`/flashcards/${id}`);
        return response.data.data || response.data.data as unknown as Flashcard;
    },

    async update(id: string, data: UpdateFlashcardRequest): Promise<Flashcard> {
        const response = await apiClient.put<ApiResponse<Flashcard>>(`/flashcards/${id}`, data);
        return response.data.data || response.data.data as unknown as Flashcard;
    },

    async delete(id: string): Promise<void> {
        await apiClient.delete<ApiResponse<void>>(`/flashcards/${id}`);
    },

    async generate(collectionId: string, inputType: 'text' | 'topic', content: string): Promise<Flashcard[]> {
        const response = await apiClient.post<ApiResponse<{flashcards: Flashcard[]}>>('/flashcards/generate', {
            input_type: inputType,
            content,
            collection_id: collectionId,
        });
        return response.data.data?.flashcards || [];
    },

    async getGenerationStats(): Promise<GenerationStats> {
        const response = await apiClient.get<ApiResponse<GenerationStats>>('/flashcards/generation-logs');
        return response.data.data || {} as GenerationStats;
    },
};