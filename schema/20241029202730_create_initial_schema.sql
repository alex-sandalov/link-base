-- +goose Up
CREATE TABLE users(
    user_id uuid PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL
);

CREATE TABLE referral(
    user_id uuid NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    referred_by_user_id uuid REFERENCES users(user_id),
    PRIMARY KEY (user_id)
);

CREATE TABLE refresh_token(
    user_id uuid NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    refresh_token VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    PRIMARY KEY (user_id, refresh_token)
);

CREATE TABLE referral_code(
    user_id uuid NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    code VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    PRIMARY KEY (user_id, code)
);

-- Indexes
CREATE INDEX idx_user_email ON users (email);
CREATE INDEX idx_referral_referred_by ON referral (user_id);
CREATE INDEX idx_refresh_token_expires_at ON refresh_token (expires_at);
CREATE INDEX idx_referral_code_expires_at ON referral_code (expires_at);

-- +goose Down
DROP TABLE IF EXISTS refresh_token;
DROP TABLE IF EXISTS referral;
DROP TABLE IF EXISTS users;