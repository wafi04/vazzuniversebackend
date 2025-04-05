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

CREATE INDEX idx_users_role ON users(role);         -- Berguna untuk filter by role

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