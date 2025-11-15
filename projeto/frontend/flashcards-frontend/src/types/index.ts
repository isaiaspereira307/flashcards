export type Backend = 'golang' | 'java' | 'fastapi' | 'django';

export type UserPlan = 'free' | 'pro' | 'admin';

export interface User {
    id: string;
    email: string;
    plan: UserPlan;
    created_at: string;
}

export interface AuthResponse {
    success: boolean;
    message?: string;
    data: {
        token: string;
        user: User;
    };
}

export interface LoginRequest {
    email: string;
    password: string;
}

export interface RegisterRequest {
    email: string;
    password: string;
}

export interface Collection {
    id: string;
    user_id: string;
    name: string;
    description?: string;
    is_public: boolean;
    max_cards: number;
    card_count?: number;
    created_at: string;
    updated_at: string;
}

export interface CreateCollectionRequest {
    name: string;
    description?: string;
    is_public?: boolean;
}

export interface UpdateCollectionRequest {
    name?: string;
    description?: string;
    is_public?: boolean;
}

export interface Flashcard {
    id: string;
    collection_id: string;
    front: string;
    back: string;
    created_by_ia: boolean;
    created_at: string;
    updated_at: string;
}

export interface CreateFlashcardRequest {
    front: string;
    back: string;
}

export interface UpdateFlashcardRequest {
    front?: string;
    back?: string;
}

export interface GenerateFlashcardsRequest {
    input_type: 'text' | 'topic';
    content: string;
    collection_id: string;
    quantity?: number;
}

export interface GenerationLog {
    id: string;
    user_id: string;
    collection_id: string;
    generated_count: number;
    status: 'pending' | 'completed' | 'failed';
    created_at: string;
}

export interface GenerationStats {
    generated_today: number;
    daily_limit: number;
    remaining: number;
}

export interface ApiResponse<T = any> {
    success: boolean;
    message?: string;
    data?: T;
    error?: string;
    errors?: Record<string, string[]>;
}

export interface CollectionsResponse {
    success: boolean;
    message?: string;
    data: {
        collections: Collection[];
    };
}

export interface FlashcardsResponse {
    success: boolean;
    message?: string;
    data: {
        flashcards: Flashcard[];
    };
}

export interface Share {
    id: string;
    collection_id: string;
    shared_by: string;
    shared_with?: string;
    share_token: string;
    access_type: 'view' | 'edit';
    created_at: string;
}

export interface Subscription {
    id: string;
    user_id: string;
    plan: UserPlan;
    status: 'active' | 'cancelled' | 'expired';
    started_at: string;
    expires_at?: string;
}