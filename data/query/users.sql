-- Tabel users (sudah ada)
CREATE TABLE users (
    user_id VARCHAR(30) PRIMARY KEY,
    full_name TEXT NULL,
    username VARCHAR(200) UNIQUE,  -- Sudah otomatis ada indeks dari UNIQUE
    email VARCHAR(200) UNIQUE,     -- Sudah otomatis ada indeks dari UNIQUE
    password TEXT,
    role VARCHAR(50) DEFAULT 'Member',
    is_deleted BOOLEAN DEFAULT FALSE,
    balance INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


ALTER  TABLE  users  ADD COLUMN whatsapp VARCHAR(20);
CREATE INDEX idx_users_role ON users(role); 


CREATE INDEX  idx_users_whatsapp ON users(whatsapp);        -- Berguna untuk filter by role

-- Indeks yang lebih mungkin sering digunakan dalam query:
-- 1. Untuk pencarian kombinasi status aktif + role (contoh: user aktif dengan role 'Admin')
CREATE INDEX idx_users_active_role ON users(role) WHERE is_deleted = FALSE;

-- 2. Untuk pencarian nama (jika sering dipakai di WHERE atau LIKE)
CREATE INDEX idx_users_full_name ON users(full_name);

-- 3. Untuk sorting atau range query balance (contoh: user dengan balance > 1000)
CREATE INDEX idx_users_balance ON users(balance);

-- 4. Composite index untuk pencarian role + status aktif + sorting by created_at
CREATE INDEX idx_users_role_active_created ON users(role, created_at) WHERE is_deleted = FALSE;

CREATE INDEX idx_users_created_at ON users(created_at);



-- Sessions table
CREATE TABLE sessions (
    session_id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(30) NOT NULL REFERENCES users(user_id),
    access_token TEXT NOT NULL,
    ip_address VARCHAR(45),
    user_agent TEXT,
    device_info TEXT,
    last_activity TIMESTAMP,
    expires_at TIMESTAMP,
    is_access BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for sessions table
CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_access_token ON sessions(access_token);
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);
CREATE INDEX idx_sessions_active_sessions ON sessions(user_id, is_access) WHERE is_access = TRUE;
CREATE INDEX idx_sessions_last_activity ON sessions(last_activity);

-- Verification tokens table for password resets and email verification
CREATE TABLE verification_tokens (
    token_id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(30) NOT NULL REFERENCES users(user_id),
    token TEXT NOT NULL,
    token_type VARCHAR(20) NOT NULL CHECK (token_type IN ('email_verification', 'password_reset')),
    expires_at TIMESTAMP NOT NULL,
    is_used BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for verification_tokens table
CREATE INDEX idx_verification_tokens_user_id ON verification_tokens(user_id);
CREATE INDEX idx_verification_tokens_token ON verification_tokens(token);
CREATE INDEX idx_verification_tokens_type ON verification_tokens(token_type);
CREATE INDEX idx_verification_tokens_expires ON verification_tokens(expires_at);

CREATE OR REPLACE FUNCTION immutable_current_timestamp()
RETURNS timestamp AS $$
BEGIN
    RETURN CURRENT_TIMESTAMP;
END;
$$ LANGUAGE plpgsql IMMUTABLE;

CREATE INDEX idx_verification_tokens_active ON verification_tokens(token, token_type, expires_at) 
    WHERE is_used = FALSE AND expires_at > immutable_current_timestamp();