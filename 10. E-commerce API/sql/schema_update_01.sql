-- Run: psql -h localhost -U postgres -d ecommerce_db -f sql/schema_update_01.sql

-- Chỉ hiển thị cảnh báo hoặc lỗi, bỏ qua các thông báo vụn vặt
SET client_min_messages TO warning;

-- Các lệnh bên dưới sẽ không còn hiện NOTICE nữa
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
-- ...

-- 1. Bảng lưu Refresh Token (Cho phép 1 user đăng nhập nhiều thiết bị)
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    token TEXT NOT NULL, -- Hash của token hoặc chính token
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    is_revoked BOOLEAN DEFAULT FALSE -- Dùng để vô hiệu hóa token (Logout)
);
CREATE INDEX idx_refresh_tokens_token ON refresh_tokens(token);

-- 2. Thêm cột cho tính năng Forgot Password vào bảng Users
ALTER TABLE users 
ADD COLUMN IF NOT EXISTS reset_token VARCHAR(255),
ADD COLUMN IF NOT EXISTS reset_token_exp TIMESTAMPTZ;