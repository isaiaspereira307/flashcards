-- Schema PostgreSQL do Banco de Dados Flashcards
-- Criado pelos backends Django/FastAPI e verificado pelo Golang

-- ============================================
-- USERS TABLE
-- ============================================
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    username VARCHAR(150) UNIQUE,
    first_name VARCHAR(150),
    last_name VARCHAR(150),
    plan VARCHAR(20) DEFAULT 'free',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);

-- ============================================
-- COLLECTIONS TABLE
-- ============================================
CREATE TABLE IF NOT EXISTS collections (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    is_public BOOLEAN DEFAULT FALSE,
    max_cards INTEGER DEFAULT 10,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_collections_user_id ON collections(user_id);
CREATE INDEX IF NOT EXISTS idx_collections_is_public ON collections(is_public);

-- ============================================
-- FLASHCARDS TABLE
-- ============================================
CREATE TABLE IF NOT EXISTS flashcards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    collection_id UUID NOT NULL REFERENCES collections(id) ON DELETE CASCADE,
    front TEXT NOT NULL,
    back TEXT NOT NULL,
    extra JSONB,
    created_by_ia BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_flashcards_collection_id ON flashcards(collection_id);

-- ============================================
-- SHARES TABLE
-- ============================================
CREATE TABLE IF NOT EXISTS shares (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    collection_id UUID NOT NULL REFERENCES collections(id) ON DELETE CASCADE,
    shared_with_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    permissions VARCHAR(20) DEFAULT 'read',
    share_id VARCHAR(50) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_shares_collection_id ON shares(collection_id);
CREATE INDEX IF NOT EXISTS idx_shares_shared_with_id ON shares(shared_with_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_shares_unique_collection_user ON shares(collection_id, shared_with_id);

-- ============================================
-- GENERATION_LOGS TABLE
-- ============================================
CREATE TABLE IF NOT EXISTS generation_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    date DATE DEFAULT CURRENT_DATE,
    count INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_generation_logs_user_id ON generation_logs(user_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_generation_logs_unique_user_date ON generation_logs(user_id, date);

-- ============================================
-- PAYMENTS TABLE
-- ============================================
CREATE TABLE IF NOT EXISTS payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    subscription_id VARCHAR(255),
    status VARCHAR(20) DEFAULT 'pending',
    start_date TIMESTAMP,
    end_date TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_payments_user_id ON payments(user_id);
CREATE INDEX IF NOT EXISTS idx_payments_status ON payments(status);

-- ============================================
-- GRANT PERMISSIONS
-- ============================================
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO flashcards_user;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO flashcards_user;
